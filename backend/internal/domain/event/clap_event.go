package event

import (
	"time"

	"github.com/google/uuid"
)

type ClapSendEvent struct {
	RoomID       uuid.UUID
	BattleID     uuid.UUID
	UserID       uuid.UUID // 拍手した人
	TargetUserID uuid.UUID // 拍手された画像の投稿者
	ClapCount    int
	OccurredOn   time.Time
}

func (e ClapSendEvent) EventType() string {
	return "clap_send"
}

func (e ClapSendEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

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
