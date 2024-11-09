package title

import (
	"errors"
	"fmt"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

type Command struct {
	gctx global.Context
}

func NewTitleCommand(gctx global.Context) *Command {
	return &Command{
		gctx: gctx,
	}
}

func (c *Command) Name() string {
	return "title"
}

func (c *Command) Aliases() []string {
	return []string{}
}

func (c *Command) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *Command) Description() string {
	return "Get the title of the stream."
}

func (c *Command) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Gets the current title of the stream.",
		"<br/>",
		fmt.Sprintf("<code>%vtitle</code>", prefix),
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
		return "", errors.New("error getting the stream status")
	}

	if stream.Title.String == "" || !stream.Title.Valid {
		return fmt.Sprintf("@%v the title is not set to anything", target), nil
	}

	return fmt.Sprintf("@%v current title: %v", target, stream.Title.String), nil
}
