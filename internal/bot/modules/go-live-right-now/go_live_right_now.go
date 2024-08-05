package goliverightnow

import (
	"fmt"
	"time"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/go-co-op/gocron"
	"golang.org/x/exp/slog"
)

type GoLiveRightNowModule struct {
	scheduler *gocron.Scheduler
}

func NewGoLiveRightNowModule(gctx global.Context, client *twitch.Client) *GoLiveRightNowModule {
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
		streamStatus, err := gctx.Crate().Turso.Queries().GetMostRecentStreamStatus(gctx)
		if err != nil {
			slog.Error("[go-live-right-now] Error getting most recent stream status", "error", err)
			return
		}

		// Say GOLIVERIGHTNOWMADGE if the stream isn't live
		if !streamStatus.Live {
			client.Say(gctx.Config().Twitch.Bot.Channel, "GOLIVERIGHTNOWMADGE")
		}
	}

	// Schedule the job to run every day at 12:00 PM CST
	_, err = s.Every(1).Day().At("12:00").Do(job)
	if err != nil {
		fmt.Println("Error scheduling job:", err)
		return nil
	}

	// Start the scheduler in async mode
	s.StartAsync()

	// Listen for context done signal to stop the scheduler
	go func() {
		<-gctx.Done()
		s.Stop()
		fmt.Println("Scheduler stopped")
	}()

	return &GoLiveRightNowModule{
		scheduler: s,
	}
}
