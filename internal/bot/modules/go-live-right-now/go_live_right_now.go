package goliverightnow

import (
	"fmt"
	"time"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/go-co-op/gocron"
)

type GoLiveRightNowModule struct{}

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
		client.Say(gctx.Config().Twitch.Bot.Channel, "GOLIVERIGHTNOWMADGE")
	}

	// Schedule the job to run every day at 12:00 PM CST
	_, err = s.Every(1).Day().At("12:00").Do(job)
	if err != nil {
		fmt.Println("Error scheduling job:", err)
		return nil
	}

	// Start the scheduler
	s.StartBlocking()

	go func() {
		<-gctx.Done()
		s.Stop()
	}()

	return &GoLiveRightNowModule{}
}
