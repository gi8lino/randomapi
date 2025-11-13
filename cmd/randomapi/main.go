package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gi8lino/randomapi/internal/app"
)

var (
	Version string = "dev"
	Commit  string = "none"
)

func main() {
	// Create a root context
	ctx := context.Background()

	if err := app.Run(ctx, Version, os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
