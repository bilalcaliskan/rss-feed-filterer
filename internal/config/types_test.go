//go:build unit

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_SetAccessCredentialsFromEnv(t *testing.T) {
	// Set environment variables
	err := os.Setenv("STORAGE_AWS_ACCESS_KEY", "testAccessKey")
	assert.Nil(t, err)
	err = os.Setenv("STORAGE_AWS_SECRET_KEY", "testSecretKey")
	assert.Nil(t, err)

	defer func(t *testing.T) {
		err := os.Unsetenv("STORAGE_AWS_SECRET_KEY")
		assert.Nil(t, err)

		err = os.Unsetenv("STORAGE_AWS_SECRET_KEY")
		assert.Nil(t, err)
	}(t)

	storage := &Storage{}
	if err := storage.SetAccessCredentialsFromEnv("aws"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if storage.S3.AccessKey != "testAccessKey" {
		t.Errorf("expected %s, got %s", "testAccessKey", storage.S3.AccessKey)
	}
}
