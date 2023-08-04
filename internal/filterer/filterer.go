package filterer

import (
	config2 "github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/feed"
	"time"
)

func Filter() {
	config := config2.ReadConfig()
	for _, repo := range config.Repositories {
		go feed.CheckGithubReleases(repo, time.Duration(repo.CheckIntervalMinutes)*time.Minute) // Start the goroutine to check GitHub releases
	}

	select {}
}
