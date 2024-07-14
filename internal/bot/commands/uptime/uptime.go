package uptime

import (
	"fmt"
	"time"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/nicklaw5/helix/v2"
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

func (c *UptimeCommand) Description() string {
	return "Get the title of the stream."
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

	res, err := c.gctx.Crate().Helix.Client().GetStreams(&helix.StreamsParams{
		UserIDs: []string{c.gctx.Config().Twitch.Bot.ChannelID},
	})
	if err != nil {
		return "", err
	}

	// Check if the response responded with an unauthorized error or some other error
	if res.Error != "" {
		return fmt.Sprintf("@%v, sorry, the Twitch API threw an error... Susge", user.Name), nil
	}

	if len(res.Data.Streams) == 0 {
		return fmt.Sprintf("@%v, the stream is offline Sadge", target), nil
	}

	return fmt.Sprintf("@%v, the stream has been live for %v", target, utils.TimeDifference(res.Data.Streams[0].StartedAt, time.Now(), true)), nil
}
