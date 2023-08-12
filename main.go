package main

import (
	"context"
	"os"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/feed"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/logging"
)

func main() {
	logger := logging.GetLogger()
	cfg, err := config.ReadConfig()
	if err != nil {
		os.Exit(1)
	}
	//cfg := config.GetConfig()

	if err := feed.Filter(context.Background(), cfg); err != nil {
		logger.Error().Msg(err.Error())
		os.Exit(1)
	}
}
