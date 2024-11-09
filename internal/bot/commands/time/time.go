package time

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

func NewTimeCommand(gctx global.Context) *Command {
	cmd := &Command{
		gctx: gctx,
	}

	return cmd
}

func (c *Command) Name() string {
	return "time"
}

func (c *Command) Aliases() []string {
	return []string{}
}

func (c *Command) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *Command) Description() string {
	return "Returns the current time of Esfand in CST."
}

func (c *Command) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Gets the current time of Esfands local time in CST and military time.",
		"<br/>",
		fmt.Sprintf("<code>%vtime</code>", prefix),
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

	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	currentTime := time.Now().In(location)

	return fmt.Sprintf(
		"@%v Esfand's local time is %v CST KKona (%v)",
		target, currentTime.Format("03:04 PM"),
		currentTime.Format("15:04"),
	), nil
}
