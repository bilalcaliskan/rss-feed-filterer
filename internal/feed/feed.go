package feed

import (
	"context"
	"fmt"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce/slack"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/logging"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/storage/aws"
)

const (
	maxRetries         = 3
	defaultSemverRegex = `/(v?\d+\.\d+\.\d+)$`
	releaseFileKey     = "releases.json"
)

// create a channel to act as a semaphore
var sem = make(chan struct{}, 5) // allow up to 5 concurrent access

func Filter(ctx context.Context, cfg config.Config) error {
	client, err := aws.CreateClient(cfg.AccessKey, cfg.SecretKey, cfg.Region)
	if err != nil {
		return err
	}

	if !aws.IsBucketExists(client, cfg.BucketName) {
		return fmt.Errorf("bucket %s not found", cfg.BucketName)
	}

	var announcer announce.Announcer

	if cfg.Announcer.Slack.Enabled {
		announcer = slack.NewSlackAnnouncer(cfg.Announcer.Slack.WebhookUrl, true)
	} else {
		announcer = &announce.NoopAnnouncer{}
	}

	for _, repo := range cfg.Repositories {
		go func(repo config.Repository) {
			checker := NewReleaseChecker(client, repo, cfg.BucketName, logging.GetLogger(), announcer)

			// acquire the semaphore
			sem <- struct{}{}
			checker.CheckGithubReleases(ctx, sem)
		}(repo)
	}

	<-ctx.Done()

	return nil
}
