package commands

import (
	"github.com/esfands/retpaladinbot/internal/rest/v1/respond"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/errors"
	"github.com/esfands/retpaladinbot/pkg/utils"
)

func convertToPermissions(strings []string) []domain.Permission {
	var permissions []domain.Permission
	for _, str := range strings {
		permissions = append(permissions, domain.Permission(str))
	}
	return permissions
}

type GetCommandsResponse struct {
	DefaultCommands []domain.Command       `json:"default_commands"`
	CustomCommands  []domain.CustomCommand `json:"custom_commands"`
}

func (rg *RouteGroup) GetCommands(ctx *respond.Ctx) error {
	storedDefaultCommands, err := rg.gctx.Crate().Turso.Queries().GetAllDefaultCommands(ctx.Context())
	if err != nil {
		return errors.ErrInternalServerError().SetDetail(err.Error())
	}

	var defaultCommands []domain.Command
	for _, storedDefaultCommand := range storedDefaultCommands {
		convertedAliases, err := utils.ConvertJSONStringToSlice(storedDefaultCommand.Aliases)
		if err != nil {
			return errors.ErrInternalServerError().SetDetail(err.Error())
		}

		convertedDynamicDescription, err := utils.ConvertJSONStringToSlice(storedDefaultCommand.DynamicDescription)
		if err != nil {
			return errors.ErrInternalServerError().SetDetail(err.Error())
		}

		// Retrieve permissions
		convertedPermissions, err := utils.ConvertJSONStringToSlice(storedDefaultCommand.Permissions)
		if err != nil {
			return errors.ErrInternalServerError().SetDetail(err.Error())
		}

		defaultCommands = append(defaultCommands, domain.Command{
			Name:               storedDefaultCommand.Name,
			Aliases:            convertedAliases,
			Description:        storedDefaultCommand.Description,
			DynamicDescription: convertedDynamicDescription,
			GlobalCooldown:     storedDefaultCommand.GlobalCooldown,
			UserCooldown:       storedDefaultCommand.UserCooldown,
			EnabledOffline:     storedDefaultCommand.EnabledOffline == 1,
			EnabledOnline:      storedDefaultCommand.EnabledOnline == 1,
			UsageCount:         storedDefaultCommand.UsageCount,
			Permissions:        convertToPermissions(convertedPermissions),
		})
	}

	storedCustomCommands, err := rg.gctx.Crate().Turso.Queries().GetAllCustomCommands(ctx.Context())
	if err != nil {
		return errors.ErrInternalServerError().SetDetail(err.Error())
	}

	var customCommands []domain.CustomCommand
	for _, storedCustomCommand := range storedCustomCommands {
		customCommands = append(customCommands, domain.CustomCommand{
			Name:       storedCustomCommand.Name,
			Response:   storedCustomCommand.Response,
			UsageCount: storedCustomCommand.UsageCount,
		})
	}

	return ctx.JSON(GetCommandsResponse{
		DefaultCommands: defaultCommands,
		CustomCommands:  customCommands,
	})
}

type GetCommandResponse struct {
	DefaultCommand *domain.Command       `json:"default_command,omitempty"`
	CustomCommand  *domain.CustomCommand `json:"custom_command,omitempty"`
}

func (rg *RouteGroup) GetCommandByName(ctx *respond.Ctx) error {
	name := ctx.Params("name")

	// Query the default commands
	storedDefaultCommand, err := rg.gctx.Crate().Turso.Queries().GetDefaultCommandByName(ctx.Context(), name)
	if err == nil && storedDefaultCommand != nil {
		convertedAliases, err := utils.ConvertJSONStringToSlice(storedDefaultCommand.Aliases)
		if err != nil {
			return errors.ErrInternalServerError().SetDetail(err.Error())
		}

		convertedDynamicDescription, err := utils.ConvertJSONStringToSlice(storedDefaultCommand.DynamicDescription)
		if err != nil {
			return errors.ErrInternalServerError().SetDetail(err.Error())
		}

		// Retrieve permissions
		convertedPermissions, err := utils.ConvertJSONStringToSlice(storedDefaultCommand.Permissions)
		if err != nil {
			return errors.ErrInternalServerError().SetDetail(err.Error())
		}

		command := domain.Command{
			Name:               storedDefaultCommand.Name,
			Aliases:            convertedAliases,
			Permissions:        convertToPermissions(convertedPermissions),
			Description:        storedDefaultCommand.Description,
			DynamicDescription: convertedDynamicDescription,
			GlobalCooldown:     storedDefaultCommand.GlobalCooldown,
			UserCooldown:       storedDefaultCommand.UserCooldown,
			EnabledOffline:     storedDefaultCommand.EnabledOffline == 1,
			EnabledOnline:      storedDefaultCommand.EnabledOnline == 1,
			UsageCount:         storedDefaultCommand.UsageCount,
		}

		return ctx.JSON(GetCommandResponse{DefaultCommand: &command, CustomCommand: nil})
	}

	// Query the custom commands
	storedCustomCommand, err := rg.gctx.Crate().Turso.Queries().GetCustomCommandByName(ctx.Context(), name)
	if err == nil && storedCustomCommand != nil {
		command := domain.CustomCommand{
			Name:       storedCustomCommand.Name,
			Response:   storedCustomCommand.Response,
			UsageCount: storedCustomCommand.UsageCount,
		}

		return ctx.JSON(GetCommandResponse{DefaultCommand: nil, CustomCommand: &command})
	}

	// If no command was found
	return errors.ErrNotFound().SetDetail("Command not found")
}
