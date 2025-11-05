package repository

import (
	"context"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/google/uuid"
)

// AlbumRepository はアルバムの永続化を担当するリポジトリインターフェース
type AlbumRepository interface {
	// Create は新しいアルバムを作成する
	Create(ctx context.Context, album *entity.Album) error

	// FindByID は指定されたIDのアルバムを取得する
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Album, error)

	// Update はアルバム情報を更新する
	Update(ctx context.Context, album *entity.Album) error

	// Delete はアルバムを削除する
	Delete(ctx context.Context, id uuid.UUID) error

	// List は全てのアルバムを取得する
	List(ctx context.Context) ([]*entity.Album, error)
}
