package game

import (
	"errors"
	"fmt"
	"strings"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

type Command struct {
	gctx global.Context
}

func NewGameCommand(gctx global.Context) *Command {
	return &Command{
		gctx: gctx,
	}
}

func (c *Command) Name() string {
	return "game"
}

func (c *Command) Aliases() []string {
	return []string{"category"}
}

func (c *Command) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *Command) Description() string {
	return "Get the current category of the stream."
}

func (c *Command) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Gives the current category of the stream",
		"<br/>",
		fmt.Sprintf("<code>%vgame</code>", prefix),
		"<br/>",
		fmt.Sprintf("<code>%vcategory</code>", prefix),
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

	if stream.GameName.String == "" || !stream.GameID.Valid {
		return fmt.Sprintf("@%v, Esfand isn't under a specific category", target), nil
	}

	if strings.ToLower(stream.GameName.String) == "just chatting" {
		return fmt.Sprintf("@%v, Esfand is under the category: %v", target, stream.GameName.String), nil
	}

	return fmt.Sprintf("@%v, Esfand is playing %v", target, stream.GameName.String), nil
}
