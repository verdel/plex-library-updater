package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/verdel/plex-library-updater/pkg/config"
	"github.com/verdel/plex-library-updater/pkg/plex"
)

func run(ctx context.Context, app config.Config) {
	ticker := time.NewTicker(app.UpdateInterval)
	defer ticker.Stop()

	log.Info().
		Str("interval", app.UpdateInterval.String()).
		Int("retries", app.MaxRetries).
		Msg("plex refresher started")

	for {
		log.Info().Msg("starting library refresh")
		err := plex.RefreshAllLibraries(ctx, app)
		if err != nil {
			log.Error().Err(err).Msg("refresh cycle failed")
		}

		select {
		case <-ctx.Done():
			log.Info().Msg("shutdown requested")
			return

		case <-ticker.C:
		}
	}
}

func main() {
	app := config.CreateConfig()
	log.Logger = config.CreateLogger(app.LogFormat)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-sig
		log.Info().
			Str("signal", s.String()).
			Msg("received shutdown signal")
		cancel()
	}()

	run(ctx, app)
	log.Info().Msg("shutdown complete")
}
