package filterer

import (
	"context"
	"fmt"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/feed"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/logging"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/storage/aws"
)

// create a channel to act as a semaphore
var sem = make(chan struct{}, 5) // allow up to 5 concurrent access

func Filter(ctx context.Context) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	client, err := aws.CreateClient(cfg.AccessKey, cfg.SecretKey, cfg.Region)
	if err != nil {
		return err
	}

	if !aws.IsBucketExists(client, cfg.BucketName) {
		return fmt.Errorf("bucket %s not found", cfg.BucketName)
	}

	for _, repo := range cfg.Repositories {
		go func(repo config.Repository) {
			checker := feed.NewReleaseChecker(client, repo, cfg.BucketName, logging.GetLogger())

			// acquire the semaphore
			sem <- struct{}{}
			checker.CheckGithubReleases(ctx, sem)
		}(repo)
	}

	<-ctx.Done()

	return nil
}
