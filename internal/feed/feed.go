package feed

import (
	"context"
	"fmt"

	"github.com/mmcdole/gofeed"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/logging"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/storage/aws"
)

const (
	maxRetries = 3
	//defaultRegex       = `^https://github\.com/[^/]+/[^/]+/releases/tag/(v?\d+\.\d+\.\d+)$`
	defaultSemverRegex = `/(v?\d+\.\d+\.\d+)$`
	releaseFileKey     = "releases.json"
)

// create a channel to act as a semaphore
var sem = make(chan struct{}, 5) // allow up to 5 concurrent access

func Filter(ctx context.Context, cfg config.Config, client aws.S3ClientAPI, announcer announce.Announcer) error {
	logger := logging.GetLogger()

	if !aws.IsBucketExists(client, cfg.BucketName) {
		err := fmt.Errorf("bucket %s not found", cfg.BucketName)
		logger.Error().Err(err).Str("bucketName", cfg.BucketName).Err(err).Msg("an error occurred while checking existence of bucket")
		return err
	}

	for _, repo := range cfg.Repositories {
		go func(repo config.Repository) {
			checker := NewReleaseChecker(client, repo, gofeed.NewParser(), cfg.BucketName, logging.GetLogger(), announcer)

			projectName, err := checker.extractProjectName()
			if err != nil {
				logger.Error().Err(err).Msg("failed to extract project name")
				return
			}

			sem <- struct{}{} // acquire the semaphore
			checker.CheckGithubReleases(ctx, sem, projectName)
		}(repo)
	}

	<-ctx.Done()

	return nil
}
