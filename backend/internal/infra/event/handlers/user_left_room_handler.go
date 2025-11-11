package handlers

import (
	"context"
	"log"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
)

type UserLeftRoomHandler struct {
	roomManager    service.RoomManager
	eventPublisher service.EventPublisher
}

func NewUserLeftRoomHandler(
	roomManager service.RoomManager,
	eventPublisher service.EventPublisher,
) *UserLeftRoomHandler {
	return &UserLeftRoomHandler{
		roomManager:    roomManager,
		eventPublisher: eventPublisher,
	}
}

// UserLeftRoomEventを処理
func (h *UserLeftRoomHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	e, ok := evt.(event.UserLeftRoomEvent)
	if !ok {
		log.Printf("Invalid event type: expected UserLeftRoomEvent, got %T", evt)
		return nil
	}

	// WebSocket部屋から退出
	h.roomManager.LeaveRoom(e.UserID, e.RoomID)

	// WebSocket通知を部屋のメンバーに送信 (通知のみ、詳細はREST APIで取得)
	broadcastEvent := service.BroadcastEvent{
		Type: "user_left",
		Payload: map[string]interface{}{
			"room_id": e.RoomID.String(), // 最小限の情報のみ
		},
	}

	if err := h.eventPublisher.PublishToRoom(ctx, e.RoomID, broadcastEvent); err != nil {
		log.Printf("Failed to publish user_left event: %v", err)
		// WebSocket送信エラーは処理を止めない
	}

	log.Printf("User %s left room %s (was_host: %v, dissolved: %v)",
		e.UserID, e.RoomID, e.WasHost, e.RoomDissolved)

	return nil
}
