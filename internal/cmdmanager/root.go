package cmdmanager

import "github.com/esfands/retpaladinbot/pkg/domain"

// CommandManagerInterface defines the methods for managing commands.
type CommandManagerInterface interface {
	AddCustomCommand(cmd domain.CustomCommand) error
	UpdateCustomCommand(cmd domain.CustomCommand) error
	DeleteCustomCommand(name string) error
	CustomCommandExists(name string) bool
	GetCustomCommands() []domain.CustomCommand
}
