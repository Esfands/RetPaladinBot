package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/esfands/retpaladinbot/internal/cmdmanager"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/gempir/go-twitch-irc/v4"
)

type CommandCommand struct {
	gctx    global.Context
	manager cmdmanager.CommandManagerInterface
}

func NewCommandCommand(gctx global.Context, manager cmdmanager.CommandManagerInterface) *CommandCommand {
	return &CommandCommand{
		gctx:    gctx,
		manager: manager,
	}
}

func (c *CommandCommand) Name() string {
	return "command"
}

func (c *CommandCommand) Aliases() []string {
	return []string{
		"cmd",
	}
}

func (c *CommandCommand) Permissions() []domain.Permission {
	return []domain.Permission{
		domain.PermissionBroadcaster,
		domain.PermissionModerator,
	}
}

func (c *CommandCommand) Description() string {
	return "Create/edit/delete custom commands."
}

func (c *CommandCommand) DynamicDescription() []string {
	prefix := c.gctx.Config().Twitch.Bot.Prefix

	return []string{
		"Create a command with a name and message",
		"<br/>",
		fmt.Sprintf("<code>%vcommand create (name) (response)</code>", prefix),
		"<br/>",
		"Edit a command with a name and message",
		"<br/>",
		fmt.Sprintf("<code>%vcommand edit (name) (response)</code>", prefix),
		"<br/>",
		"Delete a command with a name",
		"<br/>",
		fmt.Sprintf("<code>%vcommand delete (name)</code>", prefix),
	}
}

func (c *CommandCommand) Conditions() domain.DefaultCommandConditions {
	return domain.DefaultCommandConditions{
		EnabledOnline:  true,
		EnabledOffline: true,
	}
}

func (c *CommandCommand) UserCooldown() int {
	return 30
}

func (c *CommandCommand) GlobalCooldown() int {
	return 10
}

func (c *CommandCommand) Code(user twitch.User, context []string) (string, error) {
	if len(context) < 2 {
		return "", errors.New("invalid command format")
	}

	action := context[0]
	name := context[1]

	// Remove the "!" prefix if it exists in the command name
	name = strings.TrimPrefix(name, "!")

	// Lowercase the command name
	name = strings.ToLower(name)

	response := strings.Join(context[2:], " ")

	switch action {
	case "create":
		return c.createCommand(name, response)
	case "edit":
		return c.editCommand(name, response)
	case "delete":
		return c.deleteCommand(name)
	default:
		return "", errors.New("invalid action specified")
	}
}

func (c *CommandCommand) createCommand(name, response string) (string, error) {
	// Check if the command already exists
	if c.manager.CustomCommandExists(name) {
		return "", errors.New("command already exists")
	}

	// Add the new command to the manager's CustomCommands slice
	err := c.manager.AddCustomCommand(domain.CustomCommand{
		Name:     name,
		Response: response,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Command '%s' created with response: %s", name, response), nil
}

func (c *CommandCommand) editCommand(name, response string) (string, error) {
	// Update the command's response
	err := c.manager.UpdateCustomCommand(domain.CustomCommand{
		Name:     name,
		Response: response,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Command '%s' updated with new response: %s", name, response), nil
}

func (c *CommandCommand) deleteCommand(name string) (string, error) {
	// Delete the command from the manager's CustomCommands slice
	err := c.manager.DeleteCustomCommand(name)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Command '%s' deleted", name), nil
}
