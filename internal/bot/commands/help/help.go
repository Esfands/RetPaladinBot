package help

import (
	"fmt"

	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
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
	target := utils.GetTarget(user, context)

	return fmt.Sprintf("@%v, RetPaladinBot was created for EsfandTV and developed by Mahcksimus. Current version: %v", target, c.version), nil
}
