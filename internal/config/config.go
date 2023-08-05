package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func ReadConfig() (conf Config, err error) {
	file, err := os.ReadFile("resources/sample_config.yaml")
	if err != nil {
		return conf, errors.Wrap(err, "an error occurred while reading config file")
	}

	if err := yaml.Unmarshal(file, &conf); err != nil {
		return conf, errors.Wrap(err, "an error occurred while unmarshaling config file")
	}

	if err := conf.Storage.SetAccessCredentialsFromEnv(); err != nil {
		return conf, errors.Wrap(err, "an error occurred while setting credentials with env variables")
	}

	return conf, nil
}
