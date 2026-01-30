package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
	texttemplate "text/template"

	"github.com/rs/zerolog/log"
)

// SMTPSender implements Sender interface using SMTP
type SMTPSender struct {
	config    Config
	templates map[string]Template
}

// NewSMTPSender creates a new SMTP email sender
func NewSMTPSender(config Config) *SMTPSender {
	return &SMTPSender{
		config:    config,
		templates: DefaultTemplates,
	}
}

// SetTemplates allows overriding default templates
func (s *SMTPSender) SetTemplates(templates map[string]Template) {
	s.templates = templates
}

// Send sends an email via SMTP
func (s *SMTPSender) Send(ctx context.Context, email *Email) error {
	// Use default from if not specified
	from := email.From
	if from == "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromAddress)
	}

	// Build message
	var msg bytes.Buffer

	// Headers
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ", ")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))

	if email.ReplyTo != "" {
		msg.WriteString(fmt.Sprintf("Reply-To: %s\r\n", email.ReplyTo))
	}

	// Custom headers
	for key, value := range email.Headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// Content type
	if email.HTMLBody != "" {
		// Multipart message with both text and HTML
		boundary := "==BOUNDARY=="
		msg.WriteString("MIME-Version: 1.0\r\n")
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary))
		msg.WriteString("\r\n")

		// Plain text part
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
		msg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
		msg.WriteString("\r\n")
		msg.WriteString(email.Body)
		msg.WriteString("\r\n")

		// HTML part
		msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		msg.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
		msg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
		msg.WriteString("\r\n")
		msg.WriteString(email.HTMLBody)
		msg.WriteString("\r\n")

		msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else {
		// Plain text only
		msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
		msg.WriteString("\r\n")
		msg.WriteString(email.Body)
	}

	// Connect and send
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)

	var auth smtp.Auth
	if s.config.SMTPUser != "" {
		auth = smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)
	}

	if s.config.SMTPUseTLS {
		return s.sendTLS(addr, auth, s.config.FromAddress, email.To, msg.Bytes())
	}

	return smtp.SendMail(addr, auth, s.config.FromAddress, email.To, msg.Bytes())
}

// sendTLS sends email using TLS connection
func (s *SMTPSender) sendTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// Connect with TLS
	tlsConfig := &tls.Config{
		ServerName: s.config.SMTPHost,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	// Authenticate if credentials provided
	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Set sender
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// Send message body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data connection: %w", err)
	}

	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close data connection: %w", err)
	}

	return client.Quit()
}

// SendTemplate sends an email using a named template
func (s *SMTPSender) SendTemplate(ctx context.Context, to []string, templateName string, data map[string]interface{}) error {
	tmpl, ok := s.templates[templateName]
	if !ok {
		return fmt.Errorf("template not found: %s", templateName)
	}

	// Parse and execute subject template
	subjectTmpl, err := texttemplate.New("subject").Parse(tmpl.Subject)
	if err != nil {
		return fmt.Errorf("failed to parse subject template: %w", err)
	}
	var subjectBuf bytes.Buffer
	if err := subjectTmpl.Execute(&subjectBuf, data); err != nil {
		return fmt.Errorf("failed to execute subject template: %w", err)
	}

	// Parse and execute body template
	bodyTmpl, err := texttemplate.New("body").Parse(tmpl.Body)
	if err != nil {
		return fmt.Errorf("failed to parse body template: %w", err)
	}
	var bodyBuf bytes.Buffer
	if err := bodyTmpl.Execute(&bodyBuf, data); err != nil {
		return fmt.Errorf("failed to execute body template: %w", err)
	}

	// Parse and execute HTML template
	var htmlBuf bytes.Buffer
	if tmpl.HTML != "" {
		htmlTmpl, err := template.New("html").Parse(tmpl.HTML)
		if err != nil {
			return fmt.Errorf("failed to parse HTML template: %w", err)
		}
		if err := htmlTmpl.Execute(&htmlBuf, data); err != nil {
			return fmt.Errorf("failed to execute HTML template: %w", err)
		}
	}

	email := &Email{
		To:       to,
		Subject:  subjectBuf.String(),
		Body:     bodyBuf.String(),
		HTMLBody: htmlBuf.String(),
	}

	log.Debug().
		Strs("to", to).
		Str("template", templateName).
		Msg("Sending template email")

	return s.Send(ctx, email)
}
