package commands

import (
	"github.com/esfands/retpaladinbot/internal/rest/v1/respond"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"github.com/esfands/retpaladinbot/pkg/errors"
	"github.com/esfands/retpaladinbot/pkg/utils"
)

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
