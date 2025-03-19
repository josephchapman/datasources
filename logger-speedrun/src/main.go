package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	// Define the -nocron flag to accept an integer value for minutes
	sleepMinutes := flag.Int("nocron", 0, "Run in a loop with a sleep interval in minutes")
	flag.Parse()

	for {
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

		// If -nocron flag is not set, exit the loop
		if *sleepMinutes <= 0 {
			break
		}

		// Sleep for the specified number of minutes before the next iteration
		time.Sleep(time.Duration(*sleepMinutes) * time.Minute)
	}
}
