package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

func ReadConfig() Config {
	var config Config
	file, err := os.ReadFile("resources/sample_config.yaml")
	if err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error parsing configuration file: %v", err)
	}
	return config
}
