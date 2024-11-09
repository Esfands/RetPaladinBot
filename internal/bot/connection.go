package bot

import (
	"log/slog"

	"github.com/esfands/retpaladinbot/config"
	"github.com/esfands/retpaladinbot/internal/bot/commands"
	"github.com/esfands/retpaladinbot/internal/bot/modules"
	"github.com/esfands/retpaladinbot/internal/bot/variables"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/gempir/go-twitch-irc/v4"
)

type Connection struct {
	client         *twitch.Client
	CommandManager *commands.CommandManager
	ModuleManager  *modules.ModuleManager
	Variables      variables.ServiceI
}

func StartBot(gctx global.Context, cfg *config.Config, version string) {
	conn := &Connection{}
	var err error

	conn.client = twitch.NewClient(cfg.Twitch.Bot.Username, cfg.Twitch.Bot.OAuth)

	// Register variables
	conn.Variables = variables.NewService(gctx)

	conn.ModuleManager, err = modules.NewModuleManager(gctx, conn.client)
	if err != nil {
		slog.Error("Error setting up bot modules", "error", err.Error())
		return
	}
	commandManger := commands.NewCommandManager(gctx, version)

	conn.client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		conn.OnPrivateMessage(gctx, message, commandManger, conn.Variables)
	})

	conn.client.OnUserNoticeMessage(func(message twitch.UserNoticeMessage) {
		OnUserNoticeMessage(conn.client, message)
	})

	conn.client.Join(cfg.Twitch.Bot.Channel)

	go func() {
		<-gctx.Done()

		slog.Info("Twitch bot shutting down...")
		err := conn.client.Disconnect()
		if err != nil {
			return
		}
	}()

	err = conn.client.Connect()
	if err != nil {
		return
	}
}
