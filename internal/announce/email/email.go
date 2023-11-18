package email

import (
	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
)

// Sender interface ensures that any specific email service (like SMTP, SES, etc.)
// can be integrated into the EmailAnnouncer.
type Sender interface {
	Send(to, cc, bcc []string, from, projectName, version, url string) error
}

// EmailPayload is the payload that is sent to the email service.
// It contains the subject and content of the email.
type EmailPayload struct {
	Subject string
	Content string
}

// EmailAnnouncer is the announcer that sends the email. It contains the sender and the email addresses.
type EmailAnnouncer struct {
	Sender
	From string
	To   []string
	Cc   []string
	Bcc  []string
}

// NewEmailAnnouncer creates a new EmailAnnouncer. It requires a sender and the email addresses.
func NewEmailAnnouncer(sender Sender, from string, to, cc, bcc []string) *EmailAnnouncer {
	return &EmailAnnouncer{
		sender,
		from,
		to,
		cc,
		bcc,
	}
}

// Notify sends the email to the recipients.
func (e *EmailAnnouncer) Notify(payload *announce.AnnouncerPayload) error {
	return e.Send(e.To, e.Cc, e.Bcc, e.From, payload.ProjectName, payload.Version, payload.URL)
}

// IsEnabled checks if the EmailAnnouncer is enabled.
func (e *EmailAnnouncer) IsEnabled() bool {
	// This can be more dynamic, for now, I'm assuming if a sender and "to" address is present, it's enabled.
	return e.Sender != nil && len(e.To) > 0
}
