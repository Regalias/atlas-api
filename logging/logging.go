package logging

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// New configures the zerolog global options such as log level
// returns a zerolog.Logger
func New(level string, appname string, useConsole bool) (*zerolog.Logger, error) {
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

	var stream io.Writer
	if useConsole {
		// FOR DEBUG
		stream = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	} else {
		stream = os.Stdout
	}

	appLogger := zerolog.New(stream).With().
		Timestamp().
		Str("svc", appname).
		Str("host", host).
		Logger()

	return &appLogger, nil
}
