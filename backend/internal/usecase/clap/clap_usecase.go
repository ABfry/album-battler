package clap

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
	roomRepo    repository.RoomRepository
	battleRepo  repository.BattleRepository
	clapCounter service.ClapCounter
	dispatcher  service.EventDispatcher
}

func NewClapSendUseCase(roomRepo repository.RoomRepository, battleRepo repository.BattleRepository, clapCounter service.ClapCounter, dispatcher service.EventDispatcher) *ClapSendUseCase {
	return &ClapSendUseCase{
		roomRepo:    roomRepo,
		battleRepo:  battleRepo,
		clapCounter: clapCounter,
		dispatcher:  dispatcher,
	}
}

func (uc *ClapSendUseCase) Execute(ctx context.Context, input ClapSendInput) error {
	// 拍手数が不正なら弾く
	if input.Count <= 0 || input.Count > 10 {
		return errors.New("count must be greater than zero and less than or equal to 10")
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
		fmt.Printf("User not found in battle: %v", input.UserID)
		return errors.New("user not found in battle")
	}

	// 拍手数を追加
	if _, err := uc.clapCounter.Add(input.BattleID, input.UserID, input.Count); err != nil {
		fmt.Printf("Failed to add clap count: %v", err)
		return errors.New("failed to add clap count")
	}

	// ドメインイベントを記録
	battle.RecordClapCounted(input.UserID, input.Count)

	// イベントをディスパッチ
	events := battle.PopEvents()
	if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
		fmt.Println("failed to dispatch domain events", "error", err)
		return err
	}

	return nil
}
