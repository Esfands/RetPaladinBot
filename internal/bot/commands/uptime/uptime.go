package uptime

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

type Command struct {
	gctx global.Context
}

func NewUptimeCommand(gctx global.Context) *Command {
	return &Command{
		gctx: gctx,
	}
}

func (c *Command) Name() string {
	return "uptime"
}

func (c *Command) Aliases() []string {
	return []string{}
}

func (c *Command) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *Command) Description() string {
	return "Gets the current time elapsed since the stream started."
}

func (c *Command) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Gets the current time elapsed since the stream started.",
		"<br/>",
		fmt.Sprintf("<code>%vuptime</code>", prefix),
	}
}

func (c *Command) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  true,
		EnabledOffline: true,
	}
}

func (c *Command) UserCooldown() int {
	return 30
}

func (c *Command) GlobalCooldown() int {
	return 10
}

func (c *Command) Code(user twitch.User, context []string) (string, error) {
	target := utils.GetTarget(user, context)

	stream, err := c.gctx.Crate().Turso.Queries().GetMostRecentStreamStatus(c.gctx)
	if err != nil {
		slog.Error("[uptime-cmd] error getting most recent stream status", "error", err.Error())
		return "", err
	}

	// Get the uptime since there's no end time
	if !stream.EndedAt.Valid {
		// Parse the start time
		parsedStartTime, err := time.Parse(time.RFC3339, stream.StartedAt)
		if err != nil {
			slog.Error("[uptime-cmd] error parsing stream start time", "error", err.Error())
			return "", errors.New("error parsing the stream start time")
		}

		uptime := utils.TimeDifference(parsedStartTime, time.Now(), true)
		return fmt.Sprintf("@%v, the stream has been live for %v", target, uptime), nil
	} else {
		parsedEndTime, err := time.Parse(time.RFC3339, stream.EndedAt.String)
		if err != nil {
			slog.Error("[uptime-cmd] error parsing stream end time", "error", err.Error())
			return "", errors.New("error parsing the stream end time")
		}

		downtime := utils.TimeDifference(time.Now(), parsedEndTime, true)
		return fmt.Sprintf("@%v, the stream has been offline for %v", target, downtime), nil
	}
}
