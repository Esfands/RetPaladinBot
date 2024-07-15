package goliverightnow

import (
	"fmt"
	"time"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/go-co-op/gocron"
	"github.com/nicklaw5/helix/v2"
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
		channelInfo, err := gctx.Crate().Helix.Client().GetStreams(&helix.StreamsParams{
			UserIDs: []string{gctx.Config().Twitch.Bot.ChannelID},
		})
		if err != nil {
			fmt.Println("Error getting channel information in go live right now module:", err)
			return
		}

		// Say GOLIVERIGHTNOWMADGE if the stream isn't live
		if len(channelInfo.Data.Streams) == 0 {
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
