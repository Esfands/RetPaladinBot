package db

import (
	"context"
	"database/sql"
	"log/slog"
)

type StreamStatus struct {
	ID        string
	StreamID  string
	GameID    sql.NullString
	GameName  sql.NullString
	Live      bool
	Title     sql.NullString
	StartedAt string
	EndedAt   sql.NullString
}

// InsertStream inserts a new stream into the database
func (q *Queries) InsertStream(ctx context.Context, stream StreamStatus) error {
	stmt, err := q.db.Prepare("INSERT INTO stream_status (stream_id, game_id, game_name, live, title, started_at, ended_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			slog.Error("Failed to close statement", "error", err)
		}
	}(stmt)

	_, err = stmt.Exec(
		stream.StreamID,
		stream.GameID,
		stream.GameName,
		stream.Live,
		stream.Title,
		stream.StartedAt,
		stream.EndedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetLiveStream returns the currently live stream
func (q *Queries) GetLiveStream(ctx context.Context) (StreamStatus, error) {
	var stream StreamStatus
	err := q.db.QueryRow("SELECT * FROM stream_status WHERE live = 1").Scan(
		&stream.ID,
		&stream.StreamID,
		&stream.GameID,
		&stream.GameName,
		&stream.Live,
		&stream.Title,
		&stream.StartedAt,
		&stream.EndedAt,
	)
	if err != nil {
		return stream, err
	}

	return stream, nil
}

// StreamWentOffline updates the last ID of the stream that went offline
func (q *Queries) StreamWentOffline(ctx context.Context, streamID string, timeWentOffline sql.NullString) error {
	stmt, err := q.db.Prepare("UPDATE stream_status SET live = 0, ended_at=? WHERE stream_id = ?")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			slog.Error("Failed to close statement", "error", err)
		}
	}(stmt)

	_, err = stmt.Exec(timeWentOffline, streamID)
	if err != nil {
		return err
	}

	return nil
}

// GetMostRecentStreamStatus returns the most recent stream based off the the `id` which keeps track of the most recent stream
func (q *Queries) GetMostRecentStreamStatus(ctx context.Context) (StreamStatus, error) {
	var stream StreamStatus
	err := q.db.QueryRow("SELECT * FROM stream_status ORDER BY id DESC LIMIT 1").Scan(
		&stream.ID,
		&stream.StreamID,
		&stream.GameID,
		&stream.GameName,
		&stream.Live,
		&stream.Title,
		&stream.StartedAt,
		&stream.EndedAt,
	)
	if err != nil {
		return stream, err
	}

	return stream, nil
}

// UpdateStreamInfo updates the stream information in the database
func (q *Queries) UpdateStreamInfo(ctx context.Context, stream StreamStatus) error {
	stmt, err := q.db.Prepare("UPDATE stream_status SET game_id = ?, game_name = ?, title = ? WHERE stream_id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(stream.GameID, stream.GameName, stream.Title, stream.StreamID)
	if err != nil {
		return err
	}

	return nil
}
