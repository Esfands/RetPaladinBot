package help

import (
	"fmt"
	"strings"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/gempir/go-twitch-irc/v4"
)

type Command struct {
	gctx    global.Context
	version string
}

func NewHelpCommand(gctx global.Context, version string) *Command {
	return &Command{
		gctx:    gctx,
		version: version,
	}
}

func (c *Command) Name() string {
	return "help"
}

func (c *Command) Aliases() []string {
	return []string{"about", "commands"}
}

func (c *Command) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *Command) Description() string {
	return "Gives you information about the bot."
}

func (c *Command) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Gives you information about the bot as well as the current version.",
		"<br/>",
		fmt.Sprintf("<code>%vhelp</code>", prefix),
		"<br/>",
		fmt.Sprintf("<code>%vabout</code>", prefix),
		"<br/>",
		fmt.Sprintf("<code>%vcommands</code>", prefix),
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
	if len(context) >= 1 {
		url := fmt.Sprintf("https://www.retpaladinbot.com/commands/%v", context[0])
		return fmt.Sprintf(`@%v help for the command "%v": %v`, user.Name, strings.ToLower(context[0]), url), nil
	}

	return fmt.Sprintf("@%v, created for EsfandTV and developed by Mahcksimus. Current version: %v, commands: https://www.retpaladinbot.com/", user.Name, c.version), nil
}
