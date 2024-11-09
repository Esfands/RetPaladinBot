package ping

import (
	"fmt"
	"time"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

type Command struct {
	gctx global.Context
}

func NewPingCommand(gctx global.Context) *Command {
	return &Command{
		gctx: gctx,
	}
}

func (c *Command) Name() string {
	return "ping"
}

func (c *Command) Aliases() []string {
	return []string{}
}

func (c *Command) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *Command) Description() string {
	return "Ping the bot."
}

func (c *Command) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Pings the bot and returns the uptime.",
		"<br/>",
		fmt.Sprintf("<code>%vping</code>", prefix),
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
	uptime := utils.TimeDifference(c.gctx.Config().Timestamp, time.Now(), true)

	return fmt.Sprintf("@%v, FeelsOkayMan 🏓 Uptime: %v", user.Name, uptime), nil
}
