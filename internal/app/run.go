package app

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/gi8lino/randomapi/internal/data"
	"github.com/gi8lino/randomapi/internal/flag"
	"github.com/gi8lino/randomapi/internal/logging"
	"github.com/gi8lino/randomapi/internal/routes"

	"github.com/containeroo/httpgrace/server"
	"github.com/containeroo/tinyflags"
)

// Run is the main function of the application.
func Run(ctx context.Context, version string, argv []string, stdOut, stdErr io.Writer) error {
	// Parse command-line flags
	flags, err := flag.ParseArgs(version, argv, stdOut)
	if err != nil {
		if tinyflags.IsHelpRequested(err) || tinyflags.IsVersionRequested(err) {
			_, _ = fmt.Fprint(stdOut, err.Error())
			return nil
		}
		_, _ = fmt.Fprint(stdErr, err.Error())
		return err
	}

	logger := logging.SetupLogger(flags.LogFormat, flags.Debug, stdOut)
	setupLog := logger.With("component", "setup")
	setupLog.Info(
		"Starting randomAPI",
		"version", version,
	)

	// Record any CLI overrides to aid debugging.
	if len(flags.OverriddenValues) > 0 {
		setupLog.Info(
			"cli overrides",
			"overrides", flags.OverriddenValues,
		)
	}

	// Load data elements from JSON file
	elements, err := data.LoadElements(flags.DataPath)
	if err != nil {
		setupLog.Error("load elements", "path", flags.DataPath, "err", err)
		return err
	}
	if len(elements) == 0 {
		setupLog.Error("no elements available after load", "path", flags.DataPath)
		return errors.New("no elements available")

	}
	setupLog.Debug("loaded elements", "count", len(elements))

	// HTTP server
	serverLog := logger.With("component", "server")
	router := routes.NewRouter(
		serverLog,
		flags.RoutePrefix,
		elements,
	)

	ctx, stop := server.SignalContext(ctx)
	defer stop()

	if err := server.Run(ctx, flags.ListenAddr, router, serverLog); err != nil {
		setupLog.Error("server run", "listen_address", flags.ListenAddr, "error", err)
		return err
	}

	return nil
}
