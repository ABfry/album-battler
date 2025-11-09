package entity

import (
	"errors"

	"github.com/google/uuid"
)

type BattleUser struct {
	BattleID uuid.UUID
	UserID   uuid.UUID
}

func NewBattleUser(battleID, userID uuid.UUID) (*BattleUser, error) {
	if battleID == uuid.Nil {
		return nil, errors.New("battleID is required")
	}
	if userID == uuid.Nil {
		return nil, errors.New("userID is required")
	}

	return &BattleUser{
		BattleID: battleID,
		UserID:   userID,
	}, nil
}
