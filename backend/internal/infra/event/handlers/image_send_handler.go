package handlers

import (
	"context"
	"log"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
)

type ImageSendHandler struct {
	eventPublisher service.EventPublisher
}

func NewImageSendHandler(
	eventPublisher service.EventPublisher,
) *ImageSendHandler {
	return &ImageSendHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *ImageSendHandler) Handle(ctx context.Context, evt event.DomainEvent) error {
	e, ok := evt.(event.ImageSendEvent)
	if !ok {
		log.Printf("Invalid event type: expected ImageSendEvent, got %T", evt)
		return nil
	}

	// WebSocket通知を部屋のメンバーに送信
	broadcastEvent := service.BroadcastEvent{
		Type:    e.EventType(),
		Payload: map[string]interface{}{},
	}

	if err := h.eventPublisher.PublishToRoom(ctx, e.RoomID, broadcastEvent); err != nil {
		log.Printf("Failed to publish image_send event: %v", err)
		// WebSocket送信エラーは処理を止めない
	}

	log.Printf("Image sent by user %s in room %s", e.UserID, e.RoomID)

	return nil
}
