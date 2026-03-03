package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gi8lino/randomapi/internal/data"
	"github.com/gi8lino/randomapi/internal/flag"
	"github.com/gi8lino/randomapi/internal/logging"
	"github.com/gi8lino/randomapi/internal/routes"

	"github.com/containeroo/httpgrace/server"
	"github.com/containeroo/tinyflags"
)

// Run is the main function of the application.
func Run(ctx context.Context, version string, argv []string, output io.Writer) error {
	// Create a new context that listens for interrupt signals
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Parse command-line flags
	flags, err := flag.ParseArgs(version, argv, output)

	logger := logging.SetupLogger(flags.LogFormat, flags.Debug, output)

	logger.Info(
		"Starting randomAPI",
		"version", version,
	)

	if err != nil {
		if tinyflags.IsHelpRequested(err) || tinyflags.IsVersionRequested(err) {
			_, _ = fmt.Fprint(output, err.Error())
			return nil
		}
		logger.Error("Configuration error", "error", err)
		return err
	}

	// Record any CLI overrides to aid debugging.
	if len(flags.OverriddenValues) > 0 {
		logger.Info(
			"cli overrides",
			"overrides", flags.OverriddenValues,
		)
	}

	// Load data elements from JSON file
	elements, err := data.LoadElements(flags.DataPath)
	if err != nil {
		logger.Error("load elements", "path", flags.DataPath, "err", err)
		return fmt.Errorf("load elements: %w", err)
	}
	if len(elements) == 0 {
		err := errors.New("no elements available")
		logger.Error("no elements available after load", "path", flags.DataPath)
		return err
	}

	logger.Debug("loaded elements", "count", len(elements))

	// Create server and run forever
	router := routes.NewRouter(logger, flags.RoutePrefix, elements)
	if err := server.Run(ctx, flags.ListenAddr, router, logger); err != nil {
		logger.Error("server run", "listen_address", flags.ListenAddr, "err", err)
		return fmt.Errorf("server: %w", err)
	}

	// Small grace for prober loop to finish a tick if shutting down extremely fast
	time.Sleep(100 * time.Millisecond)
	return nil
}
