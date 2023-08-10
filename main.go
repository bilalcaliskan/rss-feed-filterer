package main

import (
	"context"
	"os"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/feed"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/storage/aws"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/logging"
)

func main() {
	// use below approach on running tests
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//defer cancel()

	logger := logging.GetLogger()
	cfg, err := config.ReadConfig()
	if err != nil {
		os.Exit(1)
	}
	//cfg := config.GetConfig()

	client, err := aws.CreateClient(cfg.AccessKey, cfg.SecretKey, cfg.Region)
	if err != nil {
		os.Exit(1)
	}

	if err := feed.Filter(context.Background(), cfg, client); err != nil {
		logger.Error().Msg(err.Error())
		os.Exit(1)
	}
}
