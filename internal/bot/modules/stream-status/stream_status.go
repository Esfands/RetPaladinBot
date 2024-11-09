package streamstatus

import "github.com/esfands/retpaladinbot/internal/global"

type Module struct{}

func NewStreamStatusModule(gctx global.Context) *Module {
	return &Module{}
}
