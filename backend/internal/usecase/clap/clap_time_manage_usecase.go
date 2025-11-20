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

	// 各ユーザーのフェーズを順番に実行
	for i, userID := range battle.UserIDs {
		// 拍手時間待機
		time.Sleep(defaultClapDelay)

		if i == len(battle.UserIDs)-1 {
			// 最後の場合はイベント通知しない
			break
		}

		// ドメインイベントを記録
		battle.RecordClapUserChanged(userID)

		// ドメインイベントを配信
		events := battle.PopEvents()
		if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
			fmt.Println("failed to dispatch domain events", "error", err)
			return err
		}
	}

	// todo : resultUsecase呼び出す

	return nil
}
