package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	_ "time/tzdata"

	"github.com/parrahex/enchantech/internal/config"
)

const (
	weatherEndpoint = "https://api.open-meteo.com/v1/forecast"
	weatherRefresh  = 15 * time.Minute
	weatherTimeout  = 10 * time.Second
)

type conditions struct {
	Known       bool
	Night       bool
	Code        int
	Cover       int
	Wind        float64
	Temperature float64
}

var unknownConditions = conditions{Cover: 50, Wind: 3}

type weather struct {
	place    config.City
	location *time.Location
	client   *http.Client

	mutex   sync.RWMutex
	current conditions

	cancel context.CancelFunc
	done   chan struct{}
}

func newWeather(place config.City) *weather {
	location := time.UTC

	if place.Configured() {
		loaded, err := time.LoadLocation(place.TimeZone)
		if err != nil {
			slog.Warn("unknown time zone, using UTC", "timezone", place.TimeZone, "error", err)
		} else {
			location = loaded
		}
	}

	return &weather{
		place:    place,
		location: location,
		client:   &http.Client{Timeout: weatherTimeout},
		current:  unknownConditions,
		done:     make(chan struct{}),
	}
}

func (service *weather) start() {
	runContext, cancel := context.WithCancel(context.Background())
	service.cancel = cancel

	if !service.place.Configured() {
		slog.Info("no city configured, using the default sky")

		close(service.done)

		return
	}

	go func() {
		defer close(service.done)

		service.refresh(runContext)

		ticker := time.NewTicker(weatherRefresh)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				service.refresh(runContext)
			case <-runContext.Done():
				return
			}
		}
	}()
}

func (service *weather) stop() {
	service.cancel()
	<-service.done
}

func (service *weather) reading() conditions {
	service.mutex.RLock()
	defer service.mutex.RUnlock()

	return service.current
}

type moment struct {
	hour   int
	moon   float64
	season string
}

func (service *weather) moment() moment {
	now := time.Now().In(service.location)

	return moment{
		hour:   now.Hour(),
		moon:   moonPhase(now, service.place.Latitude),
		season: seasonFor(now, service.place.Latitude),
	}
}

func (service *weather) refresh(ctx context.Context) {
	current, err := service.read(ctx)
	if err != nil {
		slog.Warn("weather unavailable, using default sky", "error", err)

		service.store(unknownConditions)

		return
	}

	service.store(current)
}

func (service *weather) store(current conditions) {
	service.mutex.Lock()
	defer service.mutex.Unlock()

	service.current = current
}

func (service *weather) read(ctx context.Context) (conditions, error) {
	requestContext, cancel := context.WithTimeout(ctx, weatherTimeout)
	defer cancel()

	query := url.Values{
		"latitude":      {strconv.FormatFloat(service.place.Latitude, 'f', -1, 64)},
		"longitude":     {strconv.FormatFloat(service.place.Longitude, 'f', -1, 64)},
		"timezone":      {service.place.TimeZone},
		"current":       {"is_day,weather_code,cloud_cover,wind_speed_10m,temperature_2m"},
		"forecast_days": {"1"},
	}

	request, err := http.NewRequestWithContext(
		requestContext,
		http.MethodGet,
		weatherEndpoint+"?"+query.Encode(),
		nil,
	)
	if err != nil {
		return conditions{}, fmt.Errorf("build weather request: %w", err)
	}

	response, err := service.client.Do(request)
	if err != nil {
		return conditions{}, fmt.Errorf("call weather service: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return conditions{}, fmt.Errorf("weather service returned %s", response.Status)
	}

	var payload struct {
		Current struct {
			Time        string  `json:"time"`
			IsDay       int     `json:"is_day"`
			Code        int     `json:"weather_code"`
			Cover       int     `json:"cloud_cover"`
			Wind        float64 `json:"wind_speed_10m"`
			Temperature float64 `json:"temperature_2m"`
		} `json:"current"`
	}

	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return conditions{}, fmt.Errorf("decode weather response: %w", err)
	}

	if payload.Current.Time == "" {
		return conditions{}, fmt.Errorf("weather response has no current reading")
	}

	return conditions{
		Known:       true,
		Night:       payload.Current.IsDay == 0,
		Code:        payload.Current.Code,
		Cover:       payload.Current.Cover,
		Wind:        payload.Current.Wind,
		Temperature: payload.Current.Temperature,
	}, nil
}
