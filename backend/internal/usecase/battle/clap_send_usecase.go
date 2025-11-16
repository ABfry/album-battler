package battle

import (
	"context"
	"errors"
	"fmt"

	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

type ClapSendInput struct {
	BattleID uuid.UUID
	UserID   uuid.UUID
	Count    int
}

type ClapSendUseCase struct {
	battleRepo repository.BattleRepository
	dispatcher service.EventDispatcher
}

func NewClapSendUseCase(battleRepo repository.BattleRepository, dispatcher service.EventDispatcher) *ClapSendUseCase {
	return &ClapSendUseCase{
		battleRepo: battleRepo,
		dispatcher: dispatcher,
	}
}

func (uc *ClapSendUseCase) Execute(ctx context.Context, input ClapSendInput) error {
	// 拍手数が不正なら弾く
	if input.Count <= 0 {
		return errors.New("count must be greater than zero")
	}

	// バトルを取得
	battle, err := uc.battleRepo.FindByID(ctx, input.BattleID)
	if err != nil {
		fmt.Printf("Failed to find battle: %v", err)
		return errors.New("failed to find battle")
	}

	// ユーザがいるか確認
	found := false
	for _, userID := range battle.UserIDs {
		if userID == input.UserID {
			found = true
			break
		}
	}
	if !found {
		return errors.New("user not found in battle")
	}

	// 拍手数を記録
	battle.RecordClapSent(input.UserID, input.Count)

	events := battle.PopEvents()
	if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
		fmt.Println("failed to dispatch domain events", "error", err)
		return err
	}

	return nil
}
