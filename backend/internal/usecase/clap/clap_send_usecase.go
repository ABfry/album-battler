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
	BattleID     uuid.UUID
	UserID       uuid.UUID // 拍手した人
	TargetUserID uuid.UUID // 拍手された画像の投稿者
	Count        int
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
	if input.Count <= 0 {
		return errors.New("count must be greater than zero")
	} else if input.Count > 10 {
		input.Count = 10 // 上限10拍手
	}

	// バトルを取得
	battle, err := uc.battleRepo.FindByID(ctx, input.BattleID)
	if err != nil {
		fmt.Printf("Failed to find battle: %v", err)
		return errors.New("failed to find battle")
	}

	// 拍手した人がバトル参加者か確認
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

	// 拍手された人（ターゲット）がバトル参加者か確認
	targetFound := false
	for _, userID := range battle.UserIDs {
		if userID == input.TargetUserID {
			targetFound = true
			break
		}
	}
	if !targetFound {
		fmt.Printf("Target user not found in battle: %v", input.TargetUserID)
		return errors.New("target user not found in battle")
	}

	// 拍手数を追加（拍手される側で集計）
	if _, err := uc.clapCounter.Add(input.BattleID, input.TargetUserID, input.Count); err != nil {
		fmt.Printf("Failed to add clap count: %v", err)
		return errors.New("failed to add clap count")
	}

	// ドメインイベントを記録
	battle.RecordClapCounted(input.UserID, input.TargetUserID, input.Count)

	// イベントをディスパッチ
	events := battle.PopEvents()
	if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
		fmt.Println("failed to dispatch domain events", "error", err)
		return err
	}

	return nil
}
