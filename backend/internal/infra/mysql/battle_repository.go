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
	query := `
		SELECT b.id, b.room_id, b.started_at, b.theme,
		       b.current_phase, b.selecting_started_at, b.clap_phase_started_at,
		       b.clap_current_user_index, b.result_started_at,
		       GROUP_CONCAT(bu.user_id ORDER BY bu.user_id SEPARATOR ',') as user_ids
		FROM battles b
		LEFT JOIN battle_users bu ON b.id = bu.battle_id
		WHERE b.id = ?
		GROUP BY b.id, b.room_id, b.started_at, b.theme, b.current_phase,
		         b.selecting_started_at, b.clap_phase_started_at,
		         b.clap_current_user_index, b.result_started_at
	`

	row := r.db.QueryRowContext(ctx, query, id.String())
	return r.createBattle(row)
}

// FindByRoomID はルーム ID に紐づく最新のバトルを取得する。
// why: ルームごとの進行状況確認で RoomID からバトルを逆引きするため。
func (r *mysqlBattleRepository) FindByRoomID(ctx context.Context, roomID uuid.UUID) (*entity.Battle, error) {
	query := `
		SELECT b.id, b.room_id, b.started_at, b.theme,
		       b.current_phase, b.selecting_started_at, b.clap_phase_started_at,
		       b.clap_current_user_index, b.result_started_at,
		       GROUP_CONCAT(bu.user_id ORDER BY bu.user_id SEPARATOR ',') as user_ids
		FROM battles b
		LEFT JOIN battle_users bu ON b.id = bu.battle_id
		WHERE b.room_id = ?
		GROUP BY b.id, b.room_id, b.started_at, b.theme, b.current_phase,
		         b.selecting_started_at, b.clap_phase_started_at,
		         b.clap_current_user_index, b.result_started_at
		ORDER BY b.started_at DESC
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, roomID.String())
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
		battle.Theme,
		battle.CurrentPhase,
		battle.SelectingStartedAt,
		battle.ClapPhaseStartedAt,
		battle.ClapCurrentUserIndex,
		battle.ResultStartedAt,
	)
}

// --- private helpers ---

// createBattle は 1 行分のデータを entity.Battle に変換する。
// why: スキーマに変更が入っても変換ロジックを 1 箇所で管理するため。
func (r *mysqlBattleRepository) createBattle(scanner rowScanner) (*entity.Battle, error) {
	var (
		idStr                  string
		roomIDStr              string
		startedAt              time.Time
		theme                  string
		currentPhase           sql.NullString
		selectingStartedAt     sql.NullTime
		clapPhaseStartedAt     sql.NullTime
		clapCurrentUserIndex   sql.NullInt32
		resultStartedAt        sql.NullTime
		userIDsStr             sql.NullString // GROUP_CONCAT の結果は NULL の可能性がある
	)

	if err := scanner.Scan(&idStr, &roomIDStr, &startedAt, &theme,
		&currentPhase, &selectingStartedAt, &clapPhaseStartedAt,
		&clapCurrentUserIndex, &resultStartedAt, &userIDsStr); err != nil {
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

	// battle_users テーブルから取得した user_id 一覧をパース
	var userIDs []uuid.UUID
	if userIDsStr.Valid && userIDsStr.String != "" {
		userIDStrings := parseCommaSeparated(userIDsStr.String)
		userIDs = make([]uuid.UUID, 0, len(userIDStrings))
		for _, uidStr := range userIDStrings {
			uid, err := uuid.Parse(uidStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse user id %s: %w", uidStr, err)
			}
			userIDs = append(userIDs, uid)
		}
	}

	// NullTime/NullString/NullInt32 をポインタに変換
	var selectingStartedAtPtr *time.Time
	if selectingStartedAt.Valid {
		selectingStartedAtPtr = &selectingStartedAt.Time
	}

	var clapPhaseStartedAtPtr *time.Time
	if clapPhaseStartedAt.Valid {
		clapPhaseStartedAtPtr = &clapPhaseStartedAt.Time
	}

	var clapCurrentUserIndexPtr *int
	if clapCurrentUserIndex.Valid {
		val := int(clapCurrentUserIndex.Int32)
		clapCurrentUserIndexPtr = &val
	}

	var resultStartedAtPtr *time.Time
	if resultStartedAt.Valid {
		resultStartedAtPtr = &resultStartedAt.Time
	}

	currentPhaseStr := "selecting" // デフォルト値
	if currentPhase.Valid {
		currentPhaseStr = currentPhase.String
	}

	return &entity.Battle{
		ID:                   battleID,
		RoomID:               roomID,
		StartedAt:            startedAt,
		Theme:                theme,
		CurrentPhase:         currentPhaseStr,
		SelectingStartedAt:   selectingStartedAtPtr,
		ClapPhaseStartedAt:   clapPhaseStartedAtPtr,
		ClapCurrentUserIndex: clapCurrentUserIndexPtr,
		ResultStartedAt:      resultStartedAtPtr,
		UserIDs:              userIDs,
	}, nil
}

// parseCommaSeparated はカンマ区切り文字列を分割する。
// why: GROUP_CONCAT の結果を []string に変換するヘルパー。
func parseCommaSeparated(s string) []string {
	if s == "" {
		return []string{}
	}
	result := make([]string, 0)
	for _, part := range splitByComma(s) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// splitByComma はカンマで分割する。
func splitByComma(s string) []string {
	parts := []string{}
	current := ""
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(s[i])
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// trimSpace はスペースを除去する。
func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
