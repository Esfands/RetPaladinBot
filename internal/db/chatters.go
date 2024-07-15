package db

import (
	"context"
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
		slog.Error("Failed to prepare statement", "error", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(chatter.TID, chatter.Username, chatter.DisplayName)
	if err != nil {
		slog.Error("Failed to execute statement", "error", err)
		return err
	}

	return nil
}
