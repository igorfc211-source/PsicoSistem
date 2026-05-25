package email

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"api-on/internal/shared/config"
)

type Message struct {
	To       string
	Subject  string
	TextBody string
	HTMLBody string
}

type Sender interface {
	Send(ctx context.Context, message Message) error
}

type NoopSender struct{}

func (NoopSender) Send(_ context.Context, _ Message) error {
	return nil
}

type SMTPSender struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSender(cfg *config.Config) Sender {
	if cfg == nil || cfg.SMTPHost == "" || cfg.SMTPFrom == "" {
		if cfg != nil && strings.ToLower(strings.TrimSpace(cfg.AppEnv)) != "production" {
			return LogSender{}
		}
		return NoopSender{}
	}

	return &SMTPSender{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		username: cfg.SMTPUsername,
		password: cfg.SMTPPassword,
		from:     cfg.SMTPFrom,
	}
}

type LogSender struct{}

func (LogSender) Send(_ context.Context, message Message) error {
	log.Printf(
		"SMTP not configured; development email to=%s subject=%q body=%s",
		message.To,
		message.Subject,
		firstNonEmpty(message.TextBody, message.HTMLBody),
	)
	return nil
}

func (s *SMTPSender) Send(ctx context.Context, message Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	to := strings.TrimSpace(message.To)
	if to == "" {
		return fmt.Errorf("email recipient is required")
	}

	addr := s.host + ":" + s.port
	var auth smtp.Auth
	if s.username != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	return smtp.SendMail(addr, auth, s.from, []string{to}, buildPayload(s.from, message))
}

func buildPayload(from string, message Message) []byte {
	body := message.TextBody
	contentType := "text/plain"
	if message.HTMLBody != "" {
		body = message.HTMLBody
		contentType = "text/html"
	}

	headers := []string{
		"From: " + from,
		"To: " + message.To,
		"Subject: " + message.Subject,
		"MIME-Version: 1.0",
		"Content-Type: " + contentType + "; charset=UTF-8",
	}

	return []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + body)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
