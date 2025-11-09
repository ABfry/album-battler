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
	var user entity.User
	var idStr string
	var createdAt time.Time

	query := `SELECT id, name, icon_url, hashed_password, created_at FROM users WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(&idStr, &user.Name, &user.IconUrl, &user.HashedPassword, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // ユーザーが見つからない場合はnilを返す
		}
		return nil, err
	}

	// 文字列IDをUUIDに変換
	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	user.ID = parsedID
	user.CreatedAt = createdAt

	return &user, nil
}

func (r *mysqlUserRepository) Save(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (id, name, icon_url, hashed_password, created_at) VALUES (?, ?, ?, ?, ?)
	          ON DUPLICATE KEY UPDATE name = VALUES(name), icon_url = VALUES(icon_url), hashed_password = VALUES(hashed_password)`

	_, err := r.db.ExecContext(ctx, query, user.ID.String(), user.Name, user.IconUrl, user.HashedPassword, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
