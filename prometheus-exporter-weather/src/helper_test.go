package main

import (
	"os"
	"testing"
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

func TestOsEnvVarToLocations_InvalidJSON(t *testing.T) {
	// Set up the environment variable with invalid JSON
	os.Setenv("WEATHER_LOCATIONS", `invalid json`)
	defer os.Unsetenv("WEATHER_LOCATIONS")

	// Call the function
	locations := osEnvVarToLocations()

	// Check the results
	if len(locations) != 0 {
		t.Fatalf("expected 0 locations, got %d", len(locations))
	}
}
