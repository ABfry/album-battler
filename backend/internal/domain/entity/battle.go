package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Battle struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	User1     uuid.UUID
	User2     uuid.UUID
	StartedAt time.Time
	Winner    uuid.UUID
}

func NewBattle(roomID, user1, user2 uuid.UUID) (*Battle, error) {
	if roomID == uuid.Nil {
		return nil, errors.New("roomID is required")
	}
	if user1 == uuid.Nil {
		return nil, errors.New("user1 is required")
	}
	if user2 == uuid.Nil {
		return nil, errors.New("user2 is required")
	}

	return &Battle{
		ID:        uuid.New(),
		RoomID:    roomID,
		User1:     user1,
		User2:     user2,
		StartedAt: time.Now(),
	}, nil
}

func (b *Battle) SetWinner(winner uuid.UUID) error {
	if winner != b.User1 && winner != b.User2 {
		return errors.New("winner must be either user1 or user2")
	}
	b.Winner = winner
	return nil
}
