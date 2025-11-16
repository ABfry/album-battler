package clap

import (
	"context"
	"errors"
	"fmt"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

type StartClapTimeInput struct {
	BattleID uuid.UUID
}

type StartClapTimeUseCase struct {
	battleRepo repository.BattleRepository
	roomRepo   repository.RoomRepository
	dispatcher service.EventDispatcher
}

func NewStartClapTimeUseCase(
	battleRepo repository.BattleRepository,
	roomRepo repository.RoomRepository,
	dispatcher service.EventDispatcher,
) *StartClapTimeUseCase {
	return &StartClapTimeUseCase{
		battleRepo: battleRepo,
		roomRepo:   roomRepo,
		dispatcher: dispatcher,
	}
}

func (uc *StartClapTimeUseCase) Execute(ctx context.Context, input StartClapTimeInput) error {
	// fmt.Println("StartClapTimeUseCase Execute", input.BattleID)
	// battleを取得
	battle, err := uc.battleRepo.FindByID(ctx, input.BattleID)
	if err != nil {
		fmt.Printf("Failed to find battle: %v", err)
		return errors.New("failed to find battle")
	}
	if battle == nil {
		fmt.Printf("Battle not found: %s", input.BattleID)
		return errors.New("battle not found")
	}

	// roomを取得
	room, err := uc.roomRepo.FindByID(ctx, battle.RoomID)
	if err != nil {
		fmt.Printf("Failed to find room: %v", err)
		return errors.New("failed to find room")
	}
	if room == nil {
		fmt.Printf("Room not found: %s", battle.RoomID)
		return errors.New("room not found")
	}

	// roomのステータスを更新
	if err := room.ChangeStatus(entity.ClapTime); err != nil {
		fmt.Printf("Failed to change room status: %v", err)
		return errors.New("failed to change room status")
	}

	// roomを保存
	if err := uc.roomRepo.Save(ctx, room); err != nil {
		fmt.Printf("Failed to save room: %v", err)
		return errors.New("failed to save room")
	}

	// ドメインイベントを記録
	battle.RecordClapTimeStarted()

	// ドメインイベントを配信
	events := battle.PopEvents()
	if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
		fmt.Println("failed to dispatch domain events", "error", err)
		return err
	}

	return nil
}
