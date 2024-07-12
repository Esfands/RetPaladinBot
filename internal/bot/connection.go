package bot

import (
	"github.com/esfands/retpaladinbot/config"
	"github.com/esfands/retpaladinbot/internal/bot/commands"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/gempir/go-twitch-irc/v4"
	"golang.org/x/exp/slog"
)

type Connection struct {
	client *twitch.Client
}

func StartBot(gctx global.Context, cfg *config.Config) {
	conn := &Connection{}

	conn.client = twitch.NewClient(cfg.Twitch.Bot.Username, cfg.Twitch.Bot.OAuth)

	commandManger := commands.NewCommandManager(gctx)

	conn.client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		conn.OnPrivateMessage(gctx, message, commandManger)
	})

	conn.client.Join(cfg.Twitch.Bot.Channel)

	go func() {
		<-gctx.Done()

		slog.Info("Twitch bot shutting down...")
		conn.client.Disconnect()
	}()

	conn.client.Connect()
}
