package modules

import (
	"github.com/esfands/retpaladinbot/internal/bot/modules/emotes"
	goliverightnow "github.com/esfands/retpaladinbot/internal/bot/modules/go-live-right-now"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/gempir/go-twitch-irc/v4"
)

type ModuleManager struct {
	Emotes         *emotes.EmoteModule
	GoLiveRightNow *goliverightnow.GoLiveRightNowModule
}

func NewModuleManager(gctx global.Context, client *twitch.Client, channelID string) (*ModuleManager, error) {
	mm := &ModuleManager{}
	var err error

	mm.Emotes, err = emotes.NewEmoteModule(gctx, channelID)
	if err != nil {
		return nil, err
	}

	mm.GoLiveRightNow = goliverightnow.NewGoLiveRightNowModule(gctx, client)

	return mm, nil
}
