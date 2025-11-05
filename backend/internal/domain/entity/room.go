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
	ID uuid.UUID
	RoomID int
	HostUserID uuid.UUID
	CreatedAt time.Time
	ExpiredAt time.Time
	Status RoomStatus
}

func NewRoom(roomID int, hostUserID uuid.UUID, expiredAt time.Time) *Room {
	return &Room{
		ID: uuid.New(),
		RoomID: roomID,
		HostUserID: hostUserID,
		CreatedAt: time.Now(),
		ExpiredAt: expiredAt,
		Status: WaitJoin,
	}
}

func (r *Room) CreateRoom(roomID int, hostUserID uuid.UUID, expiredAt time.Time) error {
	if roomID <= 0 {
		return errors.New("roomID is invalid")
	}
	if hostUserID == uuid.Nil {
		return errors.New("hostUserID is required")
	}
	if expiredAt.IsZero() {
		return errors.New("expiredAt is required")
	} else if expiredAt.Before(time.Now()) {
		return errors.New("expiredAt must be in the future")
	}

	r.RoomID = roomID
	r.HostUserID = hostUserID
	r.ExpiredAt = expiredAt

	return nil
}

func (r *Room) ChangeStatus(status RoomStatus) {
	r.Status = status
}