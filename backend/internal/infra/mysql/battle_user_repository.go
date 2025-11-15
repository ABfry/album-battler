// battle_user_repository.go は battle_users テーブルを扱う MySQL 実装を提供する。
// 責務: バトルとユーザーの紐づけを保存し、スコア初期値を管理する。
// 依存: database/sql 接続と dao.go の save ヘルパー。
package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/google/uuid"
)

var _ repository.BattleUserRepository = (*mysqlBattleUserRepository)(nil)

type mysqlBattleUserRepository struct {
	db *sql.DB
}

// NewBattleUserRepository は MySQL 接続を受け取り BattleUserRepository 実装を返す。
// why: 他レイヤーから具象依存を隠し、テスト時にはモックに差し替えられるようにする。
func NewBattleUserRepository(db *sql.DB) repository.BattleUserRepository {
	return &mysqlBattleUserRepository{db: db}
}

// SaveBatch は複数のユーザーをバトルに一括登録する。
func (r *mysqlBattleUserRepository) SaveBatch(ctx context.Context, battleID uuid.UUID, userIDs []uuid.UUID) error {
	if battleID == uuid.Nil {
		return errors.New("battleID is required")
	}
	if len(userIDs) == 0 {
		return errors.New("userIDs is required")
	}

	// バルクインサート用のクエリを構築
	query := "INSERT INTO battle_users (battle_id, user_id) VALUES "
	values := []interface{}{}
	placeholders := []string{}

	for _, userID := range userIDs {
		if userID == uuid.Nil {
			return errors.New("userID cannot be nil")
		}
		placeholders = append(placeholders, "(?, ?)")
		values = append(values, battleID.String(), userID.String())
	}

	query += placeholders[0]
	for i := 1; i < len(placeholders); i++ {
		query += ", " + placeholders[i]
	}

	_, err := r.db.ExecContext(ctx, query, values...)
	if err != nil {
		return errors.New("failed to save battle_users")
	}

	return nil
}
