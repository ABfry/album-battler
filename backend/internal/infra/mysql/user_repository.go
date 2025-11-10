package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/google/uuid"
)

var _ repository.UserRepository = (*mysqlUserRepository)(nil)

type mysqlUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	// 共通関数を利用して1件取得
	row, err := findByKey(ctx, r.db, UsersTable, "id", id.String())
	if err != nil {
		return nil, err
	}

	var (
		idStr          string
		name           string
		iconUrl        string
		hashedPassword string
		createdAt      time.Time
	)

	// スキャン処理
	if err := row.Scan(&idStr, &name, &iconUrl, &hashedPassword, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 見つからない場合はnilを返す
		}
		return nil, err
	}

	// UUIDへ変換
	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	// Entityを返す
	return &entity.User{
		ID:             parsedID,
		Name:           name,
		IconUrl:        iconUrl,
		HashedPassword: hashedPassword,
		CreatedAt:      createdAt,
	}, nil
}

func (r *mysqlUserRepository) Save(ctx context.Context, user *entity.User) error {
	return save(
		ctx,
		r.db,
		UsersTable,
		user.ID.String(),
		user.Name,
		user.IconUrl,
		user.HashedPassword,
		user.CreatedAt,
	)
}
