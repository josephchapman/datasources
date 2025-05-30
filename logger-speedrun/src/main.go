package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/josephchapman/datasources/cmn"
)

const applicationName = "logger-speedrun"

func main() {
	cmn.SetApplicationName(applicationName)

	sleepMinutes := flag.Int("nocron", 0, "Run in a loop with a sleep interval in minutes")
	flag.Parse()

	cmn.NoCron(runTask, *sleepMinutes)
}

func runTask() {
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
}
