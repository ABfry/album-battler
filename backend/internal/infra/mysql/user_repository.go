package mysql

import (
	"database/sql"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
)

var _ repository.UserRepository = (*mysqlUserRepository)(nil)

type mysqlUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) FindByID(id string) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRow("SELECT id, name, icon_url, hashed_password FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.IconUrl, &user.HashedPassword)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
