package filterer

import (
	"context"
	"time"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/feed"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/logging"
)

// create a channel to act as a semaphore
var sem = make(chan struct{}, 5) // allow up to 5 concurrent access

func Filter(ctx context.Context) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	logger := logging.GetLogger()

	logger.Info().
		Any("notification.slack", cfg.Notification.Slack).
		Msg("")

	logger.Info().
		Any("storage.s3", cfg.Storage.S3).
		Msg("")

	for _, repo := range cfg.Repositories {
		go func(repo config.Repository) {
			sem <- struct{}{}                                                                              // acquire the semaphore
			feed.CheckGithubReleases(ctx, sem, repo, time.Duration(repo.CheckIntervalMinutes)*time.Minute) // Start the goroutine to check GitHub releases
		}(repo)
	}

	<-ctx.Done()

	return nil
}
