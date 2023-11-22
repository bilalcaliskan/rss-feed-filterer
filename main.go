package main

import (
	"os"

	"github.com/bilalcaliskan/rss-feed-filterer/cmd/root"
)

func main() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
