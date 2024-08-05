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

type GameCommand struct {
	gctx global.Context
}

func NewGameCommand(gctx global.Context) *GameCommand {
	return &GameCommand{
		gctx: gctx,
	}
}

func (c *GameCommand) Name() string {
	return "game"
}

func (c *GameCommand) Aliases() []string {
	return []string{"category"}
}

func (c *GameCommand) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *GameCommand) Description() string {
	return "Get the current category of the stream."
}

func (c *GameCommand) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Gives the current category of the stream",
		"<br/>",
		fmt.Sprintf("<code>%vgame</code>", prefix),
		"<br/>",
		fmt.Sprintf("<code>%vcategory</code>", prefix),
	}
}

func (c *GameCommand) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  true,
		EnabledOffline: true,
	}
}

func (c *GameCommand) UserCooldown() int {
	return 30
}

func (c *GameCommand) GlobalCooldown() int {
	return 10
}

func (c *GameCommand) Code(user twitch.User, context []string) (string, error) {
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
