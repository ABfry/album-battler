package entity

import (
	"time"

	"github.com/google/uuid"
)

type Audience struct {
	ID        uuid.UUID
	BattleID  uuid.UUID
	CreatedAt time.Time
	Score     int
}

func NewAudience(battleID uuid.UUID) *Audience {
	return &Audience{
		ID:        uuid.New(),
		BattleID:  battleID,
		CreatedAt: time.Now(),
	}
}

func (a *Audience) SetScore(points int) {
	a.Score = points
}
