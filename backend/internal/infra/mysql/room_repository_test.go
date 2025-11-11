// room_repository_test.go では mysqlRoomRepository の I/O 変換をテーブル駆動で検証する。
// 目的: createRoom / status 変換 / findAll の動きがスキーマ変更後も壊れないことを保障する。
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

// TestRoomRepository_FindByID は状態/期限の組み合わせをテーブル駆動で検証する。
func TestRoomRepository_FindByID(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC().Truncate(time.Second)
	hostID := uuid.New()
	id := uuid.New()

	cases := []struct {
		name        string
		status      string
		expiredTime driver.Value
		wantStatus  entity.RoomStatus
	}{
		{name: "waiting_with_expired", status: "waiting", expiredTime: now.Add(10 * time.Minute), wantStatus: entity.WaitJoin},
		{name: "full_without_expired", status: "full", expiredTime: nil, wantStatus: entity.FullyJoined},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := newRoomRepoForTest(
				t,
				queryPlan{
					mustGetRoomQuery(t, "id"): {
						columns: roomColumns(),
						rows: [][]driver.Value{
							{id.String(), 1234, hostID.String(), now, tc.expiredTime, tc.status},
						},
					},
				},
				nil,
			)

			got, err := repo.FindByID(context.Background(), id)
			if err != nil {
				t.Fatalf("FindByID error: %v", err)
			}
			if got == nil {
				t.Fatalf("FindByID returned nil room")
			}
			if got.Status != tc.wantStatus {
				t.Fatalf("unexpected status: want=%v got=%v", tc.wantStatus, got.Status)
			}
			if got.RoomNumber != 1234 {
				t.Fatalf("room number mismatch: %+v", got)
			}
			if got.HostUserID == nil || *got.HostUserID != hostID {
				t.Fatalf("host id mismatch: %+v", got.HostUserID)
			}
			if tc.expiredTime == nil && !got.ExpiredAt.IsZero() {
				t.Fatalf("expected zero expired_at but got %v", got.ExpiredAt)
			}
		})
	}
}

// TestRoomRepository_FindAll は複数行を createRoom が処理できるかを確認する。
func TestRoomRepository_FindAll(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC().Truncate(time.Second)
	hostA := uuid.New()
	hostB := uuid.New()

	repo := newRoomRepoForTest(
		t,
		queryPlan{
			selectAllRoomsQuery(): {
				columns: roomColumns(),
				rows: [][]driver.Value{
					{uuid.New().String(), 1111, hostA.String(), now, now.Add(time.Hour), "battling"},
					{uuid.New().String(), 2222, hostB.String(), now, nil, "result"},
				},
			},
		},
		nil,
	)

	got, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("FindAll error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 rooms, got %d", len(got))
	}
	if got[0].Status != entity.InBattle || got[1].Status != entity.Result {
		t.Fatalf("status mismatch: %+v", got)
	}
}

// TestRoomRepository_FindByRoomNumber は参加コード検索が正しく動くかを確認する。
func TestRoomRepository_FindByRoomNumber(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC().Truncate(time.Second)
	hostID := uuid.New()
	roomNumber := 5678

	repo := newRoomRepoForTest(
		t,
		queryPlan{
			mustGetRoomQuery(t, "room_number"): {
				columns: roomColumns(),
				rows: [][]driver.Value{
					{uuid.New().String(), roomNumber, hostID.String(), now, nil, "waiting"},
				},
			},
		},
		nil,
	)

	got, err := repo.FindByRoomNumber(context.Background(), roomNumber)
	if err != nil {
		t.Fatalf("FindByRoomNumber error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected room, got nil")
	}
	if got.RoomNumber != roomNumber {
		t.Fatalf("room number mismatch: %d", got.RoomNumber)
	}
	if got.Status != entity.WaitJoin {
		t.Fatalf("status mismatch: %v", got.Status)
	}
}

// TestRoomRepository_FindByIDMultipleCalls は同じクエリを複数回実行しても結果が得られることを確認する。
func TestRoomRepository_FindByIDMultipleCalls(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC().Truncate(time.Second)
	hostID := uuid.New()
	roomID := uuid.New()

	plan := queryPlan{
		mustGetRoomQuery(t, "id"): {
			columns: roomColumns(),
			rows: [][]driver.Value{
				{roomID.String(), 3141, hostID.String(), now, nil, "waiting"},
			},
		},
	}
	repo := newRoomRepoForTest(t, plan, nil)

	for i := 0; i < 2; i++ {
		got, err := repo.FindByID(context.Background(), roomID)
		if err != nil {
			t.Fatalf("iteration %d: FindByID error: %v", i, err)
		}
		if got == nil || got.ID != roomID {
			t.Fatalf("iteration %d: unexpected room %#v", i, got)
		}
	}
}

// TestRoomRepository_SaveSuccess は Save が正規化した値で upsert を実行することを検証する。
func TestRoomRepository_SaveSuccess(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC().Truncate(time.Second)
	expire := now.Add(30 * time.Minute)
	host := uuid.New()
	room := &entity.Room{
		ID:         uuid.New(),
		RoomNumber: 9999,
		HostUserID: &host,
		CreatedAt:  now,
		ExpiredAt:  expire,
		Status:     entity.InBattle,
	}

	statusStr, err := toDBRoomStatus(room.Status)
	if err != nil {
		t.Fatalf("unexpected status conversion error: %v", err)
	}

	query := upsertRoomsQuery()
	execs := execPlan{
		query: {
			wantArgs: []driver.Value{
				room.ID.String(),
				int64(room.RoomNumber),
				host.String(),
				room.CreatedAt,
				room.ExpiredAt,
				statusStr,
			},
		},
	}

	repo := newRoomRepoForTest(t, nil, execs)

	if err := repo.Save(context.Background(), room); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	if !execs[query].called {
		t.Fatalf("expected exec plan for %s to be called", query)
	}
}

// TestRoomRepository_SaveValidation は前提未充足時にエラーを返すことを確認する。
func TestRoomRepository_SaveValidation(t *testing.T) {
	t.Parallel()
	repo := &mysqlRoomRepository{db: nil}

	h := uuid.New()
	room := &entity.Room{
		ID:         uuid.New(),
		RoomNumber: 2,
		Status:     entity.RoomStatus(999),
		HostUserID: &h,
		CreatedAt:  time.Now(),
	}

	err := repo.Save(context.Background(), room)
	if err == nil || !strings.Contains(err.Error(), "unknown room status") {
		t.Fatalf("expected unknown room status error, got %v", err)
	}
}

// TestRoomRepository_SaveAllowsNilHost は HostUserID が nil の場合でも保存できることを確認する。
func TestRoomRepository_SaveAllowsNilHost(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC().Truncate(time.Second)
	room := &entity.Room{
		ID:         uuid.New(),
		RoomNumber: 5555,
		CreatedAt:  now,
		Status:     entity.WaitJoin,
	}

	statusStr, err := toDBRoomStatus(room.Status)
	if err != nil {
		t.Fatalf("status conversion error: %v", err)
	}

	query := upsertRoomsQuery()
	execs := execPlan{
		query: {
			wantArgs: []driver.Value{
				room.ID.String(),
				int64(room.RoomNumber),
				nil,
				room.CreatedAt,
				nil,
				statusStr,
			},
		},
	}

	repo := newRoomRepoForTest(t, nil, execs)
	if err := repo.Save(context.Background(), room); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	if !execs[query].called {
		t.Fatalf("expected exec plan for %s to be called", query)
	}
}

// --- テストヘルパー ---

func newRoomRepoForTest(t *testing.T, queries queryPlan, execs execPlan) *mysqlRoomRepository {
	t.Helper()

	driverName := fmt.Sprintf("room_stub_%s", uuid.New().String())
	sql.Register(driverName, &stubDriver{queries: queries, execs: execs})

	db, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("stub db open failed: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	return &mysqlRoomRepository{db: db}
}

func mustGetRoomQuery(t *testing.T, key string) string {
	t.Helper()
	query, err := getFindQuery(RoomsTable, key)
	if err != nil {
		t.Fatalf("getFindQuery rooms failed: %v", err)
	}
	return query
}

func selectAllRoomsQuery() string {
	columns := tableColumns[RoomsTable]
	return fmt.Sprintf("SELECT %s FROM %s", strings.Join(columns, ", "), RoomsTable)
}

func upsertRoomsQuery() string {
	columns := tableColumns[RoomsTable]
	placeholders := make([]string, len(columns))
	assignments := make([]string, len(columns))
	for i, column := range columns {
		placeholders[i] = "?"
		assignments[i] = fmt.Sprintf("%s = VALUES(%s)", column, column)
	}
	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		RoomsTable,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(assignments, ", "),
	)
}

func roomColumns() []string {
	src := tableColumns[RoomsTable]
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}
