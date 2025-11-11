// battle_repository.go は battles テーブルを扱う MySQL 実装を提供する。
// 責務: バトルの保存と検索（ID・RoomID）を担い、行データと entity.Battle を相互変換する。
// 依存: database/sql、dao.go の共通クエリヘルパー。参加ユーザー一覧は battle_users テーブル側で管理する。
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/google/uuid"
)

var _ repository.BattleRepository = (*mysqlBattleRepository)(nil)

// mysqlBattleRepository は BattleRepository を MySQL で実装した構造体。
type mysqlBattleRepository struct {
	db *sql.DB
}

// NewBattleRepository は MySQL 接続を受け取り BattleRepository 実装を返す。
// why: 依存注入で具象実装を隠し、テスト容易性を高める。
func NewBattleRepository(db *sql.DB) repository.BattleRepository {
	return &mysqlBattleRepository{db: db}
}

// FindByID は主キー検索でバトルを 1 件だけ取得する。
// why: バトル詳細表示や集計で ID を直接指定するケースがあるため。
func (r *mysqlBattleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Battle, error) {
	row, err := findRowByKey(ctx, r.db, BattlesTable, "id", id.String())
	if err != nil {
		return nil, err
	}
	return r.createBattle(row)
}

// FindByRoomID はルーム ID に紐づく最新のバトルを取得する。
// why: ルームごとの進行状況確認で RoomID からバトルを逆引きするため。
func (r *mysqlBattleRepository) FindByRoomID(ctx context.Context, roomID uuid.UUID) (*entity.Battle, error) {
	row, err := findRowByKey(ctx, r.db, BattlesTable, "room_id", roomID.String())
	if err != nil {
		return nil, err
	}
	return r.createBattle(row)
}

// Save は battles テーブルへ upsert を行う。
// why: 1 つのメソッドで新規作成と更新の両方を賄い、カラム順序の重複定義を避ける。
func (r *mysqlBattleRepository) Save(ctx context.Context, battle *entity.Battle) error {
	return save(
		ctx,
		r.db,
		BattlesTable,
		battle.ID.String(),
		battle.RoomID.String(),
		battle.StartedAt,
	)
}

// --- private helpers ---

// createBattle は 1 行分のデータを entity.Battle に変換する。
// why: スキーマに変更が入っても変換ロジックを 1 箇所で管理するため。
func (r *mysqlBattleRepository) createBattle(scanner rowScanner) (*entity.Battle, error) {
	var (
		idStr     string
		roomIDStr string
		startedAt time.Time
	)

	if err := scanner.Scan(&idStr, &roomIDStr, &startedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	battleID, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse battle id: %w", err)
	}
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse room id: %w", err)
	}

	return &entity.Battle{
		ID:        battleID,
		RoomID:    roomID,
		StartedAt: startedAt,
		UserIDs:   []uuid.UUID{}, // battle_users 側で管理するため、ここでは空配列を返す
	}, nil
}
