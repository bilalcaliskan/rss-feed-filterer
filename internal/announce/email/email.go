package email

import (
	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
)

type EmailPayload struct {
	Subject string
	Content string
}

type EmailAnnouncer struct {
	Sender
	From string
	To   []string
	Cc   []string
	Bcc  []string
}

// Sender interface ensures that any specific email service (like SMTP, SES, etc.)
// can be integrated into the EmailAnnouncer.
type Sender interface {
	Send(to, cc, bcc []string, from, projectName, version, url string) error
}

func NewEmailAnnouncer(sender Sender, from string, to, cc, bcc []string) *EmailAnnouncer {
	return &EmailAnnouncer{
		sender,
		from,
		to,
		cc,
		bcc,
	}
}

func (e *EmailAnnouncer) Notify(payload *announce.AnnouncerPayload) error {
	//emailPayload, ok := payload.(*EmailPayload)
	//if !ok {
	//	return fmt.Errorf("invalid payload type, expected EmailPayload")
	//}

	return e.Send(e.To, e.Cc, e.Bcc, e.From, payload.ProjectName, payload.Version, payload.URL)
}

func (e *EmailAnnouncer) IsEnabled() bool {
	// This can be more dynamic, for now, I'm assuming if a sender and "to" address is present, it's enabled.
	return e.Sender != nil && len(e.To) > 0
}
