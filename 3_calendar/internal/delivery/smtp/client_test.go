package smtp

import (
	"calendar/internal/config"
	"calendar/internal/delivery"
	"context"
	"testing"
)

func TestEmailClient_NoRecipients(t *testing.T) {
	c := NewEmailClient(&config.EmailConfig{Host: "localhost", Port: 1025})
	if err := c.SendEmails(context.Background(), []*delivery.EmailMessage{
		{From: "a", To: []string{"b"}},
	}); err == nil {
		t.Fatalf("expected error when no recipients")
	}
}
