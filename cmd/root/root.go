package root

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce/email"
	internalses "github.com/bilalcaliskan/rss-feed-filterer/internal/announce/email/ses"

	"github.com/bilalcaliskan/rss-feed-filterer/cmd/root/options"
	"github.com/bilalcaliskan/rss-feed-filterer/cmd/start"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/announce/slack"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/config"
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
	rootCmd = &cobra.Command{
		Use:           "rss-feed-filterer",
		Short:         "A tool designed to efficiently monitor, filter, and notify users about new releases in software projects based on their RSS feeds",
		Version:       ver.GitVersion,
		SilenceUsage:  false,
		SilenceErrors: true,
		Long: `RSS Feed Filterer is a sophisticated tool designed to efficiently monitor, filter, and notify users about new releases
in software projects based on their RSS feeds. It seamlessly integrates with AWS S3 to persist release data and provides a comprehensive
mechanism to track multiple project releases.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flag("config-file").Value.String() == "" {
				err := errors.New("required flag 'config-file' not set")
				logger.Error().Msg(err.Error())
				return err
			}

			logger = logging.GetLogger()
			logger.Info().Str("appVersion", ver.GitVersion).Str("goVersion", ver.GoVersion).Str("goOS", ver.GoOs).
				Str("goArch", ver.GoArch).Str("gitCommit", ver.GitCommit).Str("buildDate", ver.BuildDate).
				Msg("rss-feed-filterer is started!")

			cfg, err := config.ReadConfig(cmd, opts.ConfigFilePath)
			if err != nil {
				logger.Error().Err(err).Msg("failed to read config")
				return err
			}

			// check if verbose flag passed and enable debug logging
			if cfg.Verbose {
				logger = logging.WithVerbose()
				logger.Debug().Str("foo", "bar").Msg("this is a dummy log")
			}

			client, err := aws.CreateClient(cfg.AccessKey, cfg.SecretKey, cfg.Region)
			if err != nil {
				logger.Error().Err(err).Msg("failed to create s3 client")
				return err
			}

			var announcers []announce.Announcer
			if cfg.Announcer.Slack.Enabled {
				announcer := slack.NewSlackAnnouncer(cfg.Announcer.Slack.WebhookUrl, cfg.Announcer.Slack.Username, cfg.Announcer.Slack.IconUrl, &slack.SlackService{})
				announcers = append(announcers, announcer)
			}

			if cfg.Announcer.Email.Enabled {
				var sender email.Sender

				awsCfg, err := aws.CreateConfig(cfg.Email.AccessKey, cfg.Email.SecretKey, cfg.Email.Region)
				if err != nil {
					logger.Error().Err(err).Msg("failed to create aws config")
					return err
				}

				if cfg.Announcer.Email.Type == "ses" {
					sesClient := ses.NewFromConfig(awsCfg)
					sender = internalses.NewSESSender(sesClient)
				}
				//} else if cfg.Announcer.Email.Type == "smtp" {
				//	// TODO: implement smtp announcer
				//}

				announcer := email.NewEmailAnnouncer(sender, cfg.Announcer.Email.From, cfg.Announcer.Email.To, cfg.Announcer.Email.Cc, cfg.Announcer.Email.Bcc)
				announcers = append(announcers, announcer)
			}

			//else if cfg.Announcer.Email.Enabled {
			//	announcer = &announce.NoopAnnouncer{}
			//}

			cmd.SetContext(context.WithValue(cmd.Context(), options.ConfigKey{}, cfg))
			cmd.SetContext(context.WithValue(cmd.Context(), options.S3ClientKey{}, client))
			cmd.SetContext(context.WithValue(cmd.Context(), options.AnnouncerKey{}, announcers))
			cmd.SetContext(context.WithValue(cmd.Context(), options.LoggerKey{}, logger))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {
	logger = logging.GetLogger()
	opts = options.GetRootOptions()
	opts.InitFlags(rootCmd)

	rootCmd.AddCommand(start.StartCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}
