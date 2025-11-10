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
	var idStr string
	var name string
	var iconUrl string
	var hashedPassword string
	var createdAt time.Time

	query := getFindByIdQuery(UsersTable)
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(&idStr, &name, &iconUrl, &hashedPassword, &createdAt)
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

	return &entity.User{
		ID:             parsedID,
		Name:           name,
		IconUrl:        iconUrl,
		HashedPassword: hashedPassword,
		CreatedAt:      createdAt,
	}, nil
}

func (r *mysqlUserRepository) Save(ctx context.Context, user *entity.User) error {
	query := getSaveQuery(UsersTable)

	_, err := r.db.ExecContext(ctx, query, user.ID.String(), user.Name, user.IconUrl, user.HashedPassword, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
