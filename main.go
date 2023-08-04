package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
)

const maxRetries = 3

// create a channel to act as a semaphore
var sem = make(chan struct{}, 5) // allow up to 5 concurrent access

func readConfig() Config {
	var config Config
	file, err := os.ReadFile("resources/config.yaml")
	if err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error parsing configuration file: %v", err)
	}
	return config
}

func checkGithubReleases(repo Repository, interval time.Duration) {
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

func main() {
	config := readConfig()
	for _, repo := range config.Repositories {
		go checkGithubReleases(repo, time.Duration(repo.CheckIntervalMinutes)*time.Minute) // Start the goroutine to check GitHub releases
	}

	// Block the main goroutine forever
	select {}
}
