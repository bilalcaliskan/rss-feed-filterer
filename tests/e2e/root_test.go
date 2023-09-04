//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"github.com/bilalcaliskan/rss-feed-filterer/cmd/root"
	"github.com/bilalcaliskan/rss-feed-filterer/cmd/root/options"
	"testing"
)

func TestExecuteRootCmd(t *testing.T) {
	rootOpts := options.GetRootOptions()
	ctx := context.Background()
	root.RootCmd.SetContext(ctx)

	cases := []struct {
		caseName   string
		args       []string
		shouldPass bool
	}{
		{
			"aldsfkaslfkj",
			[]string{},
			true,
		},
	}

	for _, tc := range cases {
		t.Logf("starting case %s\n", tc.caseName)

		fmt.Println(rootOpts)
	}
}
