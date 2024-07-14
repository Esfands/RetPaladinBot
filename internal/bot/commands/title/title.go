package title

import (
	"fmt"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/nicklaw5/helix/v2"
)

type TitleCommand struct {
	gctx global.Context
}

func NewTitleCommand(gctx global.Context) *TitleCommand {
	return &TitleCommand{
		gctx: gctx,
	}
}

func (c *TitleCommand) Name() string {
	return "title"
}

func (c *TitleCommand) Aliases() []string {
	return []string{}
}

func (c *TitleCommand) Description() string {
	return "Get the title of the stream."
}

func (c *TitleCommand) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  true,
		EnabledOffline: true,
	}
}

func (c *TitleCommand) UserCooldown() int {
	return 30
}

func (c *TitleCommand) GlobalCooldown() int {
	return 10
}

func (c *TitleCommand) Code(user twitch.User, context []string) (string, error) {
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

	return fmt.Sprintf("@%v current title: %v", target, res.Data.Channels[0].Title), nil
}
