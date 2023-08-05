package feed

import (
	"context"
	"time"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/logging"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog"
)

func init() {
	logger = logging.GetLogger()
}

const maxRetries = 3

var logger zerolog.Logger

func CheckGithubReleases(ctx context.Context, sem chan struct{}, repo config.Repository, interval time.Duration) {
	parser := gofeed.NewParser()

	checkFeed := func() {
		defer func() {
			<-sem
		}() // release the semaphore when done

		for retries := 0; retries < maxRetries; retries++ {
			logger.Info().Str("name", repo.Name).Msg("trying to fetch the feed")

			//feed, err := fp.ParseURL(repo.RSSURL)
			_, err := parser.ParseURL(repo.RSSURL)
			if err != nil {
				logger.Warn().
					Str("error", err.Error()).
					Str("url", repo.RSSURL).
					Msg("an error occurred while fetching feed, retrying...")
				time.Sleep(time.Second * 5)
				continue
			}

			//for _, item := range feed.Items {
			//	logger.Info().Str("name", repo.Name).Msg(item.Title)
			//}
			logger.Info().Str("name", repo.Name).Msg("fetched releases")

			break
		}
	}

	// Run once immediately since ticker will wait for first attempt
	checkFeed()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			checkFeed()
		}
	}
}
