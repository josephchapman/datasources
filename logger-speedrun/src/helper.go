package main

import (
	"datasources/cmn"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// queries a URL and returns the data as a Golang map
func queryAPI(url string) (data map[string]interface{}, err error) {

	// Get all data from URL
	resp, err := http.Get(url)
	if err != nil {
		err = fmt.Errorf("http.Get: %w", err)
		return nil, cmn.LoggedError(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("resp.StatusCode: %d", resp.StatusCode)
		return nil, cmn.LoggedError(err)
	}

	// Get the body of the response from the ReaderCloser interface into a Go variable 'body'
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("io.ReadAll: %w", err)
		return nil, cmn.LoggedError(err)
	}

	// Convert the JSON data within 'body' to a Golang map in the 'data' var
	err = json.Unmarshal(body, &data)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal: %w", err)
		return nil, cmn.LoggedError(err)
	}

	return data, cmn.LoggedError(err)
}
