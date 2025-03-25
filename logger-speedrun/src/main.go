package main

import (
	"datasources/cmn"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

const applicationName = "logger-speedrun"

func main() {
	cmn.SetApplicationName(applicationName)

	// Define the -nocron flag to accept an integer value for minutes
	sleepMinutes := flag.Int("nocron", 0, "Run in a loop with a sleep interval in minutes")
	flag.Parse()

	for {
		// Read the environment variable
		envVar := os.Getenv("SPEEDRUN_LEADERBOARDS")
		if envVar == "" {
			err := fmt.Errorf("SPEEDRUN_LEADERBOARDS environment variable is not set")
			cmn.LoggedError(err)
			panic(err)
		}

		// Unmarshal the JSON array from the environment variable
		var leaderboards []leaderboard
		err := json.Unmarshal([]byte(envVar), &leaderboards)
		if err != nil {
			err = fmt.Errorf("json.Unmarshal: %w", err)
			cmn.LoggedError(err)
			panic(err)
		}

		// Iterate through the slice and process each leaderboard
		for _, l := range leaderboards {
			l.updateAPI()
			cr, _ := l.NewCurrentRecord()
			err := cr.log()
			if err != nil {
				err = fmt.Errorf("cr.log(): %w", err)
				cmn.LoggedError(err)
			}
		}

		// If -nocron flag is not set, exit the loop
		if *sleepMinutes <= 0 {
			break
		}

		// Sleep for the specified number of minutes before the next iteration
		time.Sleep(time.Duration(*sleepMinutes) * time.Minute)
	}
}
