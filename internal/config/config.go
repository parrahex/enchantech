package config

import (
	"fmt"
	"os"
	"strconv"
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

	// City is the location whose weather decides the sky. It is unset by
	// default; without it the forecast is skipped and a calm default sky is
	// used, so the application needs no configuration to run.
	City City
}

// City locates the weather used to pick the sky scene.
type City struct {
	Latitude  float64
	Longitude float64

	// TimeZone is an IANA name such as Europe/Paris. An empty value disables
	// the forecast.
	TimeZone string
}

// Configured reports whether a location was supplied.
func (city City) Configured() bool {
	return city.TimeZone != ""
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

	city, err := cityValue()
	if err != nil {
		return Config{}, err
	}
	settings.City = city

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

// cityValue reads the location used for the forecast. All three variables must
// be supplied together, or none of them.
func cityValue() (City, error) {
	timeZone := strings.TrimSpace(os.Getenv("APP_CITY_TIMEZONE"))
	rawLatitude := strings.TrimSpace(os.Getenv("APP_CITY_LATITUDE"))
	rawLongitude := strings.TrimSpace(os.Getenv("APP_CITY_LONGITUDE"))

	if timeZone == "" && rawLatitude == "" && rawLongitude == "" {
		return City{}, nil
	}

	if timeZone == "" || rawLatitude == "" || rawLongitude == "" {
		return City{}, fmt.Errorf(
			"APP_CITY_TIMEZONE, APP_CITY_LATITUDE, and APP_CITY_LONGITUDE must be set together",
		)
	}

	latitude, err := coordinateValue("APP_CITY_LATITUDE", rawLatitude, 90)
	if err != nil {
		return City{}, err
	}

	longitude, err := coordinateValue("APP_CITY_LONGITUDE", rawLongitude, 180)
	if err != nil {
		return City{}, err
	}

	return City{Latitude: latitude, Longitude: longitude, TimeZone: timeZone}, nil
}

// coordinateValue parses one coordinate and checks it against its range.
func coordinateValue(name string, rawValue string, limit float64) (float64, error) {
	value, err := strconv.ParseFloat(rawValue, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number such as 48.85: %w", name, err)
	}

	if value < -limit || value > limit {
		return 0, fmt.Errorf("%s must be between -%g and %g", name, limit, limit)
	}

	return value, nil
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
