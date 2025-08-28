package notify

import (
	"calendar/internal/config"
	"calendar/internal/delivery"
	"calendar/internal/models"
	"context"
)

// EmailClient is the interface for the email client
type EmailClient interface {
	SendEmails(ctx context.Context, msgs []*delivery.EmailMessage) error
}

// Sender is the struct for the email sender
type Sender struct {
	emailClient EmailClient
	config      *config.EmailConfig
}

// NewSender creates a new email sender
func NewSender(emailClient EmailClient, config *config.EmailConfig) *Sender {
	return &Sender{
		emailClient: emailClient,
		config:      config,
	}
}

// Send sends an email notification
func (s *Sender) Send(ctx context.Context, notifications []models.Notification) error {
	msgs := make([]*delivery.EmailMessage, len(notifications))
	for i, notification := range notifications {
		if notification.Channel != "email" {
			continue
		}

		msgs[i] = &delivery.EmailMessage{
			From:      s.config.Username,
			To:        []string{notification.Recipient},
			Subject:   "Event Reminder",
			Body:      notification.Message,
			IsHTML:    true,
			FilePaths: nil,
		}
	}

	return s.emailClient.SendEmails(ctx, msgs)
}
