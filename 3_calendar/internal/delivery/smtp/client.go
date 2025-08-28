package smtp

import (
	"calendar/internal/config"
	"fmt"
	"net/smtp"
	"strings"
)

type EmailClient struct {
	config *config.EmailConfig
	auth   smtp.Auth
}

type EmailMessage struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

// NewEmailClient creates a new SMTP email client
func NewEmailClient(cfg *config.EmailConfig) *EmailClient {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	return &EmailClient{
		config: cfg,
		auth:   auth,
	}
}

// SendEmail sends an email using the configured SMTP server
func (c *EmailClient) SendEmail(msg *EmailMessage) error {
	if len(msg.To) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	headers := make([]string, 0)
	headers = append(headers, fmt.Sprintf("From: %s", c.config.Username))
	headers = append(headers, fmt.Sprintf("To: %s", strings.Join(msg.To, ", ")))
	headers = append(headers, fmt.Sprintf("Subject: %s", msg.Subject))

	if msg.IsHTML {
		headers = append(headers, "MIME-Version: 1.0")
		headers = append(headers, "Content-Type: text/html; charset=UTF-8")
	} else {
		headers = append(headers, "Content-Type: text/plain; charset=UTF-8")
	}

	headers = append(headers, "")
	headers = append(headers, msg.Body)

	emailBody := strings.Join(headers, "\r\n")

	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	err := smtp.SendMail(addr, c.auth, c.config.Username, msg.To, []byte(emailBody))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendSimpleEmail is a convenience method for sending simple text emails
func (c *EmailClient) SendSimpleEmail(to []string, subject, body string) error {
	msg := &EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	}

	return c.SendEmail(msg)
}

// SendHTMLEmail is a convenience method for sending HTML emails
func (c *EmailClient) SendHTMLEmail(to []string, subject, htmlBody string) error {
	msg := &EmailMessage{
		To:      to,
		Subject: subject,
		Body:    htmlBody,
		IsHTML:  true,
	}

	return c.SendEmail(msg)
}

// TestConnection tests the SMTP connection without sending an email
func (c *EmailClient) TestConnection() error {
	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	if err := client.Auth(c.auth); err != nil {
		return fmt.Errorf("failed to authenticate with SMTP server: %w", err)
	}

	return nil
}
