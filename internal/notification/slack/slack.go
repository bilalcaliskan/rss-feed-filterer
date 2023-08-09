package slack

import (
	"fmt"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	api "github.com/slack-go/slack"
)

type Message struct {
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Fallback   string `json:"fallback"`
	Color      string `json:"color"`
	Pretext    string `json:"pretext"`
	AuthorName string `json:"author_name"`
	AuthorLink string `json:"author_link"`
	Title      string `json:"title"`
	Text       string `json:"text"`
	ThumbURL   string `json:"thumb_url"`
}

func SendNotification(projectName, version, url string) error {
	msg := api.WebhookMessage{
		Attachments: []api.Attachment{},
		Username:    "GoReleaser",
		IconURL:     "https://github.com/goreleaser/goreleaser/raw/939f2b002b29d2c8df6efd2d1f1d0b85c4ac5ee0/www/docs/static/logo.png",
		Text:        fmt.Sprintf("%s %s is out! Check it out at %s", projectName, version, url),
	}

	return api.PostWebhook(config.GetConfig().Notification.WebhookUrl, &msg)
}
