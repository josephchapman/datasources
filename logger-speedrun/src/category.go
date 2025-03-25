package main

import (
	"datasources/cmn"
	"encoding/json"
	"fmt"
)

type category struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (c category) endpoint() (url string, err error) {
	url = fmt.Sprintf("https://www.speedrun.com/api/v1/categories/%s", c.Id)

	return url, nil
}

// update API Data
func (c *category) updateAPI() (err error) {
	// Get the endpoint from the location
	url, err := c.endpoint()
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

	// Unmarshal the JSON data into the category struct
	err = json.Unmarshal(jsonData, &c)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal: %w", err)
		return cmn.LoggedError(err)
	}

	return nil
}
