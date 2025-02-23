package main

import (
	"encoding/json"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	TZData    string  `json:"tzdata"`
}

func (l location) endpoint() (url string, err error) {

	url = fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%2f&longitude=%2f&current=temperature_2m,relative_humidity_2m,apparent_temperature,precipitation,rain,showers,cloud_cover,wind_speed_10m,wind_direction_10m,wind_gusts_10m&wind_speed_unit=kn&timezone=%s", l.Latitude, l.Longitude, l.TZData)

	return url, nil
}

type apiData struct {
	Latitude              float64             `json:"latitude"`
	Longitude             float64             `json:"longitude"`
	Utc_offset_seconds    float64             `json:"utc_offset_seconds"`
	Timezone              string              `json:"timezone"`
	Timezone_abbreviation string              `json:"timezone_abbreviation"`
	Elevation             float64             `json:"elevation"`
	Current_units         apiDataCurrentUnits `json:"current_units"`
	Current               apiDataCurrent      `json:"current"`
}

type metrics struct {
	TemperatureActual   *prometheus.GaugeVec `json:"temperature_actual"`
	TemperatureApparent *prometheus.GaugeVec `json:"temperature_apparent"`
	HumidityRelative    *prometheus.GaugeVec `json:"humidity_relative"`
	Precipitation       *prometheus.GaugeVec `json:"precipitation"`
	Rain                *prometheus.GaugeVec `json:"rain"`
	Showers             *prometheus.GaugeVec `json:"showers"`
	CloudCover          *prometheus.GaugeVec `json:"cloud_cover"`
	WindSpeed           *prometheus.GaugeVec `json:"wind_speed"`
	WindDirection       *prometheus.GaugeVec `json:"wind_direction"`
	WindGusts           *prometheus.GaugeVec `json:"wind_gusts"`
}

type apiDataCurrentUnits struct {
	Time                 string `json:"time"`
	Interval             string `json:"interval"`
	Temperature_2m       string `json:"temperature_2m"`
	Relative_humidity_2m string `json:"relative_humidity_2m"`
	Apparent_temperature string `json:"apparent_temperature"`
	Precipitation        string `json:"precipitation"`
	Rain                 string `json:"rain"`
	Showers              string `json:"showers"`
	Cloud_cover          string `json:"cloud_cover"`
	Wind_speed_10m       string `json:"wind_speed_10m"`
	Wind_direction_10m   string `json:"wind_direction_10m"`
	Wind_gusts_10m       string `json:"wind_gusts_10m"`
}

type apiDataCurrent struct {
	Time                 string  `json:"time"`
	Interval             float64 `json:"interval"`
	Temperature_2m       float64 `json:"temperature_2m"`
	Relative_humidity_2m float64 `json:"relative_humidity_2m"`
	Apparent_temperature float64 `json:"apparent_temperature"`
	Precipitation        float64 `json:"precipitation"`
	Rain                 float64 `json:"rain"`
	Showers              float64 `json:"showers"`
	Cloud_cover          float64 `json:"cloud_cover"`
	Wind_speed_10m       float64 `json:"wind_speed_10m"`
	Wind_direction_10m   float64 `json:"wind_direction_10m"`
	Wind_gusts_10m       float64 `json:"wind_gusts_10m"`
}

type weather struct {
	Location location `json:"location"`
	ApiData  apiData  `json:"api_data"`
	Metrics  metrics  `json:"metrics"`
}

func NewWeather(l location, m *metrics) (w weather, err error) {
	w = weather{
		Location: l,
		Metrics:  *m,
	}
	err = w.updateAPI()
	if err != nil {
		err = fmt.Errorf("w.updateAPI(): %w", err)
		return weather{}, err
	}

	err = w.updateMetrics()
	if err != nil {
		err = fmt.Errorf("w.updateMetrics(): %w", err)
		return weather{}, err
	}

	return w, nil
}

func (w weather) printToConsole() (err error) {

	fmt.Printf("Location: %v\n", w.Location)
	fmt.Printf("Latitude: %.2f\n", w.ApiData.Latitude)
	fmt.Printf("Longitude: %.2f\n", w.ApiData.Longitude)
	fmt.Printf("Utc_offset_seconds: %.0f\n", w.ApiData.Utc_offset_seconds)
	fmt.Printf("Timezone: %s\n", w.ApiData.Timezone)
	fmt.Printf("Timezone_abbreviation: %s\n", w.ApiData.Timezone_abbreviation)
	fmt.Printf("Elevation: %.0f\n", w.ApiData.Elevation)
	fmt.Printf("Time: %s %s\n", w.ApiData.Current.Time, w.ApiData.Current_units.Time)
	fmt.Printf("Interval: %.0f %s\n", w.ApiData.Current.Interval, w.ApiData.Current_units.Interval)
	fmt.Printf("Temperature_2m: %.1f %s\n", w.ApiData.Current.Temperature_2m, w.ApiData.Current_units.Temperature_2m)
	fmt.Printf("Relative_humidity_2m: %.0f %s\n", w.ApiData.Current.Relative_humidity_2m, w.ApiData.Current_units.Relative_humidity_2m)
	fmt.Printf("Apparent_temperature: %.1f %s\n", w.ApiData.Current.Apparent_temperature, w.ApiData.Current_units.Apparent_temperature)
	fmt.Printf("Precipitation: %.0f %s\n", w.ApiData.Current.Precipitation, w.ApiData.Current_units.Precipitation)
	fmt.Printf("Rain: %.0f %s\n", w.ApiData.Current.Rain, w.ApiData.Current_units.Rain)
	fmt.Printf("Showers: %.0f %s\n", w.ApiData.Current.Showers, w.ApiData.Current_units.Showers)
	fmt.Printf("Cloud_cover: %.0f %s\n", w.ApiData.Current.Cloud_cover, w.ApiData.Current_units.Cloud_cover)
	fmt.Printf("Wind_speed_10m: %.1f %s\n", w.ApiData.Current.Wind_speed_10m, w.ApiData.Current_units.Wind_speed_10m)
	fmt.Printf("Wind_direction_10m: %.0f %s\n", w.ApiData.Current.Wind_direction_10m, w.ApiData.Current_units.Wind_direction_10m)
	fmt.Printf("Wind_gusts_10m: %.1f %s\n", w.ApiData.Current.Wind_gusts_10m, w.ApiData.Current_units.Wind_gusts_10m)
	return nil
}

// update API Data
func (w *weather) updateAPI() (err error) {
	// Get the endpoint from the location
	url, err := w.Location.endpoint()
	if err != nil {
		err = fmt.Errorf("w.Location.endpoint(): %w", err)
		return err
	}

	// Query the endpoint to receive updated data
	data, err := queryAPI(url)
	if err != nil {
		err = fmt.Errorf("queryAPI(): %w", err)
		return err
	}

	// Convert the map to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("json.Marshal: %w", err)
		return err
	}

	// Unmarshal the JSON data into the weather struct
	err = json.Unmarshal(jsonData, &w.ApiData)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal: %w", err)
		return err
	}

	return nil
}

// update metrics with API data
func (w *weather) updateMetrics() (err error) {

	w.Metrics.TemperatureActual.With(
		prometheus.Labels{"location": w.Location.Name},
	).Set(w.ApiData.Current.Temperature_2m)

	w.Metrics.TemperatureApparent.With(
		prometheus.Labels{"location": w.Location.Name},
	).Set(w.ApiData.Current.Apparent_temperature)

	w.Metrics.HumidityRelative.With(
		prometheus.Labels{"location": w.Location.Name},
	).Set(w.ApiData.Current.Relative_humidity_2m)

	w.Metrics.Precipitation.With(
		prometheus.Labels{"location": w.Location.Name},
	).Set(w.ApiData.Current.Precipitation)

	w.Metrics.Rain.With(
		prometheus.Labels{"location": w.Location.Name},
	).Set(w.ApiData.Current.Rain)

	w.Metrics.Showers.With(
		prometheus.Labels{"location": w.Location.Name},
	).Set(w.ApiData.Current.Showers)

	w.Metrics.CloudCover.With(
		prometheus.Labels{"location": w.Location.Name},
	).Set(w.ApiData.Current.Cloud_cover)

	w.Metrics.WindSpeed.With(
		prometheus.Labels{"location": w.Location.Name},
	).Set(w.ApiData.Current.Wind_speed_10m)

	w.Metrics.WindDirection.With(
		prometheus.Labels{"location": w.Location.Name},
	).Set(w.ApiData.Current.Wind_direction_10m)

	w.Metrics.WindGusts.With(
		prometheus.Labels{"location": w.Location.Name},
	).Set(w.ApiData.Current.Wind_gusts_10m)

	return nil
}
