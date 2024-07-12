package services

import (
	"github.com/esfands/retpaladinbot/internal/services/helix"
	"github.com/esfands/retpaladinbot/internal/services/scheduler"
)

type Crate struct {
	Helix     helix.Service
	Scheduler scheduler.Service
}
