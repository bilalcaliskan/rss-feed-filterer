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
	maxRetries         = 3
	defaultSemverRegex = `/(v?\d+\.\d+\.\d+)$`
	releaseFileKey     = "releases.json"
)

func Filter(ctx context.Context, cfg *config.Config, client aws.S3ClientAPI, announcers []announce.Announcer) error {
	logger := logging.GetLogger()

	if !aws.IsBucketExists(client, cfg.BucketName) {
		err := fmt.Errorf("bucket %s not found", cfg.BucketName)
		logger.Error().Err(err).Str("bucketName", cfg.BucketName).Err(err).Msg("an error occurred while checking existence of bucket")
		return err
	}

	for _, repo := range cfg.Repositories {
		go func(repo config.Repository) {
			checker := NewReleaseChecker(client, repo, gofeed.NewParser(), cfg.BucketName, logging.GetLogger(), announcers)

			projectName, err := checker.extractProjectName()
			if err != nil {
				logger.Error().Err(err).Msg("failed to extract project name")
				return
			}

			checker.CheckGithubReleases(ctx, projectName)
		}(repo)
	}

	<-ctx.Done()

	return nil
}
