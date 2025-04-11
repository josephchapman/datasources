package cmn

import (
	"time"
)

func NoCron(task func(), sleepMinutes int) {
	if sleepMinutes <= 0 {
		// Run the task once if sleepMinutes is 0 or less
		task()
		return
	}

	// Run the task in a loop every x minutes
	ticker := time.NewTicker(time.Duration(sleepMinutes) * time.Minute)
	defer ticker.Stop()

	for {
		task()
		<-ticker.C // Wait for the next tick
	}
}
