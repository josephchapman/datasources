package main

import (
	"encoding/json"
	"fmt"

	"github.com/josephchapman/datasources/cmn"
)

type user struct {
	Id    string            `json:"id"`
	Names map[string]string `json:"names"`
}

func (u user) endpoint() (url string, err error) {
	url = fmt.Sprintf("https://www.speedrun.com/api/v1/users/%s", u.Id)

	return url, nil
}

// update API Data
func (u *user) updateAPI() (err error) {
	// Get the endpoint from the location
	url, err := u.endpoint()
	if err != nil {
		err = fmt.Errorf("w.Location.endpoint(): %w", err)
		return cmn.LoggedError(err)
	}

	// Query the endpoint to receive updated data
	data, err := cmn.QueryAPI(url)
	if err != nil {
		err = fmt.Errorf("cmn.QueryAPI(): %w", err)
		return cmn.LoggedError(err)
	}

	// Convert the map to JSON
	jsonData, err := json.Marshal(data["data"])
	if err != nil {
		err = fmt.Errorf("json.Marshal: %w", err)
		return cmn.LoggedError(err)
	}

	// Unmarshal the JSON data into the user struct
	err = json.Unmarshal(jsonData, &u)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal: %w", err)
		return cmn.LoggedError(err)
	}

	return nil
}
