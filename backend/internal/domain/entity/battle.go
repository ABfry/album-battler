package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Battle struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	StartedAt time.Time
	Theme     string
	UserIDs   []uuid.UUID
}

func NewBattle(roomID uuid.UUID, userIDs []uuid.UUID, theme string) (*Battle, error) {
	if roomID == uuid.Nil {
		return nil, errors.New("roomID is required")
	}

	if theme == "" {
		return nil, errors.New("theme is required")
	}

	// 人数チェック
	if len(userIDs) <= 1 {
		return nil, errors.New("at least two users are required")
	}

	for _, userID := range userIDs {
		if userID == uuid.Nil {
			return nil, errors.New("userID is required")
		}
	}

	return &Battle{
		ID:        uuid.New(),
		RoomID:    roomID,
		StartedAt: time.Now(),
		Theme:     theme,
		UserIDs:   userIDs,
	}, nil
}
