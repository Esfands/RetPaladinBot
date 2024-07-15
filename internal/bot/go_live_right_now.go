package bot

import (
	"fmt"
	"time"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/go-co-op/gocron"
)

func GoLiveRightNow(gctx global.Context, bot *Connection) error {
	// Create a new scheduler
	s := gocron.NewScheduler(time.UTC)

	// Set the timezone to CST
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return nil
	}
	s.ChangeLocation(loc)

	// Define the job
	job := func() {
		fmt.Println("Running job at 12:00 PM CST")
	}

	// Schedule the job to run every day at 12:00 PM CST
	_, err = s.Every(1).Day().At("11:38").Do(job)
	if err != nil {
		fmt.Println("Error scheduling job:", err)
		return nil
	}

	// Start the scheduler
	s.StartBlocking()

	return nil
}
