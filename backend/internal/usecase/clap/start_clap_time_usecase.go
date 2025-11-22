package clap

import (
	"context"
	"errors"
	"fmt"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/ABfry/album-battler/backend/internal/domain/service/llm"
	"github.com/google/uuid"
)

type StartClapTimeInput struct {
	BattleID uuid.UUID
}

type StartClapTimeUseCase struct {
	battleRepo       repository.BattleRepository
	roomRepo         repository.RoomRepository
	imageRepo        repository.ImageRepository
	dispatcher       service.EventDispatcher
	llmClient        llm.LLMClient
	clapTimeManageUC *ClapTimeManageUseCase
}

func NewStartClapTimeUseCase(
	battleRepo repository.BattleRepository,
	roomRepo repository.RoomRepository,
	imageRepo repository.ImageRepository,
	dispatcher service.EventDispatcher,
	llmClient llm.LLMClient,
	clapTimeManageUC *ClapTimeManageUseCase,
) *StartClapTimeUseCase {
	return &StartClapTimeUseCase{
		battleRepo:       battleRepo,
		roomRepo:         roomRepo,
		imageRepo:        imageRepo,
		dispatcher:       dispatcher,
		llmClient:        llmClient,
		clapTimeManageUC: clapTimeManageUC,
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

	// go func() {
	// 	for _, userID := range battle.UserIDs {
	// 		fmt.Printf("Clap time started for user %s in battle %s\n", userID, battle.ID)
	// 	}
	// }

	// 非同期で画像採点をする
	go func(battleID uuid.UUID) {
		// 親contextから独立
		judgeCtx := context.Background()
		if err := uc.JudgeImageAsync(judgeCtx, battleID); err != nil {
			fmt.Printf("Failed to judge images asynchronously for battle %s: %v\n", battleID, err)
		}
	}(input.BattleID)

	// 拍手時間管理を開始
	if err := uc.clapTimeManageUC.Execute(ctx, ClapTimeManageInput(input)); err != nil {
		fmt.Printf("Failed to start clap time management for battle %s: %v\n", input.BattleID, err)
		return err
	}

	return nil
}

func (uc *StartClapTimeUseCase) JudgeImageAsync(ctx context.Context, battleID uuid.UUID) error {
	// バトル情報を取得（テーマ取得のため）
	battle, err := uc.battleRepo.FindByID(ctx, battleID)
	if err != nil {
		return fmt.Errorf("failed to find battle: %w", err)
	}
	if battle == nil {
		return fmt.Errorf("battle not found: %s", battleID)
	}

	// 提出された画像一覧を取得
	images, err := uc.imageRepo.FindImagesByBattleID(ctx, battleID)
	if err != nil {
		return fmt.Errorf("failed to find images: %w", err)
	}

	if len(images) == 0 {
		fmt.Printf("No images to judge for battle: %s\n", battleID)
		return nil
	}

	// 画像URL一覧を抽出
	imageURLs := make([]string, len(images))
	for i, img := range images {
		imageURLs[i] = img.ImageURL
	}

	// LLMクライアントで画像を評価（並列実行）
	judgeReq := &llm.ImageJudgeRequest{
		Theme:     battle.Theme,
		ImageURLs: imageURLs,
	}

	judgeResp, err := uc.llmClient.JudgeImage(ctx, judgeReq)
	if err != nil {
		return fmt.Errorf("failed to judge images: %w", err)
	}

	// 各画像のAIScoreを更新してDB保存
	for i, result := range judgeResp.Results {
		img, err := uc.imageRepo.FindByID(ctx, images[i].ID)
		if err != nil {
			fmt.Printf("Warning: failed to find image %s: %v\n", images[i].ID, err)
			continue
		} else if img == nil {
			fmt.Printf("Warning: image not found: %s\n", images[i].ID)
			continue
		}

		// スコアを設定
		if err := img.SetAIScore(float64(result.Score)); err != nil {
			fmt.Printf("Warning: failed to set score for image %s: %v\n", img.ID, err)
			continue
		}
		// AIの評価説明文を設定
		img.SetAIExplanation(result.Reason)

		// DBに保存
		if err := uc.imageRepo.Save(ctx, img); err != nil {
			fmt.Printf("Warning: failed to save image %s: %v\n", img.ID, err)
			continue
		}
	}

	return nil
}
