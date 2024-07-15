package commands

import (
	"github.com/esfands/retpaladinbot/internal/bot/commands/accountage"
	"github.com/esfands/retpaladinbot/internal/bot/commands/dadjoke"
	"github.com/esfands/retpaladinbot/internal/bot/commands/followage"
	"github.com/esfands/retpaladinbot/internal/bot/commands/game"
	"github.com/esfands/retpaladinbot/internal/bot/commands/help"
	"github.com/esfands/retpaladinbot/internal/bot/commands/ping"
	"github.com/esfands/retpaladinbot/internal/bot/commands/song"
	"github.com/esfands/retpaladinbot/internal/bot/commands/time"
	"github.com/esfands/retpaladinbot/internal/bot/commands/title"
	"github.com/esfands/retpaladinbot/internal/bot/commands/uptime"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
)

type CommandManager struct {
	gctx    global.Context
	version string

	DefaultCommands []domain.DefaultCommand
	CustomCommands  []domain.Command
}

func NewCommandManager(gctx global.Context, version string) *CommandManager {
	cm := &CommandManager{
		gctx:    gctx,
		version: version,
	}

	cm.DefaultCommands = cm.loadDefaultCommands()
	// cm.saveDefaultCommands()

	return cm
}

func (cm *CommandManager) loadDefaultCommands() []domain.DefaultCommand {
	return []domain.DefaultCommand{
		ping.NewPingCommand(cm.gctx),
		accountage.NewAccountAgeCommand(cm.gctx),
		song.NewSongCommand(cm.gctx),
		time.NewTimeCommand(cm.gctx),
		title.NewTitleCommand(cm.gctx),
		game.NewGameCommand(cm.gctx),
		followage.NewFollowageCommand(cm.gctx),
		uptime.NewUptimeCommand(cm.gctx),
		dadjoke.NewDadJokeCommand(cm.gctx),
		help.NewHelpCommand(cm.gctx, cm.version),
	}
}
