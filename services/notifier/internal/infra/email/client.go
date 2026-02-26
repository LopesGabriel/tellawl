package email

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"

	"github.com/lopesgabriel/tellawl/packages/logger"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Client sends emails via SMTP (e.g. Gmail).
type Client struct {
	smtpHost string
	smtpPort string
	from     string
	password string
	logger   *logger.AppLogger
	tracer   trace.Tracer
}

type NewClientParams struct {
	SMTPHost string
	SMTPPort string
	From     string
	Password string
	Logger   *logger.AppLogger
	Tracer   trace.Tracer
}

func NewClient(params NewClientParams) *Client {
	return &Client{
		smtpHost: params.SMTPHost,
		smtpPort: params.SMTPPort,
		from:     params.From,
		password: params.Password,
		logger:   params.Logger,
		tracer:   params.Tracer,
	}
}

// SendEmail sends an email to the given recipients.
func (c *Client) SendEmail(ctx context.Context, to []string, subject, body string) error {
	ctx, span := c.tracer.Start(ctx, "email.SendEmail")
	defer span.End()

	addr := fmt.Sprintf("%s:%s", c.smtpHost, c.smtpPort)
	auth := smtp.PlainAuth("", c.from, c.password, c.smtpHost)

	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n",
		c.from,
		strings.Join(to, ", "),
		subject,
	)
	msg := []byte(headers + body)

	err := smtp.SendMail(addr, auth, c.from, to, msg)
	if err != nil {
		c.logger.Error(ctx, "Failed to send email",
			slog.Any("error", err),
			slog.String("to", strings.Join(to, ", ")),
			slog.String("subject", subject),
		)
		span.SetStatus(codes.Error, "failed to send email")
		span.RecordError(err)
		return fmt.Errorf("email: send failed: %w", err)
	}

	c.logger.Info(ctx, "Email sent successfully",
		slog.String("to", strings.Join(to, ", ")),
		slog.String("subject", subject),
	)
	span.SetStatus(codes.Ok, "success")
	return nil
}
