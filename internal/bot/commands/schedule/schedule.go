package schedule

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dghubble/sling"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
	"golang.org/x/exp/slog"
)

type ScheduleCommand struct {
	gctx global.Context
}

func NewScheduleCommand(gctx global.Context) *ScheduleCommand {
	return &ScheduleCommand{
		gctx: gctx,
	}
}

func (c *ScheduleCommand) Name() string {
	return "schedule"
}

func (c *ScheduleCommand) Aliases() []string {
	return []string{}
}

func (c *ScheduleCommand) Description() string {
	return "Get the next scheduled stream"
}

func (c *ScheduleCommand) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  true,
		EnabledOffline: true,
	}
}

func (c *ScheduleCommand) UserCooldown() int {
	return 10
}

func (c *ScheduleCommand) GlobalCooldown() int {
	return 30
}

func (c *ScheduleCommand) Code(user twitch.User, context []string) (string, error) {
	target := utils.GetTarget(user, context)

	schedule, err := fetchSchedule()
	if err != nil {
		slog.Error("Failed to fetch schedule", "error", err)
		return fmt.Sprintf("@%v failed to fetch the schedule. FeelsBadMan", user.Name), err
	}

	message := getScheduleMessage(schedule)
	return fmt.Sprintf("@%v %v", target, message), nil
}

func fetchSchedule() (*Schedule, error) {
	url := "https://twitch.otkdata.com/api/streamers/esfandtv/schedule"

	s := sling.New().Base(url).Set("Accept", "application/json")
	req, err := s.New().Get("").Request()
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get schedule: status code %d, response %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read schedule response body: %w", err)
	}

	// Log the raw response body for debugging
	slog.Info("Raw schedule response body", "body", string(body))

	var schedule Schedule
	if err := json.Unmarshal(body, &schedule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schedule response: %w", err)
	}

	return &schedule, nil
}

func getScheduleMessage(schedule *Schedule) string {
	now := time.Now()
	nextSegment := getNextSegment(schedule.Data.Segments)

	if nextSegment == nil {
		return "No upcoming streams scheduled."
	}

	if nextSegment.StartTime.After(now) {
		return fmt.Sprintf("The next stream is scheduled for %s: %s", nextSegment.StartTime.Format(time.RFC1123), nextSegment.Title)
	}

	if nextSegment.EndTime == nil || nextSegment.EndTime.IsZero() {
		return fmt.Sprintf("The stream titled '%s' should have started at %s but hasn't yet.", nextSegment.Title, nextSegment.StartTime.Format(time.RFC1123))
	}

	if now.After(nextSegment.StartTime.Time) && now.Before(nextSegment.EndTime.Time) {
		return fmt.Sprintf("The stream titled '%s' is currently live.", nextSegment.Title)
	}

	return fmt.Sprintf("The stream titled '%s' should have started at %s but hasn't yet.", nextSegment.Title, nextSegment.StartTime.Format(time.RFC1123))
}

func getNextSegment(segments []Segment) *Segment {
	now := time.Now()
	for _, segment := range segments {
		if segment.StartTime.After(now) {
			return &segment
		}
	}
	return nil
}
