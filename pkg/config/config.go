package config

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Config holds application configuration values loaded from environment
// variables such as Plex endpoint, token, update interval and logging.
type Config struct {
	PlexUrl        string
	PlexToken      string
	UpdateInterval time.Duration
	MaxRetries     int
	LogFormat      string // pretty | json
}

// CreateConfig constructs a Config by reading environment variables and
// applying defaults where necessary. It will fatal-exit if required values are
// missing.
func CreateConfig() Config {
	return Config{
		PlexUrl:        Env("PLEX_URL", Required()),
		PlexToken:      Env("PLEX_TOKEN", Required()),
		UpdateInterval: MustParseDuration(Env("UPDATE_INTERVAL", Default("1m"))),
		MaxRetries:     MustParseInt(Env("MAX_RETRIES", Default("3"))),
		LogFormat:      Env("LOG_FORMAT", Default("pretty")),
	}
}

// CreateLogger returns a configured zerolog.Logger according to `format`.
// Supported formats are "pretty" (human-friendly) and "json" (structured).
// The function exits the process for an unknown format.
func CreateLogger(format string) zerolog.Logger {
	var stderrWriter io.Writer

	switch format {
	case "pretty":
		stderrWriter = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC822}
	case "json":
		stderrWriter = os.Stderr
	default:
		fmt.Fprintf(os.Stderr, "Bad log format %s\n", format)
		os.Exit(1)
	}

	return zerolog.New(stderrWriter).With().Timestamp().Logger()
}
