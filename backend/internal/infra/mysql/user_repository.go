package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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
	row, err := findRowByKey(ctx, r.db, UsersTable, "id", id.String())
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

func (r *mysqlUserRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.User, error) {
	if len(ids) == 0 {
		return []*entity.User{}, nil
	}

	// IDをstring配列に変換
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.String()
	}

	// IN句用のプレースホルダーを生成
	placeholders := make([]string, len(idStrings))
	args := make([]interface{}, len(idStrings))
	for i, idStr := range idStrings {
		placeholders[i] = "?"
		args[i] = idStr
	}

	// クエリ構築
	query := fmt.Sprintf(
		"SELECT id, name, icon_url, hashed_password, created_at FROM users WHERE id IN (%s)",
		strings.Join(placeholders, ", "),
	)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			fmt.Printf("failed to close rows: %v\n", closeErr)
		}
	}()

	var users []*entity.User
	for rows.Next() {
		var (
			idStr          string
			name           string
			iconUrl        string
			hashedPassword string
			createdAt      time.Time
		)

		if err := rows.Scan(&idStr, &name, &iconUrl, &hashedPassword, &createdAt); err != nil {
			return nil, err
		}

		parsedID, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user ID: %w", err)
		}

		users = append(users, &entity.User{
			ID:             parsedID,
			Name:           name,
			IconUrl:        iconUrl,
			HashedPassword: hashedPassword,
			CreatedAt:      createdAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
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
