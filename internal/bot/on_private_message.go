package bot

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/esfands/retpaladinbot/internal/bot/commands"
	"github.com/esfands/retpaladinbot/internal/bot/variables"
	"github.com/esfands/retpaladinbot/internal/db"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
)

func (conn *Connection) OnPrivateMessage(gctx global.Context, message twitch.PrivateMessage, commandManager *commands.CommandManager, variables variables.ServiceI) {
	slog.Debug(fmt.Sprintf("[%v] %v: %v", message.Channel, message.User.DisplayName, message.Message))

	stringID, err := strconv.Atoi(message.User.ID)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// Insert chatter into the database
	_ = gctx.Crate().Turso.Queries().InsertChatter(gctx, db.Chatter{
		TID:         stringID,
		Username:    message.User.Name,
		DisplayName: message.User.DisplayName,
	})

	response, err := handleCommand(gctx, variables, commandManager, message.User, message.Message)
	if err != nil {
		slog.Error(err.Error())
		conn.client.Say(message.Channel, fmt.Sprintf("Something went wrong... error: %v", err.Error()))
		return
	}

	conn.client.Say(message.Channel, response)
}

// Map Twitch badges to domain permissions
var badgeToPermission = map[string]domain.Permission{
	"broadcaster": domain.PermissionBroadcaster,
	"moderator":   domain.PermissionModerator,
	"vip":         domain.PermissionVIP,
}

func isUserPermitted(user twitch.User, requiredPermissions []domain.Permission) bool {
	for badge := range user.Badges {
		if permission, exists := badgeToPermission[badge]; exists {
			for _, requiredPermission := range requiredPermissions {
				if permission == requiredPermission {
					return true
				}
			}
		}
	}
	return false
}

// Check if the command name or alias matches the input
func isCommandMatch(input string, command domain.DefaultCommand) bool {
	if input == command.Name() {
		return true
	}
	for _, alias := range command.Aliases() {
		if input == alias {
			return true
		}
	}
	return false
}

// Execute command with permission and cooldown checks
func executeCommand(gctx global.Context, user twitch.User, context []string, command domain.DefaultCommand) (string, error) {
	// Allow execution if the command has no required permissions
	if len(command.Permissions()) > 0 && !isUserPermitted(user, command.Permissions()) {
		return "", nil
	}

	response, err := command.Code(user, context[1:])
	if err != nil {
		return "", err
	}

	if !utils.CooldownCanContinue(user, strings.ToLower(context[0]), command.UserCooldown(), command.GlobalCooldown()) {
		return "", nil
	}

	// Update usage
	err = gctx.Crate().Turso.Queries().IncrementDefaultCommandUsageCount(gctx, command.Name())
	if err != nil {
		slog.Error("Failed to update default command usage", "error", err.Error())
	}

	slog.Info("Command executed", "command", command.Name(), "user", user.DisplayName, "channel", gctx.Config().Twitch.Bot.Channel)

	return response, nil
}

func handleCommand(gctx global.Context, variables variables.ServiceI, commandManager *commands.CommandManager, user twitch.User, msg string) (string, error) {
	if !strings.HasPrefix(msg, gctx.Config().Twitch.Bot.Prefix) {
		return "", nil
	}

	msg = strings.TrimPrefix(msg, gctx.Config().Twitch.Bot.Prefix)
	context := strings.Split(msg, " ")

	for _, dc := range commandManager.DefaultCommands {
		if isCommandMatch(strings.ToLower(context[0]), dc) {
			slog.Info("Command match found", "command", dc.Name())

			streamStatus, err := gctx.Crate().Turso.Queries().GetMostRecentStreamStatus(gctx)
			if err != nil {
				slog.Error("Failed to get most recent stream status", "error", err.Error())
				return "", err
			}
			slog.Info("Stream status", "live", streamStatus.Live)

			if dc.Conditions().EnabledOffline && !streamStatus.Live {
				// Command can run offline, and stream is offline
				return executeCommand(gctx, user, context, dc)
			} else if dc.Conditions().EnabledOnline && streamStatus.Live {
				// Command can run online, and stream is live
				return executeCommand(gctx, user, context, dc)
			} else {
				// Command does not meet the conditions to run
				slog.Info("Command cannot run in the current stream status", "command", dc.Name(), "live", streamStatus.Live)
				return "", nil
			}
		}
	}

	// Check for custom commands
	for _, cc := range commandManager.CustomCommands {
		if strings.ToLower(context[0]) == cc.Name {
			// Bypass the cooldown for broadcaster and moderator
			if user.Badges["broadcaster"] == 1 || user.Badges["moderator"] == 1 {
				err := gctx.Crate().Turso.Queries().IncrementCustomCommandUsageCount(gctx, strings.ToLower(cc.Name))
				if err != nil {
					slog.Error("Failed to update custom command usage", "error", err.Error())
				}

				parsedResponse := variables.ParseVariables(gctx, user, strings.Split(msg, " "), strings.Split(cc.Response, " "))

				return parsedResponse, nil
			}

			if !utils.CooldownCanContinue(user, strings.ToLower(cc.Name), 30, 10) {
				return "", nil
			}

			err := gctx.Crate().Turso.Queries().IncrementCustomCommandUsageCount(gctx, strings.ToLower(cc.Name))
			if err != nil {
				slog.Error("Failed to update custom command usage", "error", err.Error())
			}

			parsedResponse := variables.ParseVariables(gctx, user, strings.Split(msg, " "), strings.Split(cc.Response, " "))

			return parsedResponse, nil
		}
	}

	return "", nil
}
