package services

import (
	"github.com/esfands/retpaladinbot/internal/services/auth"
	"github.com/esfands/retpaladinbot/internal/services/helix"
	"github.com/esfands/retpaladinbot/internal/services/scheduler"
	"github.com/esfands/retpaladinbot/internal/services/turso"
)

type Crate struct {
	Turso     turso.Service
	Helix     helix.Service
	Scheduler scheduler.Service
	Auth      auth.Authmen
}
