package notify

import (
	"calendar/internal/config"
	"calendar/internal/delivery"
	"calendar/internal/models"
	"context"
	"testing"
)

type fakeEmailClient struct{ sent *delivery.EmailMessage }

func (f *fakeEmailClient) SendEmail(ctx context.Context, msg *delivery.EmailMessage) error {
	f.sent = msg
	return nil
}

func TestSender_Email(t *testing.T) {
	fec := &fakeEmailClient{}
	s := NewSender(fec, &config.EmailConfig{Username: "from@example.com"})
	if err := s.Send(context.Background(), models.Notification{Channel: "email", Recipient: "a@b.c", Message: "hi"}); err != nil {
		t.Fatalf("send: %v", err)
	}
	if fec.sent == nil || fec.sent.From != "from@example.com" || fec.sent.To[0] != "a@b.c" || fec.sent.Body != "hi" {
		t.Fatalf("unexpected msg: %#v", fec.sent)
	}
}
