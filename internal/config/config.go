package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	// Address is where the application listens.
	Address string

	// ShutdownTimeout limits graceful shutdown.
	ShutdownTimeout time.Duration

	// ContentDir holds the site's templates and assets. It is kept outside the
	// repository, so the application starts and serves a placeholder when the
	// directory is absent.
	ContentDir string
}

// Load reads environment variables, applies defaults, and validates required
// values. It intentionally does not load .env files itself.
func Load() (Config, error) {
	settings := Config{
		Address:         ":8080",
		ShutdownTimeout: 5 * time.Second,
		ContentDir:      "content",
	}

	if directory := strings.TrimSpace(os.Getenv("APP_CONTENT_DIR")); directory != "" {
		settings.ContentDir = directory
	}

	if address := strings.TrimSpace(os.Getenv("APP_ADDRESS")); address != "" {
		settings.Address = address
	}

	shutdownTimeout, err := durationValue(
		"APP_SHUTDOWN_TIMEOUT",
		settings.ShutdownTimeout,
	)
	if err != nil {
		return Config{}, err
	}
	settings.ShutdownTimeout = shutdownTimeout

	return settings, nil
}

// durationValue reads a positive duration or returns the supplied default.
func durationValue(name string, defaultValue time.Duration) (time.Duration, error) {
	rawValue := strings.TrimSpace(os.Getenv(name))
	if rawValue == "" {
		return defaultValue, nil
	}

	duration, err := time.ParseDuration(rawValue)
	if err != nil {
		return 0, fmt.Errorf(
			"%s must be a valid duration such as 5s or 10s: %w",
			name,
			err,
		)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", name)
	}

	return duration, nil
}
