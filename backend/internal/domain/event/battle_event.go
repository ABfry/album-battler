package event

import (
	"time"

	"github.com/google/uuid"
)

type ImageSendEvent struct {
	RoomID     uuid.UUID
	BattleID   uuid.UUID
	UserID     uuid.UUID
	ImageURL   string
	OccurredOn time.Time
}

func (e ImageSendEvent) EventType() string {
	return "image_send"
}

func (e ImageSendEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

type StartResultPhaseEvent struct {
	RoomID       uuid.UUID
	BattleID     uuid.UUID
	WinnerUserID uuid.UUID
	OccurredOn   time.Time
}

func (e StartResultPhaseEvent) EventType() string {
	return "start_result_phase"
}

func (e StartResultPhaseEvent) OccurredAt() time.Time {
	return e.OccurredOn
}
