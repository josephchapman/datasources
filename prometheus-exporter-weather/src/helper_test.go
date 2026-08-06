package main

import (
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestOsEnvVarToLocations(t *testing.T) {
	// Set up the environment variable
	os.Setenv("WEATHER_LOCATIONS", `[{"Name":"New York","Latitude":40.7128,"Longitude":-74.0060},{"Name":"Los Angeles","Latitude":34.0522,"Longitude":-118.2437}]`)
	defer os.Unsetenv("WEATHER_LOCATIONS")

	// Call the function
	locations := osEnvVarToLocations()

	// Check the results
	if len(locations) != 2 {
		t.Fatalf("expected 2 locations, got %d", len(locations))
	}

	if locations[0].Name != "New York" {
		t.Errorf("expected first location to be New York, got %s", locations[0].Name)
	}

	if locations[1].Name != "Los Angeles" {
		t.Errorf("expected second location to be Los Angeles, got %s", locations[1].Name)
	}
}

func TestUpdateMetricsIncludesAdditionalWeatherFields(t *testing.T) {
	_, m := definePrometheusRegistry()
	w := weather{
		Location: location{Name: "Test City"},
		ApiData: apiData{
			Current: apiDataCurrent{
				Is_day:           1,
				Pressure_msl:     1013.25,
				Surface_pressure: 1008.5,
				Weather_code:     3,
				Dew_point_2m:     10.5,
				Visibility:       10000,
			},
		},
		Metrics: *m,
	}

	if err := w.updateMetrics(); err != nil {
		t.Fatalf("updateMetrics returned error: %v", err)
	}

	checks := []struct {
		name     string
		metric   func() float64
		expected float64
	}{
		{name: "is_day", metric: func() float64 { return testutil.ToFloat64(m.IsDay.WithLabelValues("Test City")) }, expected: 1},
		{name: "pressure_msl", metric: func() float64 { return testutil.ToFloat64(m.PressureMSL.WithLabelValues("Test City")) }, expected: 1013.25},
		{name: "surface_pressure", metric: func() float64 { return testutil.ToFloat64(m.SurfacePressure.WithLabelValues("Test City")) }, expected: 1008.5},
		{name: "weather_code", metric: func() float64 { return testutil.ToFloat64(m.WeatherCode.WithLabelValues("Test City")) }, expected: 3},
		{name: "dew_point_2m", metric: func() float64 { return testutil.ToFloat64(m.DewPoint2m.WithLabelValues("Test City")) }, expected: 10.5},
		{name: "visibility", metric: func() float64 { return testutil.ToFloat64(m.Visibility.WithLabelValues("Test City")) }, expected: 10000},
	}

	for _, check := range checks {
		if got := check.metric(); got != check.expected {
			t.Errorf("%s: expected %v, got %v", check.name, check.expected, got)
		}
	}
}
