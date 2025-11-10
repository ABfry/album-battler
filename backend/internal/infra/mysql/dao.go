package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type TableName string

const (
	UsersTable       TableName = "users"
	RoomsTable       TableName = "rooms"
	BattlesTable     TableName = "battles"
	BattleUsersTable TableName = "battle_users"
	ImagesTable      TableName = "images"
)

var tableColumns = map[TableName][]string{
	UsersTable:       {"id", "name", "icon_url", "hashed_password", "created_at"},
	RoomsTable:       {"id", "room_number", "host_user_id", "created_at"},
	BattlesTable:     {"id", "room_id", "status", "started_at", "ended_at"},
	BattleUsersTable: {"battle_id", "user_id", "score"},
	ImagesTable:      {"id", "user_id", "battle_id", "image_url", "uploaded_at", "ai_score", "user_score"},
}

func findByKey(ctx context.Context, db *sql.DB, table TableName, key string, value interface{}) (*sql.Row, error) {
	columns, ok := tableColumns[table]
	if !ok {
		return nil, fmt.Errorf("unknown table: %s", table)
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s = ?",
		strings.Join(columns, ", "),
		table,
		key,
	)

	row := db.QueryRowContext(ctx, query, value)
	return row, nil
}

func save(ctx context.Context, db *sql.DB, table TableName, values ...interface{}) error {
	columns, ok := tableColumns[table]
	if !ok {
		return fmt.Errorf("unknown table: %s", table)
	}

	if len(values) != len(columns) {
		return fmt.Errorf("value count mismatch for table %s: want %d, got %d",
			table, len(columns), len(values))
	}

	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = "?"
	}

	updateAssignments := make([]string, len(columns))
	for i, column := range columns {
		updateAssignments[i] = fmt.Sprintf("%s = VALUES(%s)", column, column)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(updateAssignments, ", "),
	)

	_, err := db.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("failed to save %s: %w", table, err)
	}
	return nil
}
