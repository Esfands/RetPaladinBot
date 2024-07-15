package commands

import (
	"strings"

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
	"github.com/esfands/retpaladinbot/internal/db"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
	"golang.org/x/exp/slog"
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
	cm.saveDefaultCommands()

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

func (cm *CommandManager) saveDefaultCommands() {
	storedCommands, err := cm.gctx.Crate().Turso.Queries().GetAllDefaultCommands(cm.gctx)
	if err != nil {
		slog.Error("Failed to get stored default commands", "error", err)
		return
	}

	// Map to keep track of commands in the codebase
	codebaseCommands := make(map[string]db.DefaultCommand)
	for _, dc := range cm.DefaultCommands {
		codebaseCommands[dc.Name()] = db.DefaultCommand{
			Name:               dc.Name(),
			Aliases:            strings.Join(dc.Aliases(), ","),
			Permissions:        "",
			Description:        dc.Description(),
			DynamicDescription: utils.ConvertSliceToJSONString(dc.DynamicDescription()),
			GlobalCooldown:     dc.GlobalCooldown(),
			UserCooldown:       dc.UserCooldown(),
			EnabledOffline:     utils.BoolToInt(dc.Conditions().EnabledOffline),
			EnabledOnline:      utils.BoolToInt(dc.Conditions().EnabledOnline),
			UsageCount:         0,
		}
	}

	// Check for commands to update or remove
	for _, storedCommand := range storedCommands {
		if _, exists := codebaseCommands[storedCommand.Name]; !exists {
			// Command exists in database but not in the codebase, remove it
			err = cm.gctx.Crate().Turso.Queries().DeleteDefaultCommand(cm.gctx, storedCommand.Name)
			if err != nil {
				slog.Error("Failed to delete default command", "command", storedCommand.Name, "error", err)
			}
		} else {
			// Command exists in both codebase and database, update it if necessary
			codebaseCommand := codebaseCommands[storedCommand.Name]
			err = cm.gctx.Crate().Turso.Queries().UpdateDefaultCommand(cm.gctx, codebaseCommand)
			if err != nil {
				slog.Error("Failed to update default command", "command", storedCommand.Name, "error", err)
			}
			// Remove from codebaseCommands map to mark it as processed
			delete(codebaseCommands, storedCommand.Name)
		}
	}

	// Add new commands from the codebase to the database
	for _, newCommand := range codebaseCommands {
		err = cm.gctx.Crate().Turso.Queries().InsertDefaultCommand(cm.gctx, newCommand)
		if err != nil {
			slog.Error("Failed to insert default command", "command", newCommand.Name, "error", err)
		}
	}
}
