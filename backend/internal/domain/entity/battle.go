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
}

func NewBattle(roomID uuid.UUID) (*Battle, error) {
	if roomID == uuid.Nil {
		return nil, errors.New("roomID is required")
	}

	return &Battle{
		ID:        uuid.New(),
		RoomID:    roomID,
		StartedAt: time.Now(),
	}, nil
}
