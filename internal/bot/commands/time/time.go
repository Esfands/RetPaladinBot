package time

import (
	"fmt"
	"time"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

type TimeCommand struct {
	gctx global.Context
}

func NewTimeCommand(gctx global.Context) *TimeCommand {
	cmd := &TimeCommand{
		gctx: gctx,
	}

	return cmd
}

func (c *TimeCommand) Name() string {
	return "time"
}

func (c *TimeCommand) Aliases() []string {
	return []string{}
}

func (c *TimeCommand) Description() string {
	return "Returns the current time of Esfand in CST."
}

func (c *TimeCommand) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  true,
		EnabledOffline: true,
	}
}

func (c *TimeCommand) UserCooldown() int {
	return 10
}

func (c *TimeCommand) GlobalCooldown() int {
	return 30
}

func (c *TimeCommand) Code(user twitch.User, context []string) (string, error) {
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
