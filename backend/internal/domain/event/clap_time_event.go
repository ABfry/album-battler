package event

import (
	"time"

	"github.com/google/uuid"
)

type StartClapTimeEvent struct {
	RoomID     uuid.UUID
	BattleID   uuid.UUID
	OccurredOn time.Time
}

func (e StartClapTimeEvent) EventType() string {
	return "start_clap_time"
}

func (e StartClapTimeEvent) OccurredAt() time.Time {
	return e.OccurredOn
}
