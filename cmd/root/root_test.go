//go:build e2e

package root

import (
	"testing"
)

func TestExecuteFileCmd(t *testing.T) {
	var cases []struct {
		caseName   string
		args       []string
		shouldPass bool
	}

	for _, tc := range cases {
		t.Logf("starting case %s", tc.caseName)

	}
}
