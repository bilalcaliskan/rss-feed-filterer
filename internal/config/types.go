package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Repositories []Repository `yaml:"repositories"`
	Storage      `yaml:"storage"`
	Announcer    `yaml:"announcer"`
	Type         string `yaml:"type"`
}

type Repository struct {
	Name                 string `yaml:"name"`
	Description          string `yaml:"description"`
	Url                  string `yaml:"url"`
	CheckIntervalMinutes int    `yaml:"checkIntervalMinutes"`
	//FeedType             string `yaml:"feedType"`
}

type Announcer struct {
	Slack `yaml:"slack"`
	Email `yaml:"email"`
}

type Email struct {
	Enabled bool     `yaml:"enabled"`
	From    string   `yaml:"from"`
	To      []string `yaml:"to"`
	Type    string   `yaml:"type"`
	Cc      []string `yaml:"cc"`
	Bcc     []string `yaml:"bcc"`
	Ses     `yaml:"ses"`
	Smtp    `yaml:"smtp"`
}

type Ses struct {
	Region    string `yaml:"region"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
}

type Smtp struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
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
	viper.SetEnvPrefix("storage_aws")
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

func (s *Ses) SetAccessCredentialsFromEnv() error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("email_aws")
	if err := viper.BindEnv("access_key", "secret_key", "region"); err != nil {
		return err
	}

	fields := map[string]*string{
		"access_key": &s.AccessKey,
		"secret_key": &s.SecretKey,
		"region":     &s.Region,
	}

	for key, field := range fields {
		if val, ok := viper.Get(key).(string); ok && val != "" {
			*field = val
		}
	}

	return nil
}
