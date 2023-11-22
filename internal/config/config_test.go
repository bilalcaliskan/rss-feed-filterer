//go:build unit

package config

import (
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/assert"
)

// TestReadConfig function tests if ReadConfig function running properly
func TestReadConfig(t *testing.T) {
	testCases := []struct {
		name       string
		path       string
		envs       map[string]string
		shouldPass bool
	}{
		{
			"valid config path",
			"../../test/config.yaml",
			map[string]string{
				"AWS_ACCESS_KEY": "testAccessKey",
				"AWS_SECRET_KEY": "testSecretKey",
			},
			true,
		},
		{
			"second valid config path",
			"../../test/config_verbose.yaml",
			map[string]string{
				"AWS_ACCESS_KEY": "testAccessKey",
				"AWS_SECRET_KEY": "testSecretKey",
			},
			true,
		},
		{
			"invalid config path",
			"../../test/invalid-config.yaml",
			map[string]string{},
			false,
		},
	}

	// run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("starting case %s", tc.name)

			for key, value := range tc.envs {
				err := os.Setenv(key, value)
				assert.Nil(t, err)
			}

			defer func() {
				for key := range tc.envs {
					err := os.Unsetenv(key)
					assert.Nil(t, err)
				}
			}()

			c, err := ReadConfig(&cobra.Command{}, tc.path)
			if tc.shouldPass {
				assert.Nil(t, err)
				assert.NotNil(t, c)
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, c)
			}

		})
	}
}
