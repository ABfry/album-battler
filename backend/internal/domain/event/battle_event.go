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
