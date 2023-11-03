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
	IconUrl    string
	Username   string
	Service    SlackAPI
}

func NewSlackAnnouncer(url, username, iconUrl string, service SlackAPI) *SlackAnnouncer {
	return &SlackAnnouncer{
		WebhookURL: url,
		Enabled:    true,
		Service:    service,
		Username:   username,
		IconUrl:    iconUrl,
	}
}

func (sa *SlackAnnouncer) Notify(payload *announce.AnnouncerPayload) error {
	msg := api.WebhookMessage{
		Attachments: []api.Attachment{},
		Username:    sa.Username,
		IconURL:     sa.IconUrl,
		Text:        fmt.Sprintf("%s %s is out! Check it out at %s", payload.ProjectName, payload.Version, payload.URL),
	}

	return sa.Service.PostWebhook(sa.WebhookURL, &msg)
}

func (sa *SlackAnnouncer) IsEnabled() bool {
	return sa.Enabled
}
