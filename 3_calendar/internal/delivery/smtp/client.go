package smtp

import (
	"calendar/internal/config"
	"calendar/internal/delivery"
	"context"
	"fmt"

	"gopkg.in/gomail.v2"
)

// EmailClient is the SMTP email client
type EmailClient struct {
	dialer *gomail.Dialer
}

// NewEmailClient creates a new SMTP email client
func NewEmailClient(cfg *config.EmailConfig) *EmailClient {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	return &EmailClient{
		dialer: dialer,
	}
}

// SendEmail sends an email using the configured SMTP server
func (c *EmailClient) SendEmail(ctx context.Context, msg *delivery.EmailMessage) error {
	if len(msg.To) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", msg.From)
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)

	if msg.IsHTML {
		m.SetBody("text/html", msg.Body)
	} else {
		m.SetBody("text/plain", msg.Body)
	}

	for _, filePath := range msg.FilePaths {
		m.Attach(filePath)
	}

	out := make(chan error)
	go func(ctx context.Context) {
		defer close(out)

		select {
		case <-ctx.Done():
			out <- ctx.Err()
		case out <- c.dialer.DialAndSend(m):
		}
	}(ctx)

	return <-out
}

// TestConnection tests the SMTP connection without sending an email
func (c *EmailClient) TestConnection() error {
	_, err := c.dialer.Dial()
	if err != nil {
		return err
	}

	return nil
}
