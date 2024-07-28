package help

import (
	"fmt"
	"strings"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/gempir/go-twitch-irc/v4"
)

type HelpCommand struct {
	gctx    global.Context
	version string
}

func NewHelpCommand(gctx global.Context, version string) *HelpCommand {
	return &HelpCommand{
		gctx:    gctx,
		version: version,
	}
}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Aliases() []string {
	return []string{"about", "commands"}
}

func (c *HelpCommand) Permissions() []domain.Permission {
	return []domain.Permission{}
}

func (c *HelpCommand) Description() string {
	return "Gives you information about the bot."
}

func (c *HelpCommand) DynamicDescription() []string {
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

func (c *HelpCommand) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  true,
		EnabledOffline: true,
	}
}

func (c *HelpCommand) UserCooldown() int {
	return 30
}

func (c *HelpCommand) GlobalCooldown() int {
	return 10
}

func (c *HelpCommand) Code(user twitch.User, context []string) (string, error) {
	if len(context) >= 1 {
		url := fmt.Sprintf("https://www.retpaladinbot.com/commands/%v", context[0])
		return fmt.Sprintf(`@%v help for the command "%v": %v`, user.Name, strings.ToLower(context[0]), url), nil
	}

	return fmt.Sprintf("@%v, created for EsfandTV and developed by Mahcksimus. Current version: %v, commands: https://www.retpaladinbot.com/", user.Name, c.version), nil
}
