package subage

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/sling"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

type SubageCommand struct {
	gctx global.Context
}

func NewSubageCommand(gctx global.Context) *SubageCommand {
	return &SubageCommand{
		gctx: gctx,
	}
}

func (c *SubageCommand) Name() string {
	return "subage"
}

func (c *SubageCommand) Aliases() []string {
	return []string{"sa"}
}

func (c *SubageCommand) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *SubageCommand) Description() string {
	return "Get subage of a user for a specific channel. Defaults to Esfands."
}

func (c *SubageCommand) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"If no channel is given it will default to EsfandTV",
		"<br/>",
		"Check your subage to EsfandTV.",
		fmt.Sprintf("<code>%vsubage</code>", prefix),
		"<br/>",
		"Check subage for another user in Esfand's channel.",
		fmt.Sprintf("<code>%vsubage (username)</code>", prefix),
		"<br/>",
		"Check subage for another user in a specific channel.",
		fmt.Sprintf("<code>%vsubage (username) (channel)</code>", prefix),
	}
}

func (c *SubageCommand) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  true,
		EnabledOffline: true,
	}
}

func (c *SubageCommand) UserCooldown() int {
	return 30
}

func (c *SubageCommand) GlobalCooldown() int {
	return 10
}

func (c *SubageCommand) Code(user twitch.User, context []string) (string, error) {
	// Parse targetUser and targetChannel from context
	targetUser := user.Name
	if len(context) > 0 {
		targetUser = context[0]
	}

	targetChannel := "esfandtv"
	if len(context) > 1 {
		targetChannel = context[1]
	}

	// Remove "@" if it exists in targetUser or targetChannel
	if strings.HasPrefix(targetUser, "@") {
		targetUser = targetUser[1:]
	}

	if strings.HasPrefix(targetChannel, "@") {
		targetChannel = targetChannel[1:]
	}

	// Begin logic to make the API request
	url := "https://api.ivr.fi/"
	s := sling.New().Base(url).Set("Accept", "application/json")
	req, err := s.New().Get(fmt.Sprintf("v2/twitch/subage/%s/%s", targetUser, targetChannel)).Request()
	if err != nil {
		slog.Error("[subage-cmd] error getting subage", "error", err.Error())
		return "Error getting the subage FeelsBadMan", nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("[subage-cmd] error getting subage", "error", err.Error())
		return "Error getting the subage FeelsBadMan", nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("[subage-cmd] error reading subage body", "error", err.Error())
		return "Error getting the subage FeelsBadMan", nil
	}

	var subageRes SubageResponse
	err = json.Unmarshal(body, &subageRes)
	if err != nil {
		slog.Error("[subage-cmd] error unmarshalling subage", "error", err.Error())
		return "Error getting the subage FeelsBadMan", nil
	}

	oldSub := subageRes.Cumulative
	// If they're not subbed...
	if subageRes.Cumulative == nil {
		if oldSub == nil || oldSub.Months == 0 {
			return fmt.Sprintf("%s is not subbed to %s and never has been.", targetUser, targetChannel), nil
		} else {
			parsedOldTime, err := time.Parse(time.RFC3339, oldSub.End)
			if err != nil {
				slog.Error("[subage-cmd] error parsing sub end time", "error", err.Error())
				return "Error parsing the sub end time FeelsBadMan", nil
			}

			return fmt.Sprintf(
				"%s is not subbed to %s but has been previously for a total of %d months. Sub ended %s ago.",
				targetUser, targetChannel, oldSub.Months, utils.TimeDifference(time.Now(), parsedOldTime, true),
			), nil
		}
	} else {
		subData := subageRes.Meta
		subLength := subageRes.Cumulative
		subStreak := subageRes.Streak

		if subData.Tier == "Custom" {
			return fmt.Sprintf(
				"%s is subbed to %s with a permanent sub and has been subbed for a total of %d months! They are currently on a %d months streak.",
				targetUser, targetChannel, subLength.Months, subStreak.Months,
			), nil
		}
		if subData.EndsAt == "" {
			return fmt.Sprintf(
				"%s is currently subbed to %s with a Tier %s sub and has been subbed for a total of %d months! They are currently on a %d months streak. This is a permanent sub.",
				targetUser, targetChannel, subData.Tier, subLength.Months, subStreak.Months,
			), nil
		}
		if subData.Type == "prime" {
			parsedEndTime, err := time.Parse(time.RFC3339, subData.EndsAt)
			if err != nil {
				slog.Error("[subage-cmd] error parsing sub end time", "error", err.Error())
				return "Error parsing the sub end time FeelsBadMan", nil
			}

			return fmt.Sprintf(
				"%s is currently subbed to %s with a prime sub and has been subbed for a total of %d months! They are currently on a %d months streak. The sub ends/renews in %s",
				targetUser, targetChannel, subLength.Months, subStreak.Months, utils.TimeDifference(parsedEndTime, time.Now(), true),
			), nil
		}
		if subData.Type == "paid" {
			parsedEndTime, err := time.Parse(time.RFC3339, subData.EndsAt)
			if err != nil {
				slog.Error("[subage-cmd] error parsing sub end time", "error", err.Error())
				return "Error parsing the sub end time FeelsBadMan", nil
			}

			return fmt.Sprintf(
				"%s is currently subbed to %s with a paid sub and has been subbed for a total of %d months! They are currently on a %d months streak. The sub ends/renews in %s",
				targetUser, targetChannel, subLength.Months, subStreak.Months, utils.TimeDifference(parsedEndTime, time.Now(), true),
			), nil
		}
		if subData.Type == "gift" {
			parsedEndTime, err := time.Parse(time.RFC3339, subData.EndsAt)
			if err != nil {
				slog.Error("[subage-cmd] error parsing sub end time", "error", err.Error())
				return "Error parsing the sub end time FeelsBadMan", nil
			}

			return fmt.Sprintf(
				"%s is currently subbed to %s with a gifted sub from %s and has been subbed for a total of %d months! They are currently on a %d months streak. The sub ends/renews in %s",
				targetUser, targetChannel, subData.GiftMeta.Gifter.DisplayName, subLength.Months, subStreak.Months, utils.TimeDifference(parsedEndTime, time.Now(), true),
			), nil
		}
	}

	return "", nil
}
