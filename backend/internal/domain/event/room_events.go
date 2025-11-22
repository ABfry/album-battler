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
	return "player_join_room"
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
	return "player_leave_room"
}

func (e UserLeftRoomEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// ゲーム開始ボタンが押されたイベント（お題生成に時間がかかるため、ボタンを押したらローディング画面を表示させるために使用
type GameStartButtonPressedEvent struct {
	RoomID     uuid.UUID
	OccurredOn time.Time
}

func (e GameStartButtonPressedEvent) EventType() string {
	return "start_button_pressed"
}

func (e GameStartButtonPressedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// ゲームが開始されたイベント
type GameStartedEvent struct {
	RoomID     uuid.UUID
	RoomNumber int
	OccurredOn time.Time
}

func (e GameStartedEvent) EventType() string {
	return "start_game"
}

func (e GameStartedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// 部屋の設定が更新されたイベント
type RoomSettingsUpdatedEvent struct {
	RoomID                 uuid.UUID
	RoomNumber             int
	BattleTimeLimitSeconds int
	OccurredOn             time.Time
}

func (e RoomSettingsUpdatedEvent) EventType() string {
	return "room_settings_updated"
}

func (e RoomSettingsUpdatedEvent) OccurredAt() time.Time {
	return e.OccurredOn
}
