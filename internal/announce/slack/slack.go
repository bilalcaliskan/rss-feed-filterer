package slack

import (
	"fmt"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
)

type SlackAnnouncer struct {
	WebhookURL string
	//SlackPayload
	Enabled bool
}

type SlackPayload struct {
	ProjectName string
	Version     string
	URL         string
	IconUrl     string
	Username    string
}

func NewSlackAnnouncer(url string, enabled bool) *SlackAnnouncer {
	return &SlackAnnouncer{
		WebhookURL: url,
		Enabled:    enabled,
	}
}

func (sa *SlackAnnouncer) Notify(payload announce.AnnouncerPayload) error {
	slackPayload, ok := payload.(SlackPayload)
	if !ok {
		return fmt.Errorf("invalid payload type for SlackAnnouncer")
	}

	fmt.Println("inside notify")
	fmt.Println(slackPayload)
	return nil
	//slackPayload, ok := payload.(SlackPayload)
	//if !ok {
	//	return fmt.Errorf("invalid payload type for SlackAnnouncer")
	//}
	//
	//msg := api.WebhookMessage{
	//	Attachments: []api.Attachment{},
	//	Username:    slackPayload.Username,
	//	IconURL:     slackPayload.IconUrl,
	//	Text:        fmt.Sprintf("%s %s is out! Check it out at %s", slackPayload.ProjectName, slackPayload.Version, slackPayload.URL),
	//}
	//
	//return api.PostWebhook(sa.WebhookURL, &msg)
}

func (sa *SlackAnnouncer) IsEnabled() bool {
	return sa.Enabled
}

//func SendNotification(projectName, version, url string) error {
//	msg := api.WebhookMessage{
//		Attachments: []api.Attachment{},
//		Username:    "GoReleaser",
//		IconURL:     "https://github.com/goreleaser/goreleaser/raw/939f2b002b29d2c8df6efd2d1f1d0b85c4ac5ee0/www/docs/static/logo.png",
//		Text:        fmt.Sprintf("%s %s is out! Check it out at %s", projectName, version, url),
//	}
//
//	return api.PostWebhook(config.GetConfig().Notification.WebhookUrl, &msg)
//}
