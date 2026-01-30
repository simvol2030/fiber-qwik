package email

import (
	"context"
)

// Email represents an email message
type Email struct {
	To          []string          // Recipients
	Subject     string            // Subject line
	Body        string            // Plain text body
	HTMLBody    string            // HTML body (optional)
	From        string            // Sender (optional, uses default if empty)
	ReplyTo     string            // Reply-to address (optional)
	Attachments []Attachment      // File attachments (optional)
	Headers     map[string]string // Custom headers (optional)
}

// Attachment represents an email attachment
type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

// Sender defines the interface for sending emails
// Implement this interface to add support for different email providers
type Sender interface {
	// Send sends an email
	Send(ctx context.Context, email *Email) error

	// SendTemplate sends an email using a named template
	SendTemplate(ctx context.Context, to []string, templateName string, data map[string]interface{}) error
}

// Config holds email service configuration
type Config struct {
	// Provider type: "smtp", "sendgrid", "ses", "mailgun"
	Provider string

	// Default sender
	FromName    string
	FromAddress string

	// SMTP settings
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPUseTLS   bool

	// SendGrid settings
	SendGridAPIKey string

	// AWS SES settings
	SESRegion    string
	SESAccessKey string
	SESSecretKey string
}

// Template represents an email template
type Template struct {
	Subject string
	Body    string
	HTML    string
}

// Common template names
const (
	TemplatePasswordReset    = "password_reset"
	TemplateEmailVerify      = "email_verify"
	TemplateWelcome          = "welcome"
	TemplatePasswordChanged  = "password_changed"
)

// DefaultTemplates provides basic email templates
// In production, you'd load these from files or database
var DefaultTemplates = map[string]Template{
	TemplatePasswordReset: {
		Subject: "Reset Your Password",
		Body:    "Click the following link to reset your password: {{.ResetURL}}\n\nThis link expires in {{.ExpiresIn}}.\n\nIf you didn't request this, please ignore this email.",
		HTML: `
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #3b82f6;">Reset Your Password</h2>
		<p>Click the button below to reset your password:</p>
		<p style="margin: 30px 0;">
			<a href="{{.ResetURL}}" style="background-color: #3b82f6; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block;">
				Reset Password
			</a>
		</p>
		<p style="color: #666; font-size: 14px;">This link expires in {{.ExpiresIn}}.</p>
		<p style="color: #666; font-size: 14px;">If you didn't request this, please ignore this email.</p>
	</div>
</body>
</html>`,
	},
	TemplateEmailVerify: {
		Subject: "Verify Your Email",
		Body:    "Click the following link to verify your email: {{.VerifyURL}}\n\nThis link expires in {{.ExpiresIn}}.",
		HTML: `
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #3b82f6;">Verify Your Email</h2>
		<p>Click the button below to verify your email address:</p>
		<p style="margin: 30px 0;">
			<a href="{{.VerifyURL}}" style="background-color: #3b82f6; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block;">
				Verify Email
			</a>
		</p>
		<p style="color: #666; font-size: 14px;">This link expires in {{.ExpiresIn}}.</p>
	</div>
</body>
</html>`,
	},
	TemplateWelcome: {
		Subject: "Welcome to {{.AppName}}!",
		Body:    "Welcome to {{.AppName}}!\n\nYour account has been created successfully.\n\nEmail: {{.Email}}\n\nGet started: {{.DashboardURL}}",
		HTML: `
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #3b82f6;">Welcome to {{.AppName}}!</h2>
		<p>Your account has been created successfully.</p>
		<p><strong>Email:</strong> {{.Email}}</p>
		<p style="margin: 30px 0;">
			<a href="{{.DashboardURL}}" style="background-color: #3b82f6; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block;">
				Go to Dashboard
			</a>
		</p>
	</div>
</body>
</html>`,
	},
	TemplatePasswordChanged: {
		Subject: "Your Password Was Changed",
		Body:    "Your password was changed successfully.\n\nIf you didn't make this change, please contact support immediately.",
		HTML: `
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #3b82f6;">Password Changed</h2>
		<p>Your password was changed successfully.</p>
		<p style="color: #666; font-size: 14px;">If you didn't make this change, please contact support immediately.</p>
	</div>
</body>
</html>`,
	},
}
