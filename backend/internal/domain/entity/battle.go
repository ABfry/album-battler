package entity

import (
	"errors"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/google/uuid"
)

type Battle struct {
	event.AggregateRoot

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

func (b *Battle) RecordImageSent(userID uuid.UUID, imageURL string) {
	if b == nil {
		return
	}

	b.RecordEvent(event.ImageSendEvent{
		RoomID:     b.RoomID,
		BattleID:   b.ID,
		UserID:     userID,
		ImageURL:   imageURL,
		OccurredOn: time.Now(),
	})
}

func (b *Battle) RecordClapTimeStarted() {
	if b == nil {
		return
	}

	b.RecordEvent(event.StartClapTimeEvent{
		RoomID:     b.RoomID,
		BattleID:   b.ID,
		OccurredOn: time.Now(),
	})
}

func (b *Battle) RecordClapCounted(userID uuid.UUID, targetUserID uuid.UUID, count int) {
	if b == nil {
		return
	}

	b.RecordEvent(event.ClapSendEvent{
		RoomID:       b.RoomID,
		BattleID:     b.ID,
		UserID:       userID,
		TargetUserID: targetUserID,
		ClapCount:    count,
		OccurredOn:   time.Now(),
	})
}
