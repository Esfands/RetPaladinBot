package streamstatus

import "github.com/esfands/retpaladinbot/internal/global"

type StreamStatusModule struct{}

func NewStreamStatusModule(gctx global.Context) *StreamStatusModule {
	return &StreamStatusModule{}
}
