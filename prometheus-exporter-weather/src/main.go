package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Read the WEATHER_LOCATIONS environment variable
	locations := osEnvVarToLocations()

	r, m := definePrometheusRegistry()

	go func() {
		periodicCheck(locations, m)
	}()

	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	http.Handle("/metrics", handler)

	fmt.Fprintln(os.Stderr, "Listening :2112/tcp")
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
