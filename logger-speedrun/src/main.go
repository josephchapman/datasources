package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	// Read the environment variable
	envVar := os.Getenv("SPEEDRUN_LEADERBOARDS")
	if envVar == "" {
		fmt.Println("SPEEDRUN_LEADERBOARDS environment variable is not set")
		return
	}

	// Unmarshal the JSON array from the environment variable
	var leaderboards []leaderboard
	err := json.Unmarshal([]byte(envVar), &leaderboards)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}

	// Iterate through the slice and process each leaderboard
	for _, l := range leaderboards {
		l.updateAPI()
		entry, _ := l.NewLogEntry()
		entry.log()
	}
}
