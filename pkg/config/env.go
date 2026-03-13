package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func newDefaultConfig() config {
	return config{
		required:     false,
		defaultValue: "",
	}
}

// Env returns the value of the environment variable named by `name`.
// WriterOption values may be provided to mark the variable as required or to
// supply a default. If the variable is required and not set, Env logs a fatal
// error and exits.
func Env(name string, opts ...WriterOption) string {
	cfg := newDefaultConfig()
	for _, opt := range opts {
		opt.apply(&cfg)
	}

	env := strings.TrimSpace(os.Getenv(name))

	if env == "" && cfg.required {
		log.Fatal().Str("Env.Name", name).Err(fmt.Errorf("environment variable %s must be set", name)).Send()
	}

	if env == "" {
		env = cfg.defaultValue
	}

	return env
}

// WriterOption configures the behavior of Env (for example making a variable
// required or providing a default).
type WriterOption interface {
	apply(*config)
}

type optionFunc func(*config)

func (fn optionFunc) apply(c *config) { fn(c) }

type config struct {
	required     bool
	defaultValue string
}

// Required returns a WriterOption that marks an environment variable as
// required. When applied, Env will fatal-exit if the variable is not set.
func Required() WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.required = true
	})
}

// Default returns a WriterOption that supplies a default value used when the
// environment variable is not set.
func Default(defaultValue string) WriterOption {
	return optionFunc(func(cfg *config) {
		cfg.defaultValue = defaultValue
	})
}

// MustParseInt parses int from string or exits with code 1
func MustParseInt(str string) int {
	parsed, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("could not parse int from %s", str))
	}
	return int(parsed)
}

// MustParseDuration parses duration from string or exits with code 1
func MustParseDuration(str string) time.Duration {
	parsed, err := time.ParseDuration(str)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("could not parse duration from %s", str))
	}
	return parsed
}
