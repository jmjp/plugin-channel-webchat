package main

import (
	"encoding/json"
	pluginv1 "fluxa/proto/fluxa/v1"
	"github.com/jmjp/plugin-channel-webchat/internal/server"
	"github.com/jmjp/plugin-channel-webchat/internal/webhook"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	s := server.New()
	socket := os.Getenv("FLUXA_PLUGIN_SOCKET")
	if socket == "" {
		socket = "/tmp/fluxa/plugin-channel-webchat.sock"
	}
	_ = os.Remove(socket)
	l, e := net.Listen("unix", socket)
	if e != nil {
		panic(e)
	}
	g := grpc.NewServer()
	pluginv1.RegisterChannelPluginServerServer(g, s)
	go g.Serve(l)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /webhook", func(w http.ResponseWriter, r *http.Request) {
		p, e := webhook.Parse(mustRead(r))
		if e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		s.Enqueue(p.CustomerRef, p.Text)
		w.WriteHeader(202)
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	go http.ListenAndServe("localhost:9100", mux)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	g.GracefulStop()
}
func mustRead(r *http.Request) []byte {
	defer r.Body.Close()
	b := make([]byte, r.ContentLength)
	_, _ = r.Body.Read(b)
	return b
}
