package event

import (
	"time"

	"github.com/google/uuid"
)

// ユーザーが部屋に参加したイベント
type UserJoinedRoomEvent struct {
	RoomID     uuid.UUID
	RoomNumber int
	UserID     uuid.UUID
	IsHost     bool // ホストかどうか
	OccurredOn time.Time
}

func (e UserJoinedRoomEvent) EventType() string {
	return "user_joined_room"
}

func (e UserJoinedRoomEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// ユーザーが部屋から退出したイベント
type UserLeftRoomEvent struct {
	RoomID        uuid.UUID
	RoomNumber    int
	UserID        uuid.UUID
	WasHost       bool       // 退出者がホストだったか
	NewHostID     *uuid.UUID // 新しいホスト (ホストが退出し、他にメンバーがいる場合)
	RoomDissolved bool       // 部屋が解散したか (全員退出)
	OccurredOn    time.Time
}

func (e UserLeftRoomEvent) EventType() string {
	return "user_left_room"
}

func (e UserLeftRoomEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// ゲームが開始されたイベント
type GameStartedEvent struct {
	RoomID     uuid.UUID
	RoomNumber int
	OccurredOn time.Time
}

func (e GameStartedEvent) EventType() string {
	return "game_started"
}

func (e GameStartedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}
