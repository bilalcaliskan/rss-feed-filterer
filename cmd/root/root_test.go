//go:build e2e

package root

import (
	"testing"

	"github.com/bilalcaliskan/rss-feed-filterer/cmd/root/options"
	"github.com/bilalcaliskan/rss-feed-filterer/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	cases := []struct {
		caseName   string
		args       []string
		shouldPass bool
	}{
		{
			"Execute",
			[]string{"--verbose", "--config-file=../../test/config.yaml"},
			true,
		},
		{
			"Slack enabled config",
			[]string{"--config-file=../../test/config_slack_enabled.yaml"},
			true,
		},
		{
			"Email enabled config",
			[]string{"--config-file=../../test/config_email_enabled.yaml"},
			true,
		},
		{
			"Empty config path",
			[]string{"--verbose"},
			false,
		},
		{
			"Wrong config path",
			[]string{"--config-file=../../test/config_slack_enabled.yamllll"},
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

		rootCmd.SetArgs(tc.args)

		err := rootCmd.Execute()
		if !tc.shouldPass {
			assert.NotNil(t, err)
			continue
		}

		utils.SleepSeconds(10)
		assert.Nil(t, err)

		options.GetRootOptions().SetZeroValues()
	}
}

func TestOuterExecute(t *testing.T) {
	err := rootCmd.PersistentFlags().Set("verbose", "true")
	assert.Nil(t, err)

	_ = Execute()
}
