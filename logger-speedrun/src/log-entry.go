package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type logEntry struct {
	Game      string    `json:"game"`
	Category  string    `json:"category"`
	Player    string    `json:"player"`
	Time      string    `json:"time"`
	Performed performed `json:"submitted"`
}

type performed struct {
	Datetime  string `json:"date_time"`
	Timesince string `json:"time_since"`
}

func (l logEntry) log() (err error) {
	log.SetOutput(os.Stdout)
	log.SetFlags(0) // Disable date and timestamps

	// lData, err := json.MarshalIndent(l, "", "  ")
	lData, err := json.Marshal(l)
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}
	log.Println(string(lData))

	return nil
}
