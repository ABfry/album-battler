package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type RoomStatus int

const (
	WaitJoin RoomStatus = iota
	FullyJoined
	InBattle
	Result
	Closed
)

type Room struct {
	ID         uuid.UUID
	RoomID     int // 1~9999
	HostUserID uuid.UUID
	CreatedAt  time.Time
	ExpiredAt  time.Time
	Status     RoomStatus
}

func NewRoom(roomID int, hostUserID uuid.UUID, expiredAt time.Time) (*Room, error) {
	if roomID <= 0 {
		return nil, errors.New("roomID is invalid")
	}
	if hostUserID == uuid.Nil {
		return nil, errors.New("hostUserID is required")
	}
	if expiredAt.IsZero() {
		return nil, errors.New("expiredAt is required")
	} else if expiredAt.Before(time.Now()) {
		return nil, errors.New("expiredAt must be in the future")
	}

	return &Room{
		ID:         uuid.New(),
		RoomID:     roomID,
		HostUserID: hostUserID,
		CreatedAt:  time.Now(),
		ExpiredAt:  expiredAt,
		Status:     WaitJoin,
	}, nil
}

func (r *Room) ChangeStatus(status RoomStatus) {
	r.Status = status
}
