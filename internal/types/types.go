package types

import "time"

type Release struct {
	ProjectName string     `json:"projectName"`
	Version     string     `json:"version"`
	PublishedAt *time.Time `json:"publishedAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	Url         string     `json:"url"`
}
