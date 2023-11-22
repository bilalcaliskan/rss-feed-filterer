package feed

import (
	"context"
	"fmt"
	"sync"

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

// Filter function filters the feed and uploads the filtered feed to the bucket if there is a new release
func Filter(ctx context.Context, cfg *config.Config, client aws.S3ClientAPI, announcers []announce.Announcer) error {
	logger := logging.GetLogger()
	if !aws.IsBucketExists(client, cfg.BucketName) {
		err := fmt.Errorf("bucket %s not found", cfg.BucketName)
		logger.Error().Err(err).Str("bucketName", cfg.BucketName).Err(err).Msg("an error occurred while checking existence of bucket")
		return err
	}

	var wg sync.WaitGroup

	// iterate over repositories and start a goroutine for each repository to check for new releases
	for _, repo := range cfg.Repositories {
		wg.Add(1)
		go func(repo config.Repository) {
			defer wg.Done()
			checker := NewReleaseChecker(client, repo, gofeed.NewParser(), cfg.BucketName, logging.GetLogger(), announcers)

			projectName, err := checker.extractProjectName()
			if err != nil {
				logger.Error().Err(err).Msg("failed to extract project name")
				return
			}

			checker.CheckGithubReleases(ctx, projectName, cfg.OneShot)
		}(repo)
	}

	doneChan := make(chan struct{})
	// start a goroutine to wait for all other goroutines to finish their works
	go func() {
		wg.Wait()
		// notify the main goroutine that all other goroutines are finished their works
		close(doneChan)
	}()

	// block until we receive a notification from the doneChan
	<-doneChan
	logger.Info().Msg("all goroutines are finished their works, shutting down...")
	return nil
}
