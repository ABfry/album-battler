package clap

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

type ClapInput struct {
	BattleID uuid.UUID
	UserID   uuid.UUID
	Count    int
}

type ClapUseCase struct {
	roomRepo    repository.RoomRepository
	battleRepo  repository.BattleRepository
	clapCounter service.ClapCounter
	dispatcher  service.EventDispatcher
}

func NewClapUseCase(roomRepo repository.RoomRepository, battleRepo repository.BattleRepository, clapCounter service.ClapCounter, dispatcher service.EventDispatcher) *ClapUseCase {
	return &ClapUseCase{
		roomRepo:    roomRepo,
		battleRepo:  battleRepo,
		clapCounter: clapCounter,
		dispatcher:  dispatcher,
	}
}

func (uc *ClapUseCase) Execute(ctx context.Context, input ClapInput) error {
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
		fmt.Printf("User not found in battle: %v", input.UserID)
		return errors.New("user not found in battle")
	}

	// 拍手数を追加
	uc.clapCounter.Add(input.BattleID, input.UserID, input.Count)

	// ドメインイベントを記録
	battle.RecordEvent(event.ClapSendEvent{
		BattleID:   input.BattleID,
		UserID:     input.UserID,
		ClapCount:  input.Count,
		OccurredOn: time.Now(),
	})

	// イベントをディスパッチ
	events := battle.PopEvents()
	if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
		fmt.Println("failed to dispatch domain events", "error", err)
		return err
	}

	return nil
}
