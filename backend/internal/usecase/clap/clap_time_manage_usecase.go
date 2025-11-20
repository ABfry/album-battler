package clap

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

type ClapTimeManageInput struct {
	BattleID uuid.UUID
}

type ClapTimeManageUseCase struct {
	battleRepo repository.BattleRepository
	dispatcher service.EventDispatcher
}

func NewClapTimeManageUseCase(
	battleRepo repository.BattleRepository,
	dispatcher service.EventDispatcher,
) *ClapTimeManageUseCase {
	return &ClapTimeManageUseCase{
		battleRepo: battleRepo,
		dispatcher: dispatcher,
	}
}

func (uc *ClapTimeManageUseCase) Execute(ctx context.Context, input ClapTimeManageInput) error {
	// バトルを取得
	battle, err := uc.battleRepo.FindByID(ctx, input.BattleID)
	if err != nil {
		fmt.Printf("Failed to find battle: %v", err)
		return errors.New("failed to find battle")
	}

	// デフォルト5秒
	const defaultClapDelay = time.Second * 5

	// goroutineで非同期実行（HTTPレスポンスはすぐ返る）
	go func(battleID uuid.UUID, userIDs []uuid.UUID) {
		for i, userID := range userIDs {
			// 拍手時間待機
			<-time.After(defaultClapDelay)

			// 最後のユーザーはイベント通知不要
			if i >= len(userIDs)-1 {
				continue
			}

			// バトルを再取得してイベント記録
			b, err := uc.battleRepo.FindByID(context.Background(), battleID)
			if err != nil {
				fmt.Printf("Failed to find battle in goroutine: %v\n", err)
				continue
			}

			// ドメインイベントを記録
			b.RecordClapUserChanged(userID)
			events := b.PopEvents()

			// タイムアウト付きコンテキストでDispatch実行
			dispatchCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := uc.dispatcher.Dispatch(dispatchCtx, events); err != nil {
				fmt.Printf("failed to dispatch domain events: %v\n", err)
			}
		}

	}(battle.ID, battle.UserIDs)
	// TODO: StartResultUseCaseを呼ぶ

	return nil
}
