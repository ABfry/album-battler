package clap

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/ABfry/album-battler/backend/internal/usecase/battle"
	"github.com/google/uuid"
)

type ClapTimeManageInput struct {
	BattleID uuid.UUID
}

type ClapTimeManageUseCase struct {
	battleRepo    repository.BattleRepository
	dispatcher    service.EventDispatcher
	startResultUC *battle.StartResultUseCase
}

func NewClapTimeManageUseCase(
	battleRepo repository.BattleRepository,
	dispatcher service.EventDispatcher,
	startResultUC *battle.StartResultUseCase,
) *ClapTimeManageUseCase {
	return &ClapTimeManageUseCase{
		battleRepo:    battleRepo,
		dispatcher:    dispatcher,
		startResultUC: startResultUC,
	}
}

func (uc *ClapTimeManageUseCase) Execute(ctx context.Context, input ClapTimeManageInput) error {
	// バトルを取得
	b, err := uc.battleRepo.FindByID(ctx, input.BattleID)
	if err != nil {
		fmt.Printf("Failed to find battle: %v", err)
		return errors.New("failed to find battle")
	}

	// デフォルト5秒
	const defaultClapDelay = time.Second * 5

	// goroutineで非同期実行（HTTPレスポンスはすぐ返る）
	go func(battleID uuid.UUID, userIDs []uuid.UUID) {
		for i, userID := range userIDs {
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

			// 拍手ユーザーインデックスをインクリメント（WebSocket再接続時の状態復元用）
			nextIndex := i + 1
			b.ClapCurrentUserIndex = &nextIndex
			now := time.Now()
			b.ClapPhaseStartedAt = &now

			// バトルを保存
			if err := uc.battleRepo.Save(context.Background(), b); err != nil {
				fmt.Printf("Failed to save battle: %v\n", err)
				continue
			}

			// ドメインイベントを記録
			b.RecordClapUserChanged(userID)
			events := b.PopEvents()

			// タイムアウト付きコンテキストでDispatch実行
			dispatchCtx, cancel := context.WithTimeout(context.Background(), defaultClapDelay)
			if err := uc.dispatcher.Dispatch(dispatchCtx, events); err != nil {
				fmt.Printf("failed to dispatch domain events: %v\n", err)
			}
			cancel()

			// 拍手時間待機
			<-time.After(defaultClapDelay)
		}

		// 最終ユーザーの持ち時間相当を待ってから結果フェーズへ
		<-time.After(defaultClapDelay)

		resultCtx := context.Background()
		if err := uc.startResultUC.Execute(resultCtx, battle.StartResultInput{
			BattleID: battleID,
		}); err != nil {
			fmt.Printf("Failed to start result phase: %v\n", err)
		}
	}(b.ID, b.UserIDs)

	return nil
}
