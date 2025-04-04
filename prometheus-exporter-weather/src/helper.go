package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/josephchapman/datasources/cmn"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func periodicCheck(locations []location, m *metrics) (err error) {
	weathers := make([]weather, len(locations))

	// Initialize weather structs for each location
	for i, loc := range locations {
		w, err := NewWeather(loc, m)
		if err != nil {
			err = fmt.Errorf("NewWeather: %w", err)
			return cmn.LoggedError(err)
		}
		weathers[i] = w
	}

	for {
		now := time.Now()
		minute := now.Minute()
		second := now.Second()

		// API updates every 00, 15, 30, 45 mins
		// Calculate the next target time (2.5, 17.5, 32.5, 47.5 minutes)
		var nextTarget time.Time
		switch {
		case minute < 2 || (minute == 2 && second < 30):
			nextTarget = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 2, 30, 0, now.Location())
		case minute < 17 || (minute == 17 && second < 30):
			nextTarget = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 17, 30, 0, now.Location())
		case minute < 32 || (minute == 32 && second < 30):
			nextTarget = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 32, 30, 0, now.Location())
		case minute < 47 || (minute == 47 && second < 30):
			nextTarget = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 47, 30, 0, now.Location())
		default:
			nextTarget = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 2, 30, 0, now.Location())
		}

		// Calculate the sleep duration until the next target time
		sleepDuration := time.Until(nextTarget)
		str := fmt.Sprintf("Sleeping until next API update at %v (for %v).", nextTarget, sleepDuration)
		cmn.LogOut.Info(str)
		time.Sleep(sleepDuration)

		// Update API and metrics
		for i, w := range weathers {
			w.updateAPI()
			w.updateMetrics()
			weathers[i] = w
		}
	}
}

func definePrometheusRegistry() (*prometheus.Registry, *metrics) {
	r := prometheus.NewRegistry()

	metrics := &metrics{
		TemperatureActual: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "weather_temperature_actual_celcius",
			Help: "The actual temperature.",
		}, []string{"location"}),

		TemperatureApparent: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "weather_temperature_apparent_celcius",
			Help: "The apparent temperature.",
		}, []string{"location"}),

		HumidityRelative: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "weather_humidity_relative_percent",
			Help: "The relative humidity.",
		}, []string{"location"}),

		Precipitation: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "weather_precipitation_milimeters",
			Help: "The precipitation.",
		}, []string{"location"}),

		Rain: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "weather_rain_milimeters",
			Help: "The rain.",
		}, []string{"location"}),

		Showers: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "weather_showers_milimeters",
			Help: "The showers.",
		}, []string{"location"}),

		CloudCover: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "weather_cloud_cover_percent",
			Help: "The cloud cover.",
		}, []string{"location"}),

		WindSpeed: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "weather_wind_speed_knots",
			Help: "The wind speed.",
		}, []string{"location"}),

		WindDirection: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "weather_wind_direction_degrees",
			Help: "The wind direction.",
		}, []string{"location"}),

		WindGusts: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "weather_wind_gusts_knots",
			Help: "The wind gusts.",
		}, []string{"location"}),
	}

	r.MustRegister(metrics.TemperatureActual)
	r.MustRegister(metrics.TemperatureApparent)
	r.MustRegister(metrics.HumidityRelative)
	r.MustRegister(metrics.Precipitation)
	r.MustRegister(metrics.Rain)
	r.MustRegister(metrics.Showers)
	r.MustRegister(metrics.CloudCover)
	r.MustRegister(metrics.WindSpeed)
	r.MustRegister(metrics.WindDirection)
	r.MustRegister(metrics.WindGusts)

	return r, metrics
}

func osEnvVarToLocations() (locations []location) {
	weatherLocations := os.Getenv("WEATHER_LOCATIONS")
	if weatherLocations == "" {
		err := fmt.Errorf("WEATHER_LOCATIONS environment variable is not set")
		cmn.LoggedError(err)
		panic(err)
	}

	// Unmarshal the JSON data into a slice of location structs
	err := json.Unmarshal([]byte(weatherLocations), &locations)
	if err != nil {
		err := fmt.Errorf("error unmarshalling WEATHER_LOCATIONS: %v", err)
		cmn.LoggedError(err)
		panic(err)
	}
	return locations
}
