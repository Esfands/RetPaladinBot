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

type DefaultCommand struct {
	Name               string
	Aliases            string
	Permissions        string
	Description        string
	DynamicDescription string
	GlobalCooldown     int
	UserCooldown       int
	OfflineOnly        int
	OnlineOnly         int
	UsageCount         int
}

func (q *Queries) InsertDefaultCommand(ctx context.Context, command DefaultCommand) error {
	stmt, err := q.db.Prepare(
		"INSERT INTO commands (name, aliases, permissions, description, dynamic_description, global_cooldown, user_cooldown, offline_only, online_only, usage_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		command.Name,
		command.Aliases,
		command.Permissions,
		command.Description,
		command.DynamicDescription,
		command.GlobalCooldown,
		command.UserCooldown,
		command.OfflineOnly,
		command.OnlineOnly,
		command.UsageCount,
	)
	if err != nil {
		return err
	}

	return nil
}

func (q *Queries) GetAllDefaultCommands(ctx context.Context) ([]DefaultCommand, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT name, aliases, permissions, description, dynamic_description, global_cooldown, user_cooldown, offline_only, online_only, usage_count FROM commands")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []DefaultCommand
	for rows.Next() {
		var command DefaultCommand
		if err := rows.Scan(&command.Name, &command.Aliases, &command.Permissions, &command.Description, &command.DynamicDescription, &command.GlobalCooldown, &command.UserCooldown, &command.OfflineOnly, &command.OnlineOnly, &command.UsageCount); err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return commands, nil
}
