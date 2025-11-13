package flag

import (
	"io"
	"net"

	"github.com/gi8lino/randomapi/internal/logging"
	"github.com/gi8lino/randomapi/internal/server"

	"github.com/containeroo/tinyflags"
)

// Config holds all application configuration.
type Config struct {
	ListenAddr  string            // HTTP bind address (e.g. ":8080")
	LogFormat   logging.LogFormat // Log output format (text or json)
	RoutePrefix string            // Canonical path prefix ("" or "/random-api")
	DataPath    string            // Path to JSON file with elements
}

// ParseArgs parses CLI args into Config.
func ParseArgs(version string, args []string, out io.Writer) (Config, error) {
	var cfg Config
	tf := tinyflags.NewFlagSet("random-api", tinyflags.ContinueOnError)
	tf.Version(version)
	tf.EnvPrefix("RANDOMAPI")
	tf.SetOutput(out)

	// Server
	tf.StringVar(&cfg.RoutePrefix, "route-prefix", "", "Path prefix to mount the app (e.g., /random-api). Empty = root.").
		Finalize(func(input string) string {
			return server.NormalizeRoutePrefix(input) // canonical "" or "/random-api"
		}).
		Placeholder("PATH").
		Value()

	listenAddr := tf.TCPAddr("listen-address", &net.TCPAddr{IP: nil, Port: 8080}, "HTTP server listen address").
		Placeholder("ADDR:PORT").
		Value()

	dataPath := tf.String("data-path", "/app/data.json", "Path to JSON file with elements.").
		Placeholder("PATH").
		Value()

	// Logging
	logFormat := tf.String("log-format", "text", "Log format").
		Choices("text", "json").
		Short("l").
		Value()

	// Parse
	if err := tf.Parse(args); err != nil {
		return Config{}, err
	}

	// Post-parse
	cfg.LogFormat = logging.LogFormat(*logFormat)
	cfg.ListenAddr = (*listenAddr).String()
	cfg.DataPath = *dataPath

	return cfg, nil
}
