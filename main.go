package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
)

type Config struct {
	Repositories []Repository `yaml:"repositories"`
}

type Repository struct {
	Name                 string `yaml:"name"`
	Description          string `yaml:"description,omitempty"`
	RSSURL               string `yaml:"rssURL"`
	CheckIntervalMinutes int    `yaml:"checkIntervalMinutes"`
}

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
	fp := gofeed.NewParser()

	checkFeed := func() {
		feed, err := fp.ParseURL(repo.RSSURL)
		if err != nil {
			log.Printf("Error fetching feed: %v", err)
			return
		}

		for _, item := range feed.Items {
			log.Printf("Title: %v\nLink: %v\nPublished: %v\n", item.Title, item.Link, item.Published)
		}
	}

	checkFeed() // Run once immediately

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
