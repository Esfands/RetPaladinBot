package turso

import (
	"database/sql"

	"github.com/esfands/retpaladinbot/internal/db"
)

type Service interface {
	DB() *sql.DB
	Queries() *db.Queries
}

type tursoService struct {
	db      *sql.DB
	queries *db.Queries
}

func (t *tursoService) DB() *sql.DB {
	return t.db
}

func (t *tursoService) Queries() *db.Queries {
	return t.queries
}
