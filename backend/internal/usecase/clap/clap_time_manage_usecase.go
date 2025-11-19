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
	battleRepo     repository.BattleRepository
	clapsScheduler service.ClapWaitScheduler
	dispatcher     service.EventDispatcher
}

func NewClapTimeManageUseCase(
	battleRepo repository.BattleRepository,
	clapScheduler service.ClapWaitScheduler,
	dispatcher service.EventDispatcher,
) *ClapTimeManageUseCase {
	return &ClapTimeManageUseCase{
		battleRepo:     battleRepo,
		clapsScheduler: clapScheduler,
		dispatcher:     dispatcher,
	}
}

func (uc *ClapTimeManageUseCase) Execute(ctx context.Context, input ClapTimeManageInput) error {
	// バトルを取得
	battle, err := uc.battleRepo.FindByID(ctx, input.BattleID)
	if err != nil {
		fmt.Printf("Failed to find battle: %v", err)
		return errors.New("failed to find battle")
	}

	if uc.clapsScheduler == nil {
		fmt.Printf("Clap scheduler is not configured")
		return errors.New("clap scheduler is not configured")
	}

	// 拍手時間をスケジュール(デフォルト5秒)
	const defaultClapDelay = time.Second * 5
	for _, userID := range battle.UserIDs {
		isEnd := make(chan bool)
		uc.clapsScheduler.Schedule(battle.ID, userID, defaultClapDelay, isEnd)

		<-isEnd

		// ドメインイベントを記録
		battle.RecordClapUserChanged(userID)

		// ドメインイベントを配信
		events := battle.PopEvents()
		if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
			fmt.Println("failed to dispatch domain events", "error", err)
			return err
		}

		close(isEnd)
	}

	return nil
}
