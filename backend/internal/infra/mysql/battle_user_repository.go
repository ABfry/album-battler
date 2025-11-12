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

// Save は battle_users テーブルへ upsert を行う。
// why: 参加登録とスコア更新を同じエンドポイントで扱い、複合PK重複時は score を最新化するため。
func (r *mysqlBattleUserRepository) Save(ctx context.Context, battleID, userID uuid.UUID) error {
	if battleID == uuid.Nil {
		return errors.New("battleID is required")
	}
	if userID == uuid.Nil {
		return errors.New("userID is required")
	}

	return save(
		ctx,
		r.db,
		BattleUsersTable,
		battleID.String(),
		userID.String(),
		0, // スコア初期値
	)
}
