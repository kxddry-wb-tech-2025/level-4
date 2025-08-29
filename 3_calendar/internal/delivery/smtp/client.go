package smtp

import (
	"calendar/internal/config"
	"calendar/internal/delivery"
	"context"

	"gopkg.in/gomail.v2"
)

// EmailClient is the SMTP email client
type EmailClient struct {
	dialer *gomail.Dialer
}

// NewEmailClient creates a new SMTP email client
func NewEmailClient(cfg *config.EmailConfig) *EmailClient {
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	// Configure TLS based on config
	if !cfg.TLSEnabled {
		dialer.SSL = false
		dialer.TLSConfig = nil
		dialer.Username = ""
		dialer.Password = ""
	}

	return &EmailClient{
		dialer: dialer,
	}
}

// Ping checks if the SMTP connection is working
func (c *EmailClient) Ping() error {
	_, err := c.dialer.Dial()
	return err
}

// SendEmails sends multiple emails reusing the same connection for performance
func (c *EmailClient) SendEmails(ctx context.Context, msgs []*delivery.EmailMessage) error {
	s, err := c.dialer.Dial()
	if err != nil {
		return err
	}
	defer s.Close()

	for _, msg := range msgs {
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

		if err := gomail.Send(s, m); err != nil {
			return err
		}
	}

	return nil
}

// TestConnection tests the SMTP connection without sending an email
func (c *EmailClient) TestConnection() error {
	_, err := c.dialer.Dial()
	if err != nil {
		return err
	}

	return nil
}
