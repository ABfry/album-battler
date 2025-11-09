package mysql

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

const (
	dsn = "album_user:album_pass@tcp(127.0.0.1:3306)/album_battler?parseTime=true&loc=Local"
)

// allowedTestTables は、テストヘルパー関数で操作を許可するテーブル名のリストです。
var allowedTestTables = map[string]struct{}{
	"users":  {},
	"rooms":  {},
	"battles":{},
	"images": {},
}

// validateTableName は、指定されたテーブル名が許可リストに含まれているか検証します。
func validateTableName(t *testing.T, table string) {
	t.Helper()
	if _, ok := allowedTestTables[table]; !ok {
		t.Fatalf("Invalid table name used in test helper: %s", table)
	}
}

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("データベースへの接続に失敗しました: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("データベースへのPingに失敗しました: %v", err)
	}
	return db
}

func showTable(t *testing.T, db *sql.DB, table string) {
	validateTableName(t, table) // テーブル名を検証

	// fmt.Sprintfを使用して安全にクエリを構築します（検証後なので安全）
	query := fmt.Sprintf("SELECT * FROM %s", table)
	rows, err := db.Query(query)
	if err != nil {
		t.Fatalf("テーブルの表示に失敗しました: %v", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		t.Fatalf("カラムの取得に失敗しました: %v", err)
	}

	values := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			t.Fatalf("行のスキャンに失敗しました: %v", err)
		}

		var valueStr string
		for i, col := range values {
			if i > 0 {
				valueStr += ", "
			}
			valueStr += string(col)
		}
		t.Logf("Row: %s", valueStr)
	}
	if err = rows.Err(); err != nil {
		t.Fatalf("行の反復中にエラーが発生しました: %v", err)
	}
}

// assertEqual は、2つのオブジェクトが等しいかチェックし、等しくない場合はエラーを報告します。
func assertEqual(t *testing.T, want, got interface{}) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("mismatch: want=%+v, got=%+v", want, got)
	}
}

func cleanupTestUser(db *sql.DB, table string, id uuid.UUID, t *testing.T) {
	validateTableName(t, table) // テーブル名を検証

	// fmt.Sprintfを使用して安全にクエリを構築します（検証後なので安全）
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", table)
	_, err := db.Exec(query, id.String())
	if err != nil {
		t.Fatalf("テストデータのクリーンアップに失敗しました: %v", err)
	}
}