// battle_repository_test.go では mysqlBattleRepository の I/O をテーブル駆動で確認する。
// 目的: find/save の各パスが createBattle を経由して正しくエンティティ化されることを保証する。
package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/google/uuid"
)

// TestBattleRepository_FindByID は主キー検索が期待どおり動作するか検証する。
func TestBattleRepository_FindByID(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC().Truncate(time.Second)
	battleID := uuid.New()
	roomID := uuid.New()

	repo := newBattleRepoForTest(
		t,
		queryPlan{
			mustGetBattleQuery(t, "id"): {
				columns: battleColumns(),
				rows: [][]driver.Value{
					{battleID.String(), roomID.String(), now},
				},
			},
		},
		nil,
	)

	got, err := repo.FindByID(context.Background(), battleID)
	if err != nil {
		t.Fatalf("FindByID error: %v", err)
	}
	if got == nil {
		t.Fatalf("FindByID returned nil")
	}
	if got.ID != battleID || got.RoomID != roomID {
		t.Fatalf("unexpected IDs: %+v", got)
	}
	if !got.StartedAt.Equal(now) {
		t.Fatalf("startedAt mismatch: %v", got.StartedAt)
	}
}

// TestBattleRepository_FindByRoomID は room_id 検索の挙動を確認する。
func TestBattleRepository_FindByRoomID(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC().Truncate(time.Second)
	roomID := uuid.New()
	battleID := uuid.New()

	repo := newBattleRepoForTest(
		t,
		queryPlan{
			mustGetBattleQuery(t, "room_id"): {
				columns: battleColumns(),
				rows: [][]driver.Value{
					{battleID.String(), roomID.String(), now},
				},
			},
		},
		nil,
	)

	got, err := repo.FindByRoomID(context.Background(), roomID)
	if err != nil {
		t.Fatalf("FindByRoomID error: %v", err)
	}
	if got == nil || got.RoomID != roomID {
		t.Fatalf("unexpected room id: %+v", got)
	}
}

// TestBattleRepository_Save は upsert クエリが期待どおり呼ばれることを検証する。
func TestBattleRepository_Save(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC().Truncate(time.Second)
	roomID := uuid.New()
	battle := &entity.Battle{
		ID:        uuid.New(),
		RoomID:    roomID,
		StartedAt: now,
	}

	query := upsertBattlesQuery()
	execs := execPlan{
		query: {
			wantArgs: []driver.Value{
				battle.ID.String(),
				battle.RoomID.String(),
				battle.StartedAt,
			},
		},
	}

	repo := newBattleRepoForTest(t, nil, execs)

	if err := repo.Save(context.Background(), battle); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	if !execs[query].called {
		t.Fatalf("expected exec plan for %s to be called", query)
	}
}

// --- テストヘルパー ---

func newBattleRepoForTest(t *testing.T, queries queryPlan, execs execPlan) *mysqlBattleRepository {
	t.Helper()

	driverName := fmt.Sprintf("battle_stub_%s", uuid.New().String())
	sql.Register(driverName, &stubDriver{queries: queries, execs: execs})

	db, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("stub db open failed: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	return &mysqlBattleRepository{db: db}
}

func mustGetBattleQuery(t *testing.T, key string) string {
	t.Helper()
	query, err := getFindQuery(BattlesTable, key)
	if err != nil {
		t.Fatalf("getFindQuery battles failed: %v", err)
	}
	return query
}

func battleColumns() []string {
	src := tableColumns[BattlesTable]
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

func upsertBattlesQuery() string {
	columns := tableColumns[BattlesTable]
	placeholders := make([]string, len(columns))
	assignments := make([]string, len(columns))
	for i, column := range columns {
		placeholders[i] = "?"
		assignments[i] = fmt.Sprintf("%s = VALUES(%s)", column, column)
	}
	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		BattlesTable,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(assignments, ", "),
	)
}
