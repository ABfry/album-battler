package mysql

import (
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

func getFindByIdQuery(table TableName) string {
	columns, ok := tableColumns[table]
	if !ok {
		panic(fmt.Sprintf("unknown table: %s", table))
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = ?",
		strings.Join(columns, ", "),
		table,
	)

	return query
}

func getSaveQuery(table TableName) string {
	columns, ok := tableColumns[table]
	if !ok {
		panic(fmt.Sprintf("unknown table: %s", table))
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

	return query
}
