package slack

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	api "github.com/slack-go/slack"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
)

type SlackAPI interface {
	PostWebhook(url string, msg *api.WebhookMessage) error
}

type SlackService struct{}

func (r *SlackService) PostWebhook(url string, msg *api.WebhookMessage) error {
	return api.PostWebhook(url, msg)
}

type MockSlackService struct {
	mock.Mock
}

type SlackAnnouncer struct {
	WebhookURL string
	Enabled    bool
	Service    SlackAPI
}

type SlackPayload struct {
	ProjectName string
	Version     string
	URL         string
	IconUrl     string
	Username    string
}

func NewSlackAnnouncer(url string, enabled bool, service SlackAPI) *SlackAnnouncer {
	return &SlackAnnouncer{
		WebhookURL: url,
		Enabled:    enabled,
		Service:    service,
	}
}

func (sa *SlackAnnouncer) Notify(payload announce.AnnouncerPayload) error {
	slackPayload, ok := payload.(SlackPayload)
	if !ok {
		return fmt.Errorf("invalid payload type for SlackAnnouncer")
	}

	msg := api.WebhookMessage{
		Attachments: []api.Attachment{},
		Username:    slackPayload.Username,
		IconURL:     slackPayload.IconUrl,
		Text:        fmt.Sprintf("%s %s is out! Check it out at %s", slackPayload.ProjectName, slackPayload.Version, slackPayload.URL),
	}

	return sa.Service.PostWebhook(sa.WebhookURL, &msg)
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
