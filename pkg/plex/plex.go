package plex

import (
	"context"
	"encoding/xml"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/verdel/plex-library-updater/pkg/config"
)

// MediaContainer represents the XML response from Plex that contains library
// directories (sections).
type MediaContainer struct {
	Directories []Directory `xml:"Directory"`
}

// Directory represents a single Plex library section with its key and title.
type Directory struct {
	Key   string `xml:"key,attr"`
	Title string `xml:"title,attr"`
}

func retry(ctx context.Context, maxRetries int, fn func() error) error {
	var err error
	backoff := time.Second
	for attempt := 0; attempt <= maxRetries; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		if attempt == maxRetries {
			break
		}

		log.Warn().
			Err(err).
			Int("attempt", attempt+1).
			Msg("retrying request")

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
		backoff *= 2
	}
	return err
}

func getLibraries(ctx context.Context, app config.Config) ([]Directory, error) {
	var result []Directory

	err := retry(ctx, app.MaxRetries, func() error {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			app.PlexUrl+"/library/sections?X-Plex-Token="+app.PlexToken,
			nil,
		)

		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode >= 500 {
			return errors.New("plex server error")
		}

		var container MediaContainer

		if err := xml.NewDecoder(resp.Body).Decode(&container); err != nil {
			return err
		}

		result = container.Directories
		return nil
	})

	return result, err
}

func refreshLibrary(ctx context.Context, app config.Config, lib Directory) error {
	return retry(ctx, app.MaxRetries, func() error {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			app.PlexUrl+"/library/sections/"+lib.Key+"/refresh?X-Plex-Token="+app.PlexToken,
			nil,
		)

		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode >= 500 {
			return errors.New("plex server error")
		}

		log.Info().
			Str("library", lib.Title).
			Int("status", resp.StatusCode).
			Msg("refresh triggered")

		return nil
	})
}

// RefreshAllLibraries fetches all Plex library sections and triggers a refresh
// request for each one concurrently. Errors from individual refreshes are
// logged; this function returns an error only if retrieving the list of
// libraries fails.
func RefreshAllLibraries(ctx context.Context, app config.Config) error {
	libs, err := getLibraries(ctx, app)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, lib := range libs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := refreshLibrary(ctx, app, lib)
			if err != nil {
				log.Error().
					Err(err).
					Str("library", lib.Key).
					Msg("refresh failed after retries")
			}
		}()
	}
	wg.Wait()
	return nil
}
