package cmn

import (
	"flag"
	"time"
)

func NoCron(task func()) {
	// Define the -nocron flag to accept an integer value for minutes
	sleepMinutes := flag.Int("nocron", 0, "Run in a loop with a sleep interval in minutes")
	flag.Parse()

	if *sleepMinutes <= 0 {
		// Run the task once if -nocron is not set or set to 0
		task()
		return
	}

	// Run the task in a loop every x minutes
	ticker := time.NewTicker(time.Duration(*sleepMinutes) * time.Minute)
	defer ticker.Stop()

	for {
		task()
		<-ticker.C // Wait for the next tick
	}
}
