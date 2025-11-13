package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gi8lino/randomapi/internal/data"
	"github.com/gi8lino/randomapi/internal/flag"
	"github.com/gi8lino/randomapi/internal/logging"
	"github.com/gi8lino/randomapi/internal/server"

	"github.com/containeroo/tinyflags"
)

// Run is the main function of the application.
func Run(ctx context.Context, version string, argv []string, output io.Writer) error {
	// Create a new context that listens for interrupt signals
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Parse command-line flags
	cfg, err := flag.ParseArgs(version, argv, output)
	if err != nil {
		if tinyflags.IsHelpRequested(err) || tinyflags.IsVersionRequested(err) {
			fmt.Fprint(output, err.Error()) // nolint:errcheck
			return nil
		}
		return fmt.Errorf("configuration error: %w", err)
	}

	// Setup logger
	logger := logging.SetupLogger(cfg.LogFormat, output)

	// Load data elements from JSON file
	elements, err := data.LoadElements(cfg.DataPath)
	if err != nil {
		return fmt.Errorf("load elements: %w", err)
	}

	// Create server and run forever
	router := server.NewRouter(logger, cfg.RoutePrefix, elements)
	if err := server.Run(ctx, cfg.ListenAddr, router, logger); err != nil {
		return fmt.Errorf("server: %w", err)
	}

	// Small grace for prober loop to finish a tick if shutting down extremely fast
	time.Sleep(100 * time.Millisecond)
	return nil
}
