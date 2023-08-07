package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Repositories []Repository `yaml:"repositories"`
	Storage      `yaml:"storage"`
	Notification `yaml:"notification"`
}

type Repository struct {
	Name                 string `yaml:"name"`
	Description          string `yaml:"description"`
	RSSURL               string `yaml:"rssURL"`
	CheckIntervalMinutes int    `yaml:"checkIntervalMinutes"`
}

type Notification struct {
	Slack `yaml:"slack"`
}

type Slack struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookUrl string `yaml:"webhookUrl"`
}

type Storage struct {
	S3 `yaml:"s3"`
}

type S3 struct {
	Provider   string `yaml:"provider"`
	AccessKey  string `yaml:"accessKey"`
	SecretKey  string `yaml:"secretKey"`
	Region     string `yaml:"region"`
	BucketName string `yaml:"bucketName"`
}

func (s *Storage) SetAccessCredentialsFromEnv() error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("aws")
	if err := viper.BindEnv("access_key", "secret_key", "bucket_name", "region"); err != nil {
		return err
	}

	fields := map[string]*string{
		"access_key":  &s.AccessKey,
		"secret_key":  &s.SecretKey,
		"region":      &s.Region,
		"bucket_name": &s.BucketName,
	}

	for key, field := range fields {
		if val, ok := viper.Get(key).(string); ok && val != "" {
			*field = val
		}
	}

	return nil
}
