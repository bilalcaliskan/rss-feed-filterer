package feed

import (
	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/mmcdole/gofeed"
	"log"
	"time"
)

const maxRetries = 3

// create a channel to act as a semaphore
var sem = make(chan struct{}, 5) // allow up to 5 concurrent access

func CheckGithubReleases(repo config.Repository, interval time.Duration) {
	sem <- struct{}{}        // acquire the semaphore
	defer func() { <-sem }() // release the semaphore when done

	fp := gofeed.NewParser()

	checkFeed := func() {
		for retries := 0; retries < maxRetries; retries++ {
			feed, err := fp.ParseURL(repo.RSSURL)
			if err != nil {
				log.Printf("Error fetching feed: %v, retrying...", err)
				time.Sleep(time.Second * time.Duration(2<<retries))
				continue
			}

			for _, item := range feed.Items {
				log.Printf("Title: %v\nLink: %v\nPublished: %v\n", item.Title, item.Link, item.Published)
			}

			break
		}
	}

	// Run once immediately since ticker will wait for first attempt
	checkFeed()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		checkFeed()
	}
}
