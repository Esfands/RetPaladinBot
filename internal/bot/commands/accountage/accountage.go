package accountage

import (
	"fmt"
	"time"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/nicklaw5/helix/v2"
	"golang.org/x/exp/slog"
)

type AccountAgeCommand struct {
	gctx global.Context
}

func NewAccountAgeCommand(gctx global.Context) *AccountAgeCommand {
	return &AccountAgeCommand{
		gctx: gctx,
	}
}

func (c *AccountAgeCommand) Name() string {
	return "accountage"
}

func (c *AccountAgeCommand) Aliases() []string {
	return []string{}
}

func (c *AccountAgeCommand) Description() string {
	return "Check the age of your account."
}

func (c *AccountAgeCommand) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Checks the age of your Twitch account",
		"<br/>",
		fmt.Sprintf("<code>%vaccountage</code>", prefix),
	}
}

func (c *AccountAgeCommand) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  true,
		EnabledOffline: true,
	}
}

func (c *AccountAgeCommand) UserCooldown() int {
	return 30
}

func (c *AccountAgeCommand) GlobalCooldown() int {
	return 10
}

func (c *AccountAgeCommand) Code(user twitch.User, context []string) (string, error) {
	target := utils.GetTarget(user, context)

	res, err := c.gctx.Crate().Helix.Client().GetUsers(&helix.UsersParams{
		Logins: []string{target},
	})
	if err != nil {
		return "", err
	}

	// Check if the response responded with an unauthorized error or some other error
	if res.Error != "" {
		slog.Error("Twitch API error while fetching account age", "error", res.ErrorMessage)
		return fmt.Sprintf("@%v, sorry, the Twitch API threw an error... Susge", user.Name), nil
	}

	if len(res.Data.Users) == 0 {
		return fmt.Sprintf("@%v, sorry I couldn't find a user with that name!", user.Name), nil
	}

	slog.Debug("Target user test", "target", target)

	elapsed := utils.TimeDifference(res.Data.Users[0].CreatedAt.Time, time.Now(), true)

	if target != user.Name {
		return fmt.Sprintf("@%v created their account %v ago", target, elapsed), nil
	}

	return fmt.Sprintf("@%v, you created your account %v ago", target, elapsed), nil
}
