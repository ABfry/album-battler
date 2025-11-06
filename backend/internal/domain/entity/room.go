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

func GetValidTransitions(status RoomStatus) []RoomStatus {
	switch status {
	case WaitJoin:
		return []RoomStatus{FullyJoined, InBattle, Closed}
	case FullyJoined:
		return []RoomStatus{WaitJoin, InBattle, Closed}
	case InBattle:
		return []RoomStatus{Result, Closed}
	case Result:
		return []RoomStatus{Closed}
	case Closed:
		return []RoomStatus{}
	default:
		return []RoomStatus{}
	}
}

type Room struct {
	ID         uuid.UUID
	RoomNumber int // 1~9999
	HostUserID uuid.UUID
	CreatedAt  time.Time
	ExpiredAt  time.Time
	Status     RoomStatus
}

func NewRoom(roomNumber int, hostUserID uuid.UUID, expiredAt time.Time) (*Room, error) {
	if roomNumber <= 0 {
		return nil, errors.New("roomNumber is invalid")
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
		RoomNumber: roomNumber,
		HostUserID: hostUserID,
		CreatedAt:  time.Now(),
		ExpiredAt:  expiredAt,
		Status:     WaitJoin,
	}, nil
}

func (r *Room) ChangeStatus(status RoomStatus) error {
	allowedStatuses := GetValidTransitions(r.Status)
	for _, allowedStatus := range allowedStatuses {
		if status == allowedStatus {
			r.Status = status
			return nil
		}
	}

	return errors.New("invalid status transition")
}
