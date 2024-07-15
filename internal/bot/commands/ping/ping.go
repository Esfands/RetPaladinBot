package ping

import (
	"fmt"
	"time"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

type PingCommand struct {
	gctx global.Context
}

func NewPingCommand(gctx global.Context) *PingCommand {
	return &PingCommand{
		gctx: gctx,
	}
}

func (c *PingCommand) Name() string {
	return "ping"
}

func (c *PingCommand) Aliases() []string {
	return []string{}
}

func (c *PingCommand) Description() string {
	return "Ping the bot."
}

func (c *PingCommand) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Pings the bot and returns the uptime.",
		"<br/>",
		fmt.Sprintf("<code>%vping</code>", prefix),
	}
}

func (c *PingCommand) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  true,
		EnabledOffline: true,
	}
}

func (c *PingCommand) UserCooldown() int {
	return 30
}

func (c *PingCommand) GlobalCooldown() int {
	return 10
}

func (c *PingCommand) Code(user twitch.User, context []string) (string, error) {
	uptime := utils.TimeDifference(c.gctx.Config().Timestamp, time.Now(), true)

	return fmt.Sprintf("@%v, FeelsOkayMan 🏓 Uptime: %v", user.Name, uptime), nil
}
