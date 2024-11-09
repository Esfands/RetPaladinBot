package modules

import (
	goliverightnow "github.com/esfands/retpaladinbot/internal/bot/modules/go-live-right-now"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/gempir/go-twitch-irc/v4"
)

type ModuleManager struct {
	GoLiveRightNow *goliverightnow.Module
}

func NewModuleManager(gctx global.Context, client *twitch.Client) (*ModuleManager, error) {
	return &ModuleManager{
		GoLiveRightNow: goliverightnow.NewGoLiveRightNowModule(gctx, client),
	}, nil
}
