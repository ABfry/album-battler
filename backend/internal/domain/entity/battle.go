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

	// ゲーム状態管理（WebSocket再接続時の状態復元用）
	CurrentPhase         string     // "selecting" | "clap_time" | "result" | "finished"
	SelectingStartedAt   *time.Time // 画像選択開始時刻
	ClapPhaseStartedAt   *time.Time // 現在の拍手フェーズ開始時刻
	ClapCurrentUserIndex *int       // 現在何人目の拍手か (0-4)
	ResultStartedAt      *time.Time // 結果フェーズ開始時刻
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
		ID:           uuid.New(),
		RoomID:       roomID,
		StartedAt:    time.Now(),
		Theme:        theme,
		UserIDs:      userIDs,
		CurrentPhase: "selecting", // 初期状態は画像選択フェーズ
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

func (b *Battle) RecordClapUserChanged(userID uuid.UUID) {
	if b == nil {
		return
	}

	b.RecordEvent(event.ChangeClapUserEvent{
		RoomID:     b.RoomID,
		BattleID:   b.ID,
		UserID:     userID,
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

func (b *Battle) RecordResultPhaseStarted(winnerUserID uuid.UUID) {
	if b == nil {
		return
	}

	b.RecordEvent(event.StartResultPhaseEvent{
		RoomID:       b.RoomID,
		BattleID:     b.ID,
		WinnerUserID: winnerUserID,
		OccurredOn:   time.Now(),
	})
}
