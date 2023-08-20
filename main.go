package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce/slack"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/storage/aws"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/feed"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/logging"
)

func main() {
	logger := logging.GetLogger()
	cfg, err := config.ReadConfig("resources/sample_config.yaml")
	if err != nil {
		logger.Error().Err(err).Msg("failed to read config")
		os.Exit(1)
	}

	client, err := aws.CreateClient(cfg.AccessKey, cfg.SecretKey, cfg.Region)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create s3 client")
		os.Exit(1)
	}

	var announcer announce.Announcer
	if cfg.Announcer.Slack.Enabled {
		announcer = slack.NewSlackAnnouncer(cfg.Announcer.Slack.WebhookUrl, true, &slack.SlackService{})
	} else {
		announcer = &announce.NoopAnnouncer{}
	}

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for interrupts to perform a graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-sigChan
		logger.Info().Msg("Interrupt received. Shutting down...")
		cancel()
	}()

	if err := feed.Filter(ctx, cfg, client, announcer); err != nil {
		logger.Error().Err(err).Msg("filtering process failed")
		os.Exit(1)
	}
}
