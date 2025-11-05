package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Album はアルバムのドメインエンティティ
type Album struct {
	ID          uuid.UUID
	Name        string
	Description string
	ImageURL    string // S3などに保存された画像のURL
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewAlbum は新しいAlbumエンティティを作成する
func NewAlbum(name, description string) *Album {
	now := time.Now()
	return &Album{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// UpdateImage は画像URLを更新する
func (a *Album) UpdateImage(imageURL string) error {
	if imageURL == "" {
		return errors.New("imageURL is required")
	}
	a.ImageURL = imageURL
	a.UpdatedAt = time.Now()
	return nil
}

// UpdateInfo はアルバム情報を更新する
func (a *Album) UpdateInfo(name, description string) error {
	if name == "" {
		return errors.New("name is required")
	}
	a.Name = name
	a.Description = description
	a.UpdatedAt = time.Now()
	return nil
}
