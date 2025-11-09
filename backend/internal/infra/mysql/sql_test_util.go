package mysql

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

const (
	dsn = "album_user:album_pass@tcp(127.0.0.1:3306)/album_battler?parseTime=true&loc=Local"
)

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
	rows, err := db.Query("SELECT * FROM " + table)
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

func cleanupTest(db *sql.DB, table string, id uuid.UUID, t *testing.T) {
	_, err := db.Exec("DELETE FROM "+table+" WHERE id = ?", id.String())
	if err != nil {
		t.Fatalf("テストデータのクリーンアップに失敗しました: %v", err)
	}
}
