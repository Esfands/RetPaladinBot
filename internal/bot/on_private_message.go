package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/esfands/retpaladinbot/internal/bot/commands"
	"github.com/esfands/retpaladinbot/internal/db"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/gempir/go-twitch-irc/v4"
	"golang.org/x/exp/slog"
)

func (conn *Connection) OnPrivateMessage(gctx global.Context, message twitch.PrivateMessage, commandManager *commands.CommandManager) {
	slog.Debug(fmt.Sprintf("[%v] %v: %v", message.Channel, message.User.DisplayName, message.Message))

	stringID, err := strconv.Atoi(message.User.ID)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// Insert chatter into the database
	err = gctx.Crate().Turso.Queries().InsertChatter(gctx, db.Chatter{
		TID:         stringID,
		Username:    message.User.Name,
		DisplayName: message.User.DisplayName,
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	response, err := handleCommand(gctx, commandManager, message.User, message.Message)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	conn.client.Say(message.Channel, response)
}

func handleCommand(gctx global.Context, commandManager *commands.CommandManager, user twitch.User, msg string) (string, error) {
	if msg[0:len(gctx.Config().Twitch.Bot.Prefix)] == gctx.Config().Twitch.Bot.Prefix {
		msg = msg[len(gctx.Config().Twitch.Bot.Prefix):]

		context := strings.Split(msg, " ")

		for _, dc := range commandManager.DefaultCommands {
			// Found default command by name
			if context[0] == dc.Name() {
				response, err := dc.Code(user, context[1:])
				if err != nil {
					slog.Error(err.Error())
					return "", err
				}

				// Apply cooldown
				isOnCooldown := utils.CooldownCanContinue(user, strings.ToLower(context[0]), dc.UserCooldown(), dc.GlobalCooldown())
				if !isOnCooldown {
					return "", nil
				}

				return response, nil

			} else {
				for _, alias := range dc.Aliases() {
					// Found default command by alias
					if context[0] == alias {
						response, err := dc.Code(user, context[1:])
						if err != nil {
							slog.Error(err.Error())
							return "", err
						}

						// Apply cooldown
						isOnCooldown := utils.CooldownCanContinue(user, strings.ToLower(context[0]), dc.UserCooldown(), dc.GlobalCooldown())
						if !isOnCooldown {
							return "", nil
						}

						return response, nil
					}
				}
			}
		}
	}

	return "", nil
}
