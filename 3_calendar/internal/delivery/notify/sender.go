package notify

import (
	"calendar/internal/config"
	"calendar/internal/delivery"
	"calendar/internal/models"
	"context"
	"fmt"
)

// EmailClient is the interface for the email client
type EmailClient interface {
	SendEmail(ctx context.Context, msg *delivery.EmailMessage) error
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
