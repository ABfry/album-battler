package room

import (
	"context"
	"fmt"

	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

type LeaveRoomInput struct {
	UserID uuid.UUID
	RoomID uuid.UUID
}

type LeaveRoomUseCase struct {
	roomRepo   repository.RoomRepository
	dispatcher service.EventDispatcher
}

func NewLeaveRoomUseCase(
	roomRepo repository.RoomRepository,
	dispatcher service.EventDispatcher,
) *LeaveRoomUseCase {
	return &LeaveRoomUseCase{
		roomRepo:   roomRepo,
		dispatcher: dispatcher,
	}
}

func (uc *LeaveRoomUseCase) Execute(ctx context.Context, input LeaveRoomInput) error {
	// 部屋を取得
	room, err := uc.roomRepo.FindByID(ctx, input.RoomID)
	if err != nil {
		return err
	}
	if room == nil {
		return fmt.Errorf("room not found")
	}

	// ユーザーを削除 (ドメインロジック + イベント記録)
	if err := room.RemoveUser(input.UserID); err != nil {
		return err
	}

	// 永続化
	if err := uc.roomRepo.Save(ctx, room); err != nil {
		return err
	}

	// ドメインイベントを配信
	events := room.PopEvents()
	if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
		fmt.Println("failed to dispatch domain events", "error", err)
		return err
	}

	return nil
}
