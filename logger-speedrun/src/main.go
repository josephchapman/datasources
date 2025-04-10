package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/josephchapman/datasources/cmn"
)

const applicationName = "logger-speedrun"

func main() {
	cmn.SetApplicationName(applicationName)
	cmn.NoCron(runTask)
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
