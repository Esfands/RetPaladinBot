package db

import "context"

type CustomCommand struct {
	Name       string
	Response   string
	UsageCount int
}

// InsertCustomCommand inserts a new custom command into the database
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

// UpdateCustomCommand updates an existing custom command in the database
func (q *Queries) UpdateCustomCommand(ctx context.Context, command CustomCommand) error {
	stmt, err := q.db.Prepare("UPDATE custom_commands SET response = ? WHERE name = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(command.Response, command.Name)
	if err != nil {
		return err
	}

	return nil
}

func (q *Queries) IncrementCustomCommandUsageCount(ctx context.Context, name string) error {
	stmt, err := q.db.Prepare("UPDATE custom_commands SET usage_count = usage_count + 1 WHERE name = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCustomCommand deletes a custom command from the database
func (q *Queries) DeleteCustomCommand(ctx context.Context, name string) error {
	stmt, err := q.db.Prepare("DELETE FROM custom_commands WHERE name = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name)
	if err != nil {
		return err
	}

	return nil
}

// GetAllCustomCommands retrieves all custom commands from the database
func (q *Queries) GetAllCustomCommands(ctx context.Context) ([]CustomCommand, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT name, response, usage_count FROM custom_commands")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []CustomCommand
	for rows.Next() {
		var command CustomCommand
		if err := rows.Scan(&command.Name, &command.Response, &command.UsageCount); err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}
	return commands, rows.Err()
}

// GetCustomCommand retrieves a specific custom command from the database
func (q *Queries) GetCustomCommand(ctx context.Context, name string) (CustomCommand, error) {
	var command CustomCommand
	err := q.db.QueryRowContext(ctx, "SELECT name, response, usage_count FROM custom_commands WHERE name = ?", name).Scan(&command.Name, &command.Response, &command.UsageCount)
	if err != nil {
		return CustomCommand{}, err
	}
	return command, nil
}

type DefaultCommand struct {
	Name               string
	Aliases            string
	Permissions        string
	Description        string
	DynamicDescription string
	GlobalCooldown     int
	UserCooldown       int
	EnabledOffline     int
	EnabledOnline      int
	UsageCount         int
}

// InsertDefaultCommand inserts a new default command into the database
func (q *Queries) InsertDefaultCommand(ctx context.Context, command DefaultCommand) error {
	stmt, err := q.db.Prepare(
		"INSERT INTO commands (name, aliases, permissions, description, dynamic_description, global_cooldown, user_cooldown, enabled_offline, enabled_online, usage_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
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
		command.EnabledOffline,
		command.EnabledOnline,
		command.UsageCount,
	)
	return err
}

// UpdateDefaultCommand updates an existing default command in the database
func (q *Queries) UpdateDefaultCommand(ctx context.Context, command DefaultCommand) error {
	stmt, err := q.db.Prepare(
		"UPDATE commands SET aliases = ?, permissions = ?, description = ?, dynamic_description = ?, global_cooldown = ?, user_cooldown = ?, enabled_offline = ?, enabled_online = ? WHERE name = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		command.Aliases,
		command.Permissions,
		command.Description,
		command.DynamicDescription,
		command.GlobalCooldown,
		command.UserCooldown,
		command.EnabledOffline,
		command.EnabledOnline,
		command.Name,
	)
	return err
}

// IncrementDefaultCommandUsageCount increments the usage count of a default command in the database
func (q *Queries) IncrementDefaultCommandUsageCount(ctx context.Context, name string) error {
	stmt, err := q.db.Prepare("UPDATE commands SET usage_count = usage_count + 1 WHERE name = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name)
	return err
}

// DeleteDefaultCommand deletes a default command from the database
func (q *Queries) DeleteDefaultCommand(ctx context.Context, name string) error {
	stmt, err := q.db.Prepare("DELETE FROM commands WHERE name = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name)
	return err
}

// GetAllDefaultCommands retrieves all default commands from the database
func (q *Queries) GetAllDefaultCommands(ctx context.Context) ([]DefaultCommand, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT name, aliases, permissions, description, dynamic_description, global_cooldown, user_cooldown, enabled_offline, enabled_online, usage_count FROM commands")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []DefaultCommand
	for rows.Next() {
		var command DefaultCommand
		if err := rows.Scan(&command.Name, &command.Aliases, &command.Permissions, &command.Description, &command.DynamicDescription, &command.GlobalCooldown, &command.UserCooldown, &command.EnabledOffline, &command.EnabledOnline, &command.UsageCount); err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}
	return commands, rows.Err()
}
