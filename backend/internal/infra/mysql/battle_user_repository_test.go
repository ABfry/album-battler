// battle_user_repository_test.go では mysqlBattleUserRepository の Save をテーブル駆動で確認する。
// 目的: 複合PK（battle_id, user_id）に対する upsert が正しく実行され、バリデーションが機能することを保証する。
package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// TestBattleUserRepository_SaveSuccess は正常系で upsert が実行されることを確認。
func TestBattleUserRepository_SaveSuccess(t *testing.T) {
	t.Parallel()

	battleID := uuid.New()
	userID := uuid.New()

	query := upsertBattleUsersQuery()
	execs := execPlan{
		query: {
			wantArgs: []driver.Value{
				battleID.String(),
				userID.String(),
				int64(0),
			},
		},
	}

	repo := newBattleUserRepoForTest(t, execs)

	if err := repo.Save(context.Background(), battleID, userID); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	if !execs[query].called {
		t.Fatalf("expected exec plan for %s to be called", query)
	}
}

// TestBattleUserRepository_SaveValidation は UUID が nil の場合にエラーを返すことを確認する。
func TestBattleUserRepository_SaveValidation(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("mysql", "")
	if err != nil {
		t.Fatalf("open default driver failed: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	repo := &mysqlBattleUserRepository{db: db}

	cases := []struct {
		name    string
		battle  uuid.UUID
		user    uuid.UUID
		wantErr string
	}{
		{name: "missing battle", battle: uuid.Nil, user: uuid.New(), wantErr: "battleID"},
		{name: "missing user", battle: uuid.New(), user: uuid.Nil, wantErr: "userID"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Save(context.Background(), tc.battle, tc.user)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
			}
		})
	}
}

// --- テストヘルパー ---

func newBattleUserRepoForTest(t *testing.T, execs execPlan) *mysqlBattleUserRepository {
	t.Helper()

	driverName := fmt.Sprintf("battle_user_stub_%s", uuid.New().String())
	sql.Register(driverName, &stubDriver{execs: execs})

	db, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("stub db open failed: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	return &mysqlBattleUserRepository{db: db}
}

func upsertBattleUsersQuery() string {
	columns := tableColumns[BattleUsersTable]
	placeholders := make([]string, len(columns))
	assignments := make([]string, len(columns))
	for i, column := range columns {
		placeholders[i] = "?"
		assignments[i] = fmt.Sprintf("%s = VALUES(%s)", column, column)
	}
	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		BattleUsersTable,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(assignments, ", "),
	)
}
