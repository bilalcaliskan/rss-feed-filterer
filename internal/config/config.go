package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func ReadConfig(path string) (conf *Config, err error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "an error occurred while reading config file")
	}

	if err := yaml.Unmarshal(file, &conf); err != nil {
		return nil, errors.Wrap(err, "an error occurred while unmarshaling config file")
	}

	if err := conf.Storage.SetAccessCredentialsFromEnv(); err != nil {
		return nil, errors.Wrap(err, "an error occurred while setting credentials with env variables")
	}

	return conf, nil
}
