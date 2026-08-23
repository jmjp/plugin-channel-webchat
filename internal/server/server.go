package server

import (
	"context"
	pluginv1 "fluxa/proto/fluxa/v1"
	"time"
)

type Server struct {
	pluginv1.UnimplementedChannelPluginServerServer
	Inbound chan *pluginv1.InboundMessage
}

func New() *Server { return &Server{Inbound: make(chan *pluginv1.InboundMessage, 32)} }
func (s *Server) SendMessage(context.Context, *pluginv1.OutboundMessageRequest) (*pluginv1.OutboundMessageResponse, error) {
	return &pluginv1.OutboundMessageResponse{Status: "sent"}, nil
}
func (s *Server) HandleWebhook(ctx context.Context, r *pluginv1.WebhookPayloadRequest) (*pluginv1.WebhookPayloadResponse, error) {
	// Usa um timeout curto para a leitura do canal — se não há webhook
	// na fila, devolvemos resposta vazia em vez de bloquear o gRPC
	// indefinidamente. O gRPC já tem seu próprio timeout (default 3s
	// no GRPCCaller); queremos algo menor para responder rápido.
	readCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()
	select {
	case msg := <-s.Inbound:
		return &pluginv1.WebhookPayloadResponse{Messages: []*pluginv1.InboundMessage{msg}, HttpStatusCode: 200}, nil
	case <-readCtx.Done():
		return &pluginv1.WebhookPayloadResponse{HttpStatusCode: 200}, nil
	}
}
func (s *Server) HealthCheck(context.Context, *pluginv1.Empty) (*pluginv1.HealthStatus, error) {
	return &pluginv1.HealthStatus{Healthy: true}, nil
}
func (s *Server) Enqueue(customer, text string) {
	s.Inbound <- &pluginv1.InboundMessage{Channel: "webchat", CustomerRef: customer, Text: text, TimestampUnixNano: time.Now().UnixNano()}
}
