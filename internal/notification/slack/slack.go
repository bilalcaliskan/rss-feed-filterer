package slack

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SlackMessage struct {
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

func SendSlackNotification(webhookURL string, msg SlackMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
