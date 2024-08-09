package commands

import (
	"context"
	"errors"
	"log/slog"

	"github.com/esfands/retpaladinbot/internal/bot/commands/accountage"
	"github.com/esfands/retpaladinbot/internal/bot/commands/command"
	"github.com/esfands/retpaladinbot/internal/bot/commands/dadjoke"
	"github.com/esfands/retpaladinbot/internal/bot/commands/game"
	"github.com/esfands/retpaladinbot/internal/bot/commands/gdq"
	"github.com/esfands/retpaladinbot/internal/bot/commands/help"
	"github.com/esfands/retpaladinbot/internal/bot/commands/ping"
	"github.com/esfands/retpaladinbot/internal/bot/commands/song"
	"github.com/esfands/retpaladinbot/internal/bot/commands/subage"
	"github.com/esfands/retpaladinbot/internal/bot/commands/time"
	"github.com/esfands/retpaladinbot/internal/bot/commands/title"
	"github.com/esfands/retpaladinbot/internal/bot/commands/uptime"
	"github.com/esfands/retpaladinbot/internal/cmdmanager"
	"github.com/esfands/retpaladinbot/internal/db"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/utils"
)

type CommandManager struct {
	gctx    global.Context
	version string

	DefaultCommands []domain.DefaultCommand
	CustomCommands  []domain.CustomCommand
}

func NewCommandManager(gctx global.Context, version string) *CommandManager {
	cm := &CommandManager{
		gctx:    gctx,
		version: version,
	}

	// Load the default commands locally and then save them to database
	cm.DefaultCommands = cm.loadDefaultCommands()
	cm.saveDefaultCommands()

	// Load the custom commands
	err := cm.loadCustomCommands()
	if err != nil {
		slog.Error("Failed to load custom commands", "error", err)
	}

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
		uptime.NewUptimeCommand(cm.gctx),
		dadjoke.NewDadJokeCommand(cm.gctx),
		help.NewHelpCommand(cm.gctx, cm.version),
		command.NewCommandCommand(cm.gctx, cm),
		gdq.NewGDQCommand(cm.gctx),
		subage.NewSubageCommand(cm.gctx),
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
			Aliases:            utils.ConvertSliceToJSONString(dc.Aliases()),
			Permissions:        utils.ConvertSliceToJSONString(utils.ConvertPermissionsToStrings(dc.Permissions())),
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

func (cm *CommandManager) loadCustomCommands() error {
	commands, err := cm.gctx.Crate().Turso.Queries().GetAllCustomCommands(cm.gctx)
	if err != nil {
		return err
	}

	for _, customCommand := range commands {
		cm.CustomCommands = append(cm.CustomCommands, domain.CustomCommand{
			Name:     customCommand.Name,
			Response: customCommand.Response,
		})
	}

	return nil
}

func (cm *CommandManager) AddCustomCommand(cmd domain.CustomCommand) error {
	if cm.CustomCommandExists(cmd.Name) {
		return errors.New("command already exists")
	}
	cm.CustomCommands = append(cm.CustomCommands, cmd)
	// Insert into database
	err := cm.gctx.Crate().Turso.Queries().InsertCustomCommand(context.Background(), db.CustomCommand{
		Name:       cmd.Name,
		Response:   cmd.Response,
		UsageCount: 0,
	})
	return err
}

func (cm *CommandManager) UpdateCustomCommand(cmd domain.CustomCommand) error {
	for i, existingCmd := range cm.CustomCommands {
		if existingCmd.Name == cmd.Name {
			cm.CustomCommands[i].Response = cmd.Response
			// Update in database
			err := cm.gctx.Crate().Turso.Queries().UpdateCustomCommand(context.Background(), db.CustomCommand{
				Name:     cmd.Name,
				Response: cmd.Response,
			})
			return err
		}
	}
	return errors.New("command does not exist")
}

func (cm *CommandManager) DeleteCustomCommand(name string) error {
	for i, cmd := range cm.CustomCommands {
		if cmd.Name == name {
			cm.CustomCommands = append(cm.CustomCommands[:i], cm.CustomCommands[i+1:]...)
			// Delete from database
			err := cm.gctx.Crate().Turso.Queries().DeleteCustomCommand(context.Background(), name)
			return err
		}
	}
	return errors.New("command does not exist")
}

func (cm *CommandManager) CustomCommandExists(name string) bool {
	for _, cmd := range cm.CustomCommands {
		if cmd.Name == name {
			return true
		}
	}
	return false
}

func (cm *CommandManager) GetCustomCommands() []domain.CustomCommand {
	return cm.CustomCommands
}

// Ensure CommandManager implements CommandManagerInterface
var _ cmdmanager.CommandManagerInterface = (*CommandManager)(nil)
