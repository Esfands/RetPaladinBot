package db

import "context"

type CustomCommand struct {
	Name       string
	Response   string
	UsageCount int
}

func (q *Queries) InsertCustomCommand(ctx context.Context, command CustomCommand) error {
	stmt, err := q.db.Prepare("INSERT INTO custom_commands (name, response, usage_count) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(command.Name, command.Response, command.UsageCount)
	if err != nil {
		return err
	}

	return nil
}
