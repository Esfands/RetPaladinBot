package db

import (
	"context"
	"database/sql"
	"log/slog"
)

// Chatter represents the chatter model
type Chatter struct {
	TID         int
	Username    string
	DisplayName string
}

// InsertChatter inserts a new chatter into the database
func (q *Queries) InsertChatter(ctx context.Context, chatter Chatter) error {
	stmt, err := q.db.Prepare("INSERT INTO chatters (tid, username, display_name) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			slog.Error("Failed to close statement", "error", err)
		}
	}(stmt)

	_, err = stmt.Exec(chatter.TID, chatter.Username, chatter.DisplayName)
	if err != nil {
		return err
	}

	return nil
}
