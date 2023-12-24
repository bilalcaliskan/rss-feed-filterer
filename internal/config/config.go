package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ReadConfig(cmd *cobra.Command, path string) (conf *Config, err error) {
	// set config file name and path
	viper.SetConfigName(path)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "an error occurred while reading config file")
	}

	// unmarshal config file into Config struct
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, errors.Wrap(err, "an error occurred while unmarshaling config file")
	}

	// if a flag is set, it overrides the value in config file
	if cmd.Flags().Changed("verbose") {
		verbose, _ := cmd.Flags().GetBool("verbose")
		conf.Verbose = verbose
	}

	// s3 access credentials can also be set from env variables so we check them here
	if err := conf.Storage.SetAccessCredentialsFromEnv(conf.Storage.Provider); err != nil {
		return nil, errors.Wrap(err, "an error occurred while setting credentials with env variables "+
			"for storage service")
	}

	// email access credentials can also be set from env variables so we check them here
	if err := conf.Email.SetAccessCredentialsFromEnv(conf.Email.Provider); err != nil {
		return nil, errors.Wrap(err, "an error occurred while setting credentials with env variables "+
			"for email service")
	}

	return conf, nil
}
