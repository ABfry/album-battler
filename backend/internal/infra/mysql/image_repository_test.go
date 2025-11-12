// image_repository_test.go では mysqlImageRepository の振る舞いをドライバスタブで検証する。
// 目的: findBy / findAllBy が createImage の単一点実装を確実に通ることを保証し、
// 追加の回帰を防ぐ。外部依存に接続せず database/sql/driver を自前実装している。
package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestImageRepository_FindByID は単一レコード取得が createImage で整形されることを確認。
func TestImageRepository_FindByID(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	userID := uuid.New()
	battleID := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)

	repo := newImageRepoForTest(t, queryPlan{
		mustGetFindQuery(t, "id"): {
			columns: imageColumns(),
			rows: [][]driver.Value{
				{id.String(), userID.String(), battleID.String(), "https://example.com/one.png", now, 87.3, int64(12)},
			},
		},
	})

	got, err := repo.FindByID(context.Background(), id)
	if err != nil {
		t.Fatalf("FindByID error: %v", err)
	}
	if got == nil {
		t.Fatalf("FindByID returned nil without error")
	}
	if got.ID != id || got.UserID != userID || got.BattleID != battleID {
		t.Fatalf("unexpected entity: %+v", got)
	}
	if got.AIScore != 87.3 || got.UserScore != 12 {
		t.Fatalf("score mismatch: ai=%f user=%d", got.AIScore, got.UserScore)
	}
	if got.ImageURL != "https://example.com/one.png" || !got.UploadedAt.Equal(now) {
		t.Fatalf("payload mismatch: %+v", got)
	}
}

// TestImageRepository_FindByUserID は複数件取得が createImage を共有することを確認。
func TestImageRepository_FindByUserID(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	firstBattle := uuid.New()
	secondBattle := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)

	repo := newImageRepoForTest(t, queryPlan{
		mustGetFindQuery(t, "user_id"): {
			columns: imageColumns(),
			rows: [][]driver.Value{
				{uuid.New().String(), userID.String(), firstBattle.String(), "https://example.com/a.png", now, 45.5, int64(7)},
				{uuid.New().String(), userID.String(), secondBattle.String(), "https://example.com/b.png", now.Add(time.Minute), nil, nil},
			},
		},
	})

	result, err := repo.FindImagesByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("FindByUserID error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 images, got %d", len(result))
	}
	if result[0].UserID != userID || result[1].UserID != userID {
		t.Fatalf("user id mismatch: %+v", result)
	}
	if result[1].AIScore != 0 || result[1].UserScore != 0 {
		t.Fatalf("null handling failed: %+v", result[1])
	}
	if result[0].BattleID != firstBattle || result[1].BattleID != secondBattle {
		t.Fatalf("battle id mismatch: %+v", result)
	}
}

// --- テストヘルパー ---

// queryPlan はクエリ文字列ごとの期待行を表す。
type queryPlan map[string]*stubRowsTemplate

// execPlan は INSERT/UPDATE 等の実行を表す。
type execPlan map[string]*execExpectation

type execExpectation struct {
	wantArgs []driver.Value
	result   driver.Result
	err      error
	called   bool
}

// newImageRepoForTest はクエリ結果を固定した mysqlImageRepository を返す。
func newImageRepoForTest(t *testing.T, plan queryPlan) *mysqlImageRepository {
	t.Helper()

	driverName := fmt.Sprintf("image_stub_%s", uuid.New().String())
	sql.Register(driverName, &stubDriver{queries: plan})

	db, err := sql.Open(driverName, "")
	if err != nil {
		t.Fatalf("stub db open failed: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	return &mysqlImageRepository{db: db}
}

// mustGetFindQuery は dao 層のクエリ生成結果をテストでも使い回し、
// 文字列のズレによる誤検知を防ぐ。
func mustGetFindQuery(t *testing.T, key string) string {
	t.Helper()
	query, err := getFindQuery(ImagesTable, key)
	if err != nil {
		t.Fatalf("getFindQuery failed: %v", err)
	}
	return query
}

// imageColumns は images テーブルのカラムリストをコピーで返す。
func imageColumns() []string {
	src := tableColumns[ImagesTable]
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

// stubDriver は database/sql/driver を実装し、事前定義した結果を返す。
type stubDriver struct {
	queries queryPlan
	execs   execPlan
}

// Open はスタブコネクションを返す。
func (d *stubDriver) Open(string) (driver.Conn, error) {
	return &stubConn{queries: d.queries, execs: d.execs}, nil
}

// stubConn は QueryContext のみをサポートする簡易接続。
type stubConn struct {
	queries queryPlan
	execs   execPlan
}

func (c *stubConn) Prepare(string) (driver.Stmt, error) { return &stubStmt{}, nil }
func (c *stubConn) Close() error                        { return nil }
func (c *stubConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions not supported in stub")
}

// QueryContext はクエリに一致するスタブ行を返す。
func (c *stubConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	tpl, ok := c.queries[query]
	if !ok {
		return nil, fmt.Errorf("unexpected query: %s", query)
	}
	return tpl.clone(), nil
}

// ExecContext は登録済みの execPlan に従って挙動を決める。
func (c *stubConn) ExecContext(_ context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if c.execs == nil {
		return nil, fmt.Errorf("unexpected exec query without plan: %s", query)
	}
	exp, ok := c.execs[query]
	if !ok {
		return nil, fmt.Errorf("unexpected exec query: %s", query)
	}
	if exp.wantArgs != nil {
		got := make([]driver.Value, len(args))
		for i, a := range args {
			got[i] = a.Value
		}
		if !reflect.DeepEqual(got, exp.wantArgs) {
			return nil, fmt.Errorf("exec args mismatch: want=%v got=%v", exp.wantArgs, got)
		}
	}
	exp.called = true
	if exp.err != nil {
		return nil, exp.err
	}
	if exp.result != nil {
		return exp.result, nil
	}
	return stubResult{}, nil
}

// stubStmt は Prepare 時に必要となるが使用しない。
type stubStmt struct{}

func (s *stubStmt) Close() error  { return nil }
func (s *stubStmt) NumInput() int { return -1 }
func (s *stubStmt) Exec(_ []driver.Value) (driver.Result, error) {
	return nil, errors.New("exec not supported in stub")
}
func (s *stubStmt) Query(_ []driver.Value) (driver.Rows, error) {
	return nil, errors.New("query via stmt not supported")
}

// stubResult は Exec のダミー戻り値。
type stubResult struct{}

func (stubResult) LastInsertId() (int64, error) { return 0, nil }
func (stubResult) RowsAffected() (int64, error) { return 0, nil }

// stubRowsTemplate は複数回クエリされる可能性があるため clone で都度コピーする。
type stubRowsTemplate struct {
	columns []string
	rows    [][]driver.Value
}

func (tpl *stubRowsTemplate) clone() *stubRows {
	data := make([][]driver.Value, len(tpl.rows))
	for i, row := range tpl.rows {
		rowCopy := make([]driver.Value, len(row))
		copy(rowCopy, row)
		data[i] = rowCopy
	}
	cols := make([]string, len(tpl.columns))
	copy(cols, tpl.columns)
	return &stubRows{
		columns: cols,
		rows:    data,
	}
}

// stubRows は driver.Rows インターフェースを実装する。
type stubRows struct {
	columns []string
	rows    [][]driver.Value
	index   int
}

func (r *stubRows) Columns() []string { return r.columns }
func (r *stubRows) Close() error      { return nil }

func (r *stubRows) Next(dest []driver.Value) error {
	if r.index >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.index])
	r.index++
	return nil
}
