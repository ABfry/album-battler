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
		return []RoomStatus{FullyJoined, InBattle}
	case FullyJoined:
		return []RoomStatus{WaitJoin, InBattle}
	case InBattle:
		return []RoomStatus{Result}
	case Result:
		return []RoomStatus{}
	case Closed:
		return []RoomStatus{}
	default:
		return []RoomStatus{}
	}
}

type Room struct {
	ID         uuid.UUID
	RoomNumber int // 1~9999
	HostUserID *uuid.UUID
	CreatedAt  time.Time
	ExpiredAt  time.Time
	Status     RoomStatus
	UserIDs    []uuid.UUID
	MaxUsers   int
}

func NewRoom(roomNumber int, expiredAt time.Time, maxUsers int) (*Room, error) {
	if roomNumber <= 0 {
		return nil, errors.New("roomNumber is invalid")
	}
	if expiredAt.IsZero() {
		return nil, errors.New("expiredAt is required")
	} else if expiredAt.Before(time.Now()) {
		return nil, errors.New("expiredAt must be in the future")
	}

	return &Room{
		ID:         uuid.New(),
		RoomNumber: roomNumber,
		HostUserID: nil,
		CreatedAt:  time.Now(),
		ExpiredAt:  expiredAt,
		Status:     WaitJoin,
		UserIDs:    []uuid.UUID{},
	}, nil
}

func (r *Room) ChangeStatus(status RoomStatus) error {
	allowedStatuses := GetValidTransitions(r.Status)

	if status == Closed {
		return errors.New("cannot change status to closed, use Dissolve() instead")
	}

	for _, allowedStatus := range allowedStatuses {
		if status == allowedStatus {
			r.Status = status
			return nil
		}
	}

	return errors.New("invalid status transition")
}

func (r *Room) AddUser(userID uuid.UUID) error {
	if len(r.UserIDs) >= r.MaxUsers {
		return errors.New("room is full")
	}

	// 1人目ならホストにする
	if len(r.UserIDs) == 0 {
		r.HostUserID = &userID
	}

	// 重複チェック
	for _, u := range r.UserIDs {
		if u == userID {
			return errors.New("user already in room")
		}
	}

	r.UserIDs = append(r.UserIDs, userID)
	return nil
}

func (r *Room) RemoveUser(userID uuid.UUID) error {
	for i, u := range r.UserIDs {
		if u == userID {
			r.UserIDs = append(r.UserIDs[:i], r.UserIDs[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

func (r *Room) IsExpired() bool {
	return r.ExpiredAt.Before(time.Now())
}

func (r *Room) IsFull() bool {
	return len(r.UserIDs) >= r.MaxUsers
}

func (r *Room) Dissolve() error {
	r.UserIDs = []uuid.UUID{}
	r.HostUserID = nil
	r.Status = Closed
	return nil
}
