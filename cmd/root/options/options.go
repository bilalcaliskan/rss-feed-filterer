package options

import "github.com/spf13/cobra"

var rootOptions = &RootOptions{}

// RootOptions contains frequent command line and application options.
type RootOptions struct {
	// ConfigFilePath is the path for the config file to be used
	ConfigFilePath string
	// VerboseLog is the verbosity of the logging library
	VerboseLog bool
}

func (opts *RootOptions) InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opts.ConfigFilePath, "config-file", "c", "",
		"path for the config file to be used")
	cmd.Flags().BoolVarP(&opts.VerboseLog, "verbose", "", false,
		"verbose output of the logging library as 'debug' (default false)")

	_ = cmd.MarkFlagRequired("config-file")
}

// GetRootOptions returns the pointer of RootOptions
func GetRootOptions() *RootOptions {
	return rootOptions
}
