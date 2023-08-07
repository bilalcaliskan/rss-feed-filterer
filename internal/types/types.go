package types

import "time"

type Release struct {
	Name        string     `json:"name"`
	PublishedAt *time.Time `json:"publishedAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	Url         string     `json:"url"`
	IsNotified  bool       `json:"isNotified"`
}
