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

	// Initialize the Twitch client with detailed logging
	conn.client = twitch.NewClient(cfg.Twitch.Bot.Username, cfg.Twitch.Bot.OAuth)
	slog.Info("Twitch client initialized", "username", cfg.Twitch.Bot.Username)

	// Register variables service
	conn.Variables = variables.NewService(gctx)
	if conn.Variables == nil {
		slog.Error("Failed to initialize variables service")
		return
	}

	// Setup ModuleManager with error logging
	conn.ModuleManager, err = modules.NewModuleManager(gctx, conn.client)
	if err != nil {
		slog.Error("Error setting up bot modules", "error", err.Error())
		return
	}
	slog.Info("ModuleManager setup complete")

	// Setup CommandManager
	commandManager := commands.NewCommandManager(gctx, version)
	if commandManager == nil {
		slog.Error("Failed to initialize CommandManager")
		return
	}
	slog.Info("CommandManager setup complete")

	// Register message handlers with additional logging
	conn.client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		conn.OnPrivateMessage(gctx, message, commandManager, conn.Variables)
	})
	conn.client.OnUserNoticeMessage(func(message twitch.UserNoticeMessage) {
		OnUserNoticeMessage(conn.client, message)
	})

	// Attempt to join the specified channel with error handling
	slog.Info("Attempting to join channel", "channel", cfg.Twitch.Bot.Channel)
	conn.client.Join(cfg.Twitch.Bot.Channel)
	slog.Info("Successfully joined channel", "channel", cfg.Twitch.Bot.Channel)

	// Graceful shutdown handling
	go func() {
		<-gctx.Done()
		slog.Info("Twitch bot shutting down...")
		if err := conn.client.Disconnect(); err != nil {
			slog.Error("Error disconnecting Twitch client", "error", err.Error())
		} else {
			slog.Info("Twitch bot disconnected successfully")
		}
	}()

	// Attempt to connect to Twitch with enhanced error handling
	slog.Info("Attempting to connect to Twitch")
	if err := conn.client.Connect(); err != nil {
		slog.Error("Failed to connect to Twitch", "error", err.Error())
		return
	}
	slog.Info("Connected to Twitch successfully")
}
