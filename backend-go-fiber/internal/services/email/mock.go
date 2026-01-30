package email

import (
	"context"

	"github.com/rs/zerolog/log"
)

// MockSender implements Sender interface for development/testing
// Instead of sending emails, it logs them
type MockSender struct {
	config    Config
	templates map[string]Template
	SentMails []Email // Stores sent emails for testing
}

// NewMockSender creates a new mock email sender
func NewMockSender(config Config) *MockSender {
	return &MockSender{
		config:    config,
		templates: DefaultTemplates,
		SentMails: make([]Email, 0),
	}
}

// SetTemplates allows overriding default templates
func (s *MockSender) SetTemplates(templates map[string]Template) {
	s.templates = templates
}

// Send logs the email instead of actually sending it
func (s *MockSender) Send(ctx context.Context, email *Email) error {
	log.Info().
		Strs("to", email.To).
		Str("subject", email.Subject).
		Str("body_preview", truncate(email.Body, 100)).
		Msg("📧 [MOCK] Email sent")

	// Store for testing
	s.SentMails = append(s.SentMails, *email)

	return nil
}

// SendTemplate logs the template email
func (s *MockSender) SendTemplate(ctx context.Context, to []string, templateName string, data map[string]interface{}) error {
	log.Info().
		Strs("to", to).
		Str("template", templateName).
		Interface("data", data).
		Msg("📧 [MOCK] Template email sent")

	// Create email from template for storage
	tmpl, ok := s.templates[templateName]
	if ok {
		email := Email{
			To:       to,
			Subject:  tmpl.Subject,
			Body:     tmpl.Body,
			HTMLBody: tmpl.HTML,
		}
		s.SentMails = append(s.SentMails, email)
	}

	return nil
}

// GetLastEmail returns the last sent email (useful for testing)
func (s *MockSender) GetLastEmail() *Email {
	if len(s.SentMails) == 0 {
		return nil
	}
	return &s.SentMails[len(s.SentMails)-1]
}

// Clear clears all stored emails
func (s *MockSender) Clear() {
	s.SentMails = make([]Email, 0)
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}
