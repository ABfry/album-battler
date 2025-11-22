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
	RoomUsersTable   TableName = "room_users"
	BattlesTable     TableName = "battles"
	BattleUsersTable TableName = "battle_users"
	ImagesTable      TableName = "images"
)

var tableColumns = map[TableName][]string{
	UsersTable:       {"id", "name", "icon_url", "hashed_password", "created_at"},
	RoomsTable:       {"id", "room_number", "host_user_id", "created_at", "expired_at", "status", "max_users", "battle_time_limit_seconds"},
	RoomUsersTable:   {"room_id", "user_id", "joined_at"},
	BattlesTable:     {"id", "room_id", "started_at", "theme"},
	BattleUsersTable: {"battle_id", "user_id"},
	ImagesTable:      {"id", "user_id", "battle_id", "image_url", "uploaded_at", "ai_score", "user_score", "ai_explanation"},
}

func findRowByKey(ctx context.Context, db *sql.DB, table TableName, key string, value interface{}) (*sql.Row, error) {
	query, err := getFindQuery(table, key)
	if err != nil {
		return nil, err
	}

	row := db.QueryRowContext(ctx, query, value)
	return row, nil
}

func findRowsByKey(ctx context.Context, db *sql.DB, table TableName, key string, value interface{}) (*sql.Rows, error) {
	query, err := getFindQuery(table, key)
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, query, value)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// findAllRows は WHERE 句なしで指定テーブル全件を取得する。
// why: RoomRepository などで単純な全件取得を繰り返し記述するのを避けるため。
func findAllRows(ctx context.Context, db *sql.DB, table TableName) (*sql.Rows, error) {
	columns, ok := tableColumns[table]
	if !ok {
		return nil, fmt.Errorf("unknown table: %s", table)
	}

	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(columns, ", "), table)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func getFindQuery(table TableName, key string) (string, error) {
	columns, ok := tableColumns[table]
	if !ok {
		return "", fmt.Errorf("unknown table: %s", table)
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s = ?",
		strings.Join(columns, ", "),
		table,
		key,
	)
	return query, nil
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

	// insertのためのプレースホルダーを構築
	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = "?"
	}

	// update クエリ構築
	updateAssignments := make([]string, len(columns))
	for i, column := range columns {
		updateAssignments[i] = fmt.Sprintf("%s = VALUES(%s)", column, column)
	}

	// upsert クエリを構築
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(updateAssignments, ", "),
	)

	// 実行
	_, err := db.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("failed to save %s: %w", table, err)
	}
	return nil
}
