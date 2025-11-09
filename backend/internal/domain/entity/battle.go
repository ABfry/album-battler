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
	Users     []User
}

func NewBattle(roomID uuid.UUID, users []User) (*Battle, error) {
	if roomID == uuid.Nil {
		return nil, errors.New("roomID is required")
	}

	// 人数チェック
	if len(users) <= 1 {
		return nil, errors.New("at least two users are required")
	}

	return &Battle{
		ID:        uuid.New(),
		RoomID:    roomID,
		StartedAt: time.Now(),
		Users:     users,
	}, nil
}
