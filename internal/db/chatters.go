package db

import (
	"context"
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
	defer stmt.Close()

	_, err = stmt.Exec(chatter.TID, chatter.Username, chatter.DisplayName)
	if err != nil {
		return err
	}

	return nil
}
