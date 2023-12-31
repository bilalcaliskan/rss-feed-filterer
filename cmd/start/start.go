package start

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bilalcaliskan/rss-feed-filterer/cmd/root/options"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/feed"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/storage/aws"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	ctx      context.Context
	cancel   context.CancelFunc
	StartCmd = &cobra.Command{
		Use:           "start",
		Short:         "starts the main process by reading the config file",
		SilenceUsage:  false,
		SilenceErrors: true,
		Example: `# clean the desired files on target bucket
s3-manager clean --min-size-mb=1 --max-size-mb=1000 --keep-last-n-files=2 --sort-by=lastModificationDate --order=ascending
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the values from the context
			cfg := cmd.Context().Value(options.ConfigKey{}).(*config.Config)
			client := cmd.Context().Value(options.S3ClientKey{}).(aws.S3ClientAPI)
			announcers := cmd.Context().Value(options.AnnouncerKey{}).([]announce.Announcer)
			logger := cmd.Context().Value(options.LoggerKey{}).(zerolog.Logger)

			// Listen for interrupt signals and cancel the context
			go func() {
				sigChan := make(chan os.Signal, 1)
				signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
				<-sigChan
				logger.Info().Msg("interrupt signal received., shutting down...")
				cancel()
			}()

			// Start the filtering process
			if err := feed.Filter(ctx, cfg, client, announcers); err != nil {
				logger.Error().Err(err).Msg("filtering process failed, shutting down...")
				return err
			}

			return nil
		},
	}
)

func init() {
	// Create a cancellable context
	ctx, cancel = context.WithCancel(context.Background())
}
