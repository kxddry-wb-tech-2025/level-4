package notify

import (
	"calendar/internal/config"
	"calendar/internal/delivery"
	"calendar/internal/models"
	"context"
	"fmt"
)

type EmailClient interface {
	SendEmail(ctx context.Context, msg *delivery.EmailMessage) error
}

type Sender struct {
	emailClient EmailClient
	config      *config.EmailConfig
}

func NewSender(emailClient EmailClient, config *config.EmailConfig) *Sender {
	return &Sender{
		emailClient: emailClient,
		config:      config,
	}
}

func (s *Sender) Send(ctx context.Context, notification models.Notification) error {
	switch notification.Channel {
	case "email":
		return s.emailClient.SendEmail(ctx, &delivery.EmailMessage{
			From:      s.config.Username,
			To:        []string{notification.Recipient},
			Subject:   "Event Reminder",
			Body:      notification.Message,
			IsHTML:    true,
			FilePaths: nil,
		})
	default:
		return fmt.Errorf("unsupported channel: %s", notification.Channel)
	}
}
