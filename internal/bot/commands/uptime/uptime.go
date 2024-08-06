package uptime

import (
	"errors"
	"fmt"
	"time"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
	"golang.org/x/exp/slog"
)

type UptimeCommand struct {
	gctx global.Context
}

func NewUptimeCommand(gctx global.Context) *UptimeCommand {
	return &UptimeCommand{
		gctx: gctx,
	}
}

func (c *UptimeCommand) Name() string {
	return "uptime"
}

func (c *UptimeCommand) Aliases() []string {
	return []string{}
}

func (c *UptimeCommand) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *UptimeCommand) Description() string {
	return "Gets the current time elapsed since the stream started."
}

func (c *UptimeCommand) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Gets the current time elapsed since the stream started.",
		"<br/>",
		fmt.Sprintf("<code>%vuptime</code>", prefix),
	}
}

func (c *UptimeCommand) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  true,
		EnabledOffline: true,
	}
}

func (c *UptimeCommand) UserCooldown() int {
	return 30
}

func (c *UptimeCommand) GlobalCooldown() int {
	return 10
}

func (c *UptimeCommand) Code(user twitch.User, context []string) (string, error) {
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
		parsedEndTime, err := time.Parse(time.RFC3339, "2024-08-05T21:59:23Z")
		if err != nil {
			slog.Error("[uptime-cmd] error parsing stream end time", "error", err.Error())
			return "", errors.New("error parsing the stream end time")
		}

		downtime := utils.TimeDifference(time.Now(), parsedEndTime, true)
		return fmt.Sprintf("@%v, the stream has been offline for %v", target, downtime), nil
	}
}
