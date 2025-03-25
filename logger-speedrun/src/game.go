package main

import (
	"encoding/json"
	"fmt"

	"github.com/josephchapman/datasources/cmn"
)

type game struct {
	Id    string            `json:"id"`
	Names map[string]string `json:"names"`
}

func (g game) endpoint() (url string, err error) {
	url = fmt.Sprintf("https://www.speedrun.com/api/v1/games/%s", g.Id)

	return url, nil
}

// update API Data
func (g *game) updateAPI() (err error) {
	// Get the endpoint from the location
	url, err := g.endpoint()
	if err != nil {
		err = fmt.Errorf("w.Location.endpoint(): %w", err)
		return cmn.LoggedError(err)
	}

	// Query the endpoint to receive updated data
	data, err := queryAPI(url)
	if err != nil {
		err = fmt.Errorf("queryAPI(): %w", err)
		return cmn.LoggedError(err)
	}

	// Convert the map to JSON
	jsonData, err := json.Marshal(data["data"])
	if err != nil {
		err = fmt.Errorf("json.Marshal: %w", err)
		return cmn.LoggedError(err)
	}

	// Unmarshal the JSON data into the game struct
	err = json.Unmarshal(jsonData, &g)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal: %w", err)
		return cmn.LoggedError(err)
	}

	return nil
}
