package root

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/bilalcaliskan/rss-feed-filterer/cmd/root/options"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce/slack"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/feed"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/logging"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/storage/aws"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/version"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	ver     = version.Get()
	logger  zerolog.Logger
	opts    *options.RootOptions
	RootCmd = &cobra.Command{
		Use:           "rss-feed-filterer",
		Short:         "A tool designed to efficiently monitor, filter, and notify users about new releases in software projects based on their RSS feeds",
		Version:       ver.GitVersion,
		SilenceUsage:  true,
		SilenceErrors: true,
		Long: `RSS Feed Filterer is a sophisticated tool designed to efficiently monitor, filter, and notify users about new releases
in software projects based on their RSS feeds. It seamlessly integrates with AWS S3 to persist release data and provides a comprehensive
mechanism to track multiple project releases.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flag("config-file").Value.String() == "" {
				err := errors.New("required flag 'config-file' not set")
				logger.Error().Msg(err.Error())
				return err
			}

			if opts.VerboseLog {
				logging.EnableDebugLogging()
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.ReadConfig(opts.ConfigFilePath)
			if err != nil {
				logger.Error().Err(err).Msg("failed to read config")
				return err
			}

			client, err := aws.CreateClient(cfg.AccessKey, cfg.SecretKey, cfg.Region)
			if err != nil {
				logger.Error().Err(err).Msg("failed to create s3 client")
				return err
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
				return err
			}

			return nil
		},
	}
)

func init() {
	logger = logging.GetLogger()
	opts = options.GetRootOptions()
	opts.InitFlags(RootCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() error {
	return RootCmd.Execute()
}
