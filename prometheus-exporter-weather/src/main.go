package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const applicationName = "prometheus-exporter-weather"

func main() {
	// Read the WEATHER_LOCATIONS environment variable
	locations := osEnvVarToLocations()

	r, m := definePrometheusRegistry()

	go func() {
		periodicCheck(locations, m)
	}()

	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	http.Handle("/metrics", handler)

	logOut.Info("Listening :2112/tcp")
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		LoggedError(err)
	}
}
