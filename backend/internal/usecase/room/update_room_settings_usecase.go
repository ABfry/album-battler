package room

import (
	"context"
	"errors"

	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

type UpdateRoomSettingsInput struct {
	RoomID                 uuid.UUID
	UserID                 uuid.UUID // ホスト検証用
	BattleTimeLimitSeconds *int      // nilの場合は更新しない
}

type UpdateRoomSettingsUseCase struct {
	roomRepo   repository.RoomRepository
	dispatcher service.EventDispatcher
}

func NewUpdateRoomSettingsUseCase(
	roomRepo repository.RoomRepository,
	dispatcher service.EventDispatcher,
) *UpdateRoomSettingsUseCase {
	return &UpdateRoomSettingsUseCase{
		roomRepo:   roomRepo,
		dispatcher: dispatcher,
	}
}

func (uc *UpdateRoomSettingsUseCase) Execute(ctx context.Context, input UpdateRoomSettingsInput) error {
	// 部屋を取得
	room, err := uc.roomRepo.FindByID(ctx, input.RoomID)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("room not found")
	}

	// ホストであることを確認
	if room.HostUserID == nil || *room.HostUserID != input.UserID {
		return errors.New("only host can update room settings")
	}

	// バトル制限時間の更新
	if input.BattleTimeLimitSeconds != nil {
		if err := room.UpdateBattleTimeLimit(*input.BattleTimeLimitSeconds); err != nil {
			return err
		}
	}

	// 永続化
	if err := uc.roomRepo.Save(ctx, room); err != nil {
		return err
	}

	// ドメインイベントを配信
	events := room.PopEvents()
	if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
		return err
	}

	return nil
}
