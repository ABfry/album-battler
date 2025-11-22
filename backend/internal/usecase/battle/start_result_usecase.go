package battle

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

type StartResultInput struct {
	BattleID uuid.UUID
}

type StartResultUseCase struct {
	battleRepo  repository.BattleRepository
	roomRepo    repository.RoomRepository
	imageRepo   repository.ImageRepository
	clapCounter service.ClapCounter
	dispatcher  service.EventDispatcher
}

func NewStartResultUseCase(
	battleRepo repository.BattleRepository,
	roomRepo repository.RoomRepository,
	imageRepo repository.ImageRepository,
	clapCounter service.ClapCounter,
	dispatcher service.EventDispatcher,
) *StartResultUseCase {
	return &StartResultUseCase{
		battleRepo:  battleRepo,
		roomRepo:    roomRepo,
		imageRepo:   imageRepo,
		clapCounter: clapCounter,
		dispatcher:  dispatcher,
	}
}

func (uc *StartResultUseCase) Execute(ctx context.Context, input StartResultInput) error {
	// Battleを取得
	battle, err := uc.battleRepo.FindByID(ctx, input.BattleID)
	if err != nil {
		fmt.Printf("Failed to find battle: %v", err)
		return errors.New("failed to find battle")
	}
	if battle == nil {
		fmt.Printf("Battle not found: %s", input.BattleID)
		return errors.New("battle not found")
	}

	// Roomを取得
	room, err := uc.roomRepo.FindByID(ctx, battle.RoomID)
	if err != nil {
		fmt.Printf("Failed to find room: %v", err)
		return errors.New("failed to find room")
	}
	if room == nil {
		fmt.Printf("Room not found: %s", battle.RoomID)
		return errors.New("room not found")
	}

	// 拍手カウンターから最終UserScoreを取得
	clapSnapshot := uc.clapCounter.Snapshot(input.BattleID)

	// 全画像を取得してUserScoreを更新
	images, err := uc.imageRepo.FindImagesByBattleID(ctx, input.BattleID)
	if err != nil {
		fmt.Printf("Failed to find images: %v", err)
		return errors.New("failed to find images")
	}

	// UserScoreを更新してDBに保存
	for _, img := range images {
		if clapImg, ok := clapSnapshot[img.UserID]; ok {
			if err := img.SetScore(img.AIScore, clapImg.UserScore); err != nil {
				fmt.Printf("Failed to set score for image %s: %v", img.ID, err)
				continue
			}
		}

		if err := uc.imageRepo.Save(ctx, img); err != nil {
			fmt.Printf("Failed to save image %s: %v", img.ID, err)
			continue
		}
	}

	// 最終スコア計算と勝者決定
	type ImageScore struct {
		UserID     uuid.UUID
		FinalScore float64
		UploadedAt time.Time
	}

	scores := make([]ImageScore, 0, len(images))
	for _, img := range images {
		finalScore := img.AIScore + float64(img.UserScore)
		scores = append(scores, ImageScore{
			UserID:     img.UserID,
			FinalScore: finalScore,
			UploadedAt: img.UploadedAt,
		})
	}

	// スコアでソート（降順、同点なら投稿時刻の早い順）
	sort.Slice(scores, func(i, j int) bool {
		if scores[i].FinalScore != scores[j].FinalScore {
			return scores[i].FinalScore > scores[j].FinalScore
		}
		return scores[i].UploadedAt.Before(scores[j].UploadedAt)
	})

	var winnerUserID uuid.UUID
	if len(scores) > 0 {
		winnerUserID = scores[0].UserID
	}

	// RoomステータスをResultに変更
	if err := room.ChangeStatus(entity.Result); err != nil {
		fmt.Printf("Failed to change room status: %v", err)
		return errors.New("failed to change room status")
	}

	if err := uc.roomRepo.Save(ctx, room); err != nil {
		fmt.Printf("Failed to save room: %v", err)
		return errors.New("failed to save room")
	}

	// バトルの状態を結果フェーズに更新（WebSocket再接続時の状態復元用）
	battle.CurrentPhase = "result"
	now := time.Now()
	battle.ResultStartedAt = &now

	// バトルを保存
	if err := uc.battleRepo.Save(ctx, battle); err != nil {
		fmt.Printf("Failed to save battle: %v", err)
		return errors.New("failed to save battle")
	}

	// ドメインイベントを記録
	battle.RecordResultPhaseStarted(winnerUserID)

	// ドメインイベントを配信
	events := battle.PopEvents()
	if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
		fmt.Println("failed to dispatch domain events", "error", err)
		return err
	}

	// 拍手カウンターをリセット
	uc.clapCounter.Reset(input.BattleID)

	fmt.Printf("Result phase started for battle %s, winner: %s\n", input.BattleID, winnerUserID)
	return nil
}
