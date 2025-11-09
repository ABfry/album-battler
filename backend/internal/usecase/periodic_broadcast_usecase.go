package usecase

import (
	"context"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/service"
)

// 定期的にメッセージをブロードキャスト
type PeriodicBroadcastUseCase struct {
	eventPublisher service.EventPublisher
	interval       time.Duration
}

func NewPeriodicBroadcastUseCase(
	eventPublisher service.EventPublisher,
	interval time.Duration,
) *PeriodicBroadcastUseCase {
	return &PeriodicBroadcastUseCase{
		eventPublisher: eventPublisher,
		interval:       interval,
	}
}

// 定期ブロードキャストを開始
func (uc *PeriodicBroadcastUseCase) Start(ctx context.Context) {
	ticker := time.NewTicker(uc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			uc.sendPeriodicMessage(ctx)
		}
	}
}

func (uc *PeriodicBroadcastUseCase) sendPeriodicMessage(ctx context.Context) {
	event := service.BroadcastEvent{
		Type: "periodic",
		Payload: map[string]interface{}{
			"message": "Hello",
		},
	}

	_ = uc.eventPublisher.BroadcastToAll(ctx, event)
}
