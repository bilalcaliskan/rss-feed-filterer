package main

import (
	"context"
	"os"

	"github.com/bilalcaliskan/rss-feed-filterer/internal/filterer"
)

func main() {
	// use below approach on running tests
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//defer cancel()

	if err := filterer.Filter(context.Background()); err != nil {
		os.Exit(1)
	}
}
