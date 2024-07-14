package followage

import (
	"fmt"
	"time"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/nicklaw5/helix/v2"
)

type FollowageCommand struct {
	gctx global.Context
}

func NewFollowageCommand(gctx global.Context) *FollowageCommand {
	cmd := &FollowageCommand{
		gctx: gctx,
	}

	return cmd
}

func (c *FollowageCommand) Name() string {
	return "followage"
}

func (c *FollowageCommand) Aliases() []string {
	return []string{}
}

func (c *FollowageCommand) Description() string {
	return "Get your Twitch account age."
}

func (c *FollowageCommand) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  true,
		EnabledOffline: true,
	}
}

func (c *FollowageCommand) UserCooldown() int {
	return 30
}

func (c *FollowageCommand) GlobalCooldown() int {
	return 10
}

func (c *FollowageCommand) Code(user twitch.User, context []string) (string, error) {
	target := utils.GetTarget(user, context)

	var userID string
	if target == user.Name {
		userID = user.ID
	} else {
		res, err := c.gctx.Crate().Helix.Client().GetUsers(&helix.UsersParams{
			Logins: []string{target},
		})
		if err != nil {
			return fmt.Sprintf("@%v, there was an error with Twitch's API", user.Name), err
		}

		if len(res.Data.Users) == 0 {
			return fmt.Sprintf("@%v, %v does not exist", user.Name, target), nil
		}

		userID = res.Data.Users[0].ID
	}

	res, err := c.gctx.Crate().Helix.Client().GetUsersFollows(&helix.UsersFollowsParams{
		ToID:   c.gctx.Config().Twitch.Bot.ChannelID,
		FromID: userID,
	})
	if err != nil {
		return "", err
	}

	if len(res.Data.Follows) == 0 {
		if target == user.Name {
			return fmt.Sprintf("%v, you are not following the channel", user.Name), nil
		}

		return fmt.Sprintf("%v, %v is not following the channel", user.Name, target), nil
	}

	if target == user.Name {
		return fmt.Sprintf("@%v, you have been following for %v", user.Name, utils.TimeDifference(res.Data.Follows[0].FollowedAt, time.Now(), true)), nil
	}

	return fmt.Sprintf("@%v, %v has been following for %v", user.Name, target, utils.TimeDifference(res.Data.Follows[0].FollowedAt, time.Now(), true)), nil
}
