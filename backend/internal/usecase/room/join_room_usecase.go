package room

import (
	"context"
	"errors"
	"fmt"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

type JoinRoomInput struct {
	UserID     uuid.UUID
	RoomNumber int
}

type JoinRoomUseCase struct {
	roomRepo   repository.RoomRepository
	dispatcher service.EventDispatcher
}

func NewJoinRoomUseCase(
	roomRepo repository.RoomRepository,
	dispatcher service.EventDispatcher,
) *JoinRoomUseCase {
	return &JoinRoomUseCase{
		roomRepo:   roomRepo,
		dispatcher: dispatcher,
	}
}

func (uc *JoinRoomUseCase) Execute(ctx context.Context, input JoinRoomInput) (*uuid.UUID, error) {
	// 部屋を取得
	room, err := uc.roomRepo.FindByRoomNumber(ctx, input.RoomNumber)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, errors.New("room not found")
	}

	// 期限チェック
	if room.IsExpired() {
		return nil, errors.New("room has expired")
	}

	// 満員チェック
	if room.IsFull() {
		return nil, errors.New("room is full")
	}

	// ルーム状態チェック
	if room.Status != entity.WaitJoin {
		return nil, errors.New("room is not waiting for join")
	}

	// ユーザーを追加 (ドメインロジック + イベント記録)
	if err := room.AddUser(input.UserID); err != nil {
		return nil, err
	}

	// 永続化
	if err := uc.roomRepo.Save(ctx, room); err != nil {
		return nil, err
	}

	// ドメインイベントを配信
	events := room.PopEvents()
	if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
		fmt.Println("failed to dispatch domain events", "error", err)
		return nil, err
	}

	return &room.ID, nil
}
