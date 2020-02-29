package apiserver

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// initLogger configures the zerolog global options such as log level
// returns a zerolog.Logger
func initLogger(level string) (*zerolog.Logger, error) {
	switch strings.ToLower(level) {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Grab hostname for logging field
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	appLogger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("svc", "linker-api").
		Str("host", host).
		Logger()

	return &appLogger, nil
}
