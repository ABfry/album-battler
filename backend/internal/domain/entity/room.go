package entity

import (
	"errors"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/google/uuid"
)

type RoomStatus int

const (
	WaitJoin RoomStatus = iota
	FullyJoined
	InBattle
	ClapTime
	Result
	Closed
)

// roomStatusStrings は RoomStatus から文字列表現への変換マップ
var roomStatusStrings = map[RoomStatus]string{
	WaitJoin:    "waiting",
	FullyJoined: "full",
	InBattle:    "battling",
	ClapTime:    "clap_time",
	Result:      "result",
	Closed:      "closed",
}

// stringToRoomStatus は文字列から RoomStatus への変換マップ
var stringToRoomStatus = map[string]RoomStatus{
	"waiting":   WaitJoin,
	"full":      FullyJoined,
	"battling":  InBattle,
	"clap_time": ClapTime,
	"result":    Result,
	"closed":    Closed,
}

// String はデバッグやログ出力用の文字列表現を返す
func (s RoomStatus) String() string {
	if str, ok := roomStatusStrings[s]; ok {
		return str
	}
	return "unknown"
}

// MarshalJSON は JSON エンコード時に文字列として出力する
func (s RoomStatus) MarshalJSON() ([]byte, error) {
	str := s.String()
	if str == "unknown" {
		return nil, errors.New("invalid room status")
	}
	return []byte(`"` + str + `"`), nil
}

// UnmarshalJSON は JSON デコード時に文字列から RoomStatus へ変換する
func (s *RoomStatus) UnmarshalJSON(data []byte) error {
	// クォートを除去
	str := string(data)
	if len(str) < 2 || str[0] != '"' || str[len(str)-1] != '"' {
		return errors.New("invalid room status format")
	}
	str = str[1 : len(str)-1]

	status, ok := stringToRoomStatus[str]
	if !ok {
		return errors.New("unknown room status: " + str)
	}
	*s = status
	return nil
}

func GetValidTransitions(status RoomStatus) []RoomStatus {
	switch status {
	case WaitJoin:
		return []RoomStatus{FullyJoined, InBattle}
	case FullyJoined:
		return []RoomStatus{WaitJoin, InBattle}
	case InBattle:
		return []RoomStatus{ClapTime}
	case ClapTime:
		return []RoomStatus{Result}
	case Result:
		return []RoomStatus{}
	case Closed:
		return []RoomStatus{}
	default:
		return []RoomStatus{}
	}
}

type Room struct {
	event.AggregateRoot // ドメインイベント記録機能を埋め込み

	ID                      uuid.UUID
	RoomNumber              int // 1~9999
	HostUserID              *uuid.UUID
	CreatedAt               time.Time
	ExpiredAt               time.Time
	Status                  RoomStatus
	UserIDs                 []uuid.UUID
	MaxUsers                int
	BattleTimeLimitSeconds  int // 30~300秒
}

func NewRoom(roomNumber int, expiredAt time.Time, maxUsers int) (*Room, error) {
	if roomNumber <= 0 {
		return nil, errors.New("roomNumber is invalid")
	}
	if expiredAt.IsZero() {
		return nil, errors.New("expiredAt is required")
	} else if expiredAt.Before(time.Now()) {
		return nil, errors.New("expiredAt must be in the future")
	}

	return &Room{
		ID:                     uuid.New(),
		RoomNumber:             roomNumber,
		HostUserID:             nil,
		CreatedAt:              time.Now(),
		ExpiredAt:              expiredAt,
		Status:                 WaitJoin,
		UserIDs:                []uuid.UUID{},
		MaxUsers:               maxUsers,
		BattleTimeLimitSeconds: 60, // デフォルト1分
	}, nil
}

func (r *Room) ChangeStatus(status RoomStatus) error {
	allowedStatuses := GetValidTransitions(r.Status)

	if status == Closed {
		return errors.New("cannot change status to closed, use Dissolve() instead")
	}

	for _, allowedStatus := range allowedStatuses {
		if status == allowedStatus {
			r.Status = status
			return nil
		}
	}

	return errors.New("invalid status transition")
}

func (r *Room) AddUser(userID uuid.UUID) error {
	if len(r.UserIDs) >= r.MaxUsers {
		return errors.New("room is full")
	}

	// 重複チェック
	for _, u := range r.UserIDs {
		if u == userID {
			return errors.New("user already in room")
		}
	}

	// 1人目ならホストにする
	isHost := len(r.UserIDs) == 0
	if isHost {
		r.HostUserID = &userID
	}

	r.UserIDs = append(r.UserIDs, userID)

	// 満員になったら自動的にFullyJoinedへ遷移
	if len(r.UserIDs) >= r.MaxUsers {
		r.Status = FullyJoined
	}

	// ドメインイベント記録
	r.RecordEvent(event.UserJoinedRoomEvent{
		RoomID:     r.ID,
		RoomNumber: r.RoomNumber,
		UserID:     userID,
		IsHost:     isHost,
		OccurredOn: time.Now(),
	})

	return nil
}

func (r *Room) RemoveUser(userID uuid.UUID) error {
	// ユーザーが存在するか確認
	found := false
	for i, u := range r.UserIDs {
		if u == userID {
			r.UserIDs = append(r.UserIDs[:i], r.UserIDs[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return errors.New("user not found")
	}

	// 退出者がホストだったか
	wasHost := r.HostUserID != nil && *r.HostUserID == userID

	// 新しいホストの選定
	var newHostID *uuid.UUID
	if wasHost && len(r.UserIDs) > 0 {
		// 残っているユーザーの中から最初のユーザーを新ホストにする
		newHostID = &r.UserIDs[0]
		r.HostUserID = newHostID
	} else if wasHost {
		// 誰もいなくなった場合
		r.HostUserID = nil
	}

	// 満員状態から空きができたらWaitJoinに戻す
	if r.Status == FullyJoined && len(r.UserIDs) < r.MaxUsers {
		r.Status = WaitJoin
	}

	// 部屋が解散したか (全員退出)
	roomDissolved := len(r.UserIDs) == 0

	// ドメインイベント記録
	r.RecordEvent(event.UserLeftRoomEvent{
		RoomID:        r.ID,
		RoomNumber:    r.RoomNumber,
		UserID:        userID,
		WasHost:       wasHost,
		NewHostID:     newHostID,
		RoomDissolved: roomDissolved,
		OccurredOn:    time.Now(),
	})

	return nil
}

func (r *Room) IsExpired() bool {
	return r.ExpiredAt.Before(time.Now())
}

func (r *Room) IsFull() bool {
	return len(r.UserIDs) >= r.MaxUsers
}

func (r *Room) IsActive() bool {
	return r.Status == WaitJoin || r.Status == FullyJoined || r.Status == InBattle
}

func (r *Room) Dissolve() error {
	r.UserIDs = []uuid.UUID{}
	r.HostUserID = nil
	r.Status = Closed
	return nil
}

func (r *Room) UpdateBattleTimeLimit(seconds int) error {
	if seconds < 30 || seconds > 300 {
		return errors.New("battle time limit must be between 30 and 300 seconds")
	}
	r.BattleTimeLimitSeconds = seconds

	// ドメインイベント記録
	r.RecordEvent(event.RoomSettingsUpdatedEvent{
		RoomID:                 r.ID,
		RoomNumber:             r.RoomNumber,
		BattleTimeLimitSeconds: r.BattleTimeLimitSeconds,
		OccurredOn:             time.Now(),
	})

	return nil
}

func (r *Room) StartGame() error {
	// ホストが存在することを確認
	if r.HostUserID == nil {
		return errors.New("no host in room")
	}

	// 状態遷移の検証
	if err := r.ChangeStatus(InBattle); err != nil {
		return err
	}

	// ドメインイベント記録
	r.RecordEvent(event.GameStartedEvent{
		RoomID:     r.ID,
		RoomNumber: r.RoomNumber,
		OccurredOn: time.Now(),
	})

	return nil
}
