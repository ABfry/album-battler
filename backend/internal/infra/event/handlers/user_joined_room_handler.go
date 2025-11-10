package handlers

import (
	"context"
	"log"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
)

type UserJoinedRoomHandler struct {
	roomManager    service.RoomManager
	eventPublisher service.EventPublisher
}

func NewUserJoinedRoomHandler(
	roomManager service.RoomManager,
	eventPublisher service.EventPublisher,
) *UserJoinedRoomHandler {
	return &UserJoinedRoomHandler{
		roomManager:    roomManager,
		eventPublisher: eventPublisher,
	}
}

// UserJoinedRoomEventを処理
func (h *UserJoinedRoomHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	e, ok := evt.(event.UserJoinedRoomEvent)
	if !ok {
		log.Printf("Invalid event type: expected UserJoinedRoomEvent, got %T", evt)
		return nil
	}

	// WebSocket部屋に参加
	h.roomManager.JoinRoom(e.UserID, e.RoomID)

	// WebSocket通知を部屋のメンバーに送信 (通知のみ、詳細はREST APIで取得)
	broadcastEvent := service.BroadcastEvent{
		Type: "user_joined",
		Payload: map[string]interface{}{
			"room_id": e.RoomID.String(), // 最小限の情報のみ
		},
	}

	if err := h.eventPublisher.PublishToRoom(ctx, e.RoomID, broadcastEvent); err != nil {
		log.Printf("Failed to publish user_joined event: %v", err)
		// WebSocket送信エラーは処理を止めない
	}

	log.Printf("User %s joined room %s (host: %v)", e.UserID, e.RoomID, e.IsHost)

	return nil
}
