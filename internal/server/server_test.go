package server

import (
	"context"
	"testing"
	"time"

	pluginv1 "fluxa/proto/fluxa/v1"
)

// TestSendMessage devolve status sent (PRD §6.4 MessageStatusSent).
func TestSendMessage(t *testing.T) {
	s := New()
	resp, err := s.SendMessage(context.Background(), &pluginv1.OutboundMessageRequest{
		ChannelId:   "webchat_test",
		CustomerRef: "u-1",
		Text:        "olá",
	})
	if err != nil {
		t.Fatalf("SendMessage: %v", err)
	}
	if resp.GetStatus() != "sent" {
		t.Fatalf("status = %q, want sent", resp.GetStatus())
	}
}

// TestHealthCheck sempre saudável no MVP.
func TestHealthCheck(t *testing.T) {
	s := New()
	resp, err := s.HealthCheck(context.Background(), &pluginv1.Empty{})
	if err != nil {
		t.Fatalf("HealthCheck: %v", err)
	}
	if !resp.GetHealthy() {
		t.Fatal("HealthCheck unhealthy")
	}
}

// TestHandleWebhook consome uma mensagem enfileirada e devolve normalizada.
func TestHandleWebhook(t *testing.T) {
	s := New()
	s.Enqueue("u-42", "olá do webchat")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	resp, err := s.HandleWebhook(ctx, &pluginv1.WebhookPayloadRequest{
		RawPayload: []byte(`{}`),
	})
	if err != nil {
		t.Fatalf("HandleWebhook: %v", err)
	}
	if resp.GetHttpStatusCode() != 200 {
		t.Fatalf("status code = %d, want 200", resp.GetHttpStatusCode())
	}
	msgs := resp.GetMessages()
	if len(msgs) != 1 {
		t.Fatalf("messages = %d, want 1", len(msgs))
	}
	if msgs[0].GetCustomerRef() != "u-42" {
		t.Fatalf("customer_ref = %q, want u-42", msgs[0].GetCustomerRef())
	}
	if msgs[0].GetText() != "olá do webchat" {
		t.Fatalf("text = %q, want 'olá do webchat'", msgs[0].GetText())
	}
	if msgs[0].GetChannel() != "webchat" {
		t.Fatalf("channel = %q, want webchat", msgs[0].GetChannel())
	}
}

// TestHandleWebhook_EmptyQueue devolve resposta vazia quando não há
// mensagens aguardando (e não bloqueia indefinidamente).
func TestHandleWebhook_EmptyQueue(t *testing.T) {
	s := New()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	resp, err := s.HandleWebhook(ctx, &pluginv1.WebhookPayloadRequest{})
	if err != nil {
		t.Fatalf("HandleWebhook: %v", err)
	}
	if resp.GetHttpStatusCode() != 200 {
		t.Fatalf("status = %d, want 200", resp.GetHttpStatusCode())
	}
	if len(resp.GetMessages()) != 0 {
		t.Fatalf("messages = %d, want 0", len(resp.GetMessages()))
	}
}
