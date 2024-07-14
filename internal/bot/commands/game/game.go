package game

import (
	"fmt"
	"strings"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/nicklaw5/helix/v2"
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

func (c *GameCommand) Description() string {
	return "Get the title of the stream."
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

	res, err := c.gctx.Crate().Helix.Client().GetChannelInformation(&helix.GetChannelInformationParams{
		BroadcasterIDs: []string{c.gctx.Config().Twitch.Bot.ChannelID},
	})
	if err != nil {
		return "", err
	}

	// Check if the response responded with an unauthorized error or some other error
	if res.Error != "" {
		return fmt.Sprintf("@%v, sorry, the Twitch API threw an error... Susge", user.Name), nil
	}

	if strings.ToLower(res.Data.Channels[0].GameName) == "just chatting" {
		return fmt.Sprintf("@%v, Esfand is under the category: %v", target, res.Data.Channels[0].Title), nil
	}

	return fmt.Sprintf("@%v, Esfand is playing %v", target, res.Data.Channels[0].GameName), nil
}
