package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ReadConfig(cmd *cobra.Command, path string) (conf *Config, err error) {
	viper.SetConfigName(path)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "an error occurred while reading config file")
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return nil, errors.Wrap(err, "an error occurred while unmarshaling config file")
	}

	if cmd.Flags().Changed("verbose") {
		verbose, _ := cmd.Flags().GetBool("verbose")
		conf.Verbose = verbose
	}

	//file, err := os.ReadFile(path)
	//if err != nil {
	//	return nil, errors.Wrap(err, "an error occurred while reading config file")
	//}
	//
	//if err := yaml.Unmarshal(file, &conf); err != nil {
	//	return nil, errors.Wrap(err, "an error occurred while unmarshaling config file")
	//}

	// s3 access credentials can also be set from env variables so we check them here
	if err := conf.Storage.SetAccessCredentialsFromEnv(conf.Storage.Provider); err != nil {
		return nil, errors.Wrap(err, "an error occurred while setting credentials with env variables for storage service")
	}

	// email access credentials can also be set from env variables so we check them here
	if err := conf.Email.SetAccessCredentialsFromEnv(conf.Email.Provider); err != nil {
		return nil, errors.Wrap(err, "an error occurred while setting credentials with env variables for email service")
	}

	return conf, nil
}
