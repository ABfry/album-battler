package battle

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

type ImageSendInput struct {
	BattleID    uuid.UUID
	UserID      uuid.UUID
	ImageBase64 []byte
	ContentType string
}

type ImageSendUseCase struct {
	battleRepo     repository.BattleRepository
	imageRepo      repository.ImageRepository
	imageValidator service.ImageValidator
	imageStorage   service.ImageStorage
	dispatcher     service.EventDispatcher
	clapScheduler  service.ClapScheduler
}

func NewImageSendUseCase(
	battleRepo repository.BattleRepository,
	imageRepo repository.ImageRepository,
	imageValidator service.ImageValidator,
	imageStorage service.ImageStorage,
	dispatcher service.EventDispatcher,
	clapScheduler service.ClapScheduler,
) *ImageSendUseCase {
	return &ImageSendUseCase{
		battleRepo:     battleRepo,
		imageRepo:      imageRepo,
		imageValidator: imageValidator,
		imageStorage:   imageStorage,
		dispatcher:     dispatcher,
		clapScheduler:  clapScheduler,
	}
}

func (uc *ImageSendUseCase) Execute(ctx context.Context, input ImageSendInput) error {
	// 画像のバリデーション
	if err := uc.imageValidator.ValidateFormat(input.ImageBase64); err != nil {
		fmt.Printf("Invalid image format: %v", err)
		return errors.New("invalid image format")
	}
	if err := uc.imageValidator.ValidateSize(int64(len(input.ImageBase64))); err != nil {
		fmt.Printf("Image size exceeds limit: %v", err)
		return errors.New("image size exceeds limit")
	}

	// 一意なキーを生成(yyyymmdd/uuid.{ContentType})
	key := fmt.Sprintf("%s/%s.%s", time.Now().Format("20060102"), uuid.New().String(), mimeToExtension(input.ContentType))

	// 画像の保存
	imageURL, err := uc.imageStorage.Upload(ctx, key, bytes.NewReader(input.ImageBase64), input.ContentType)
	if err != nil {
		fmt.Printf("Failed to upload image: %v", err)
		return errors.New("failed to upload image")
	}

	image, err := entity.NewImage(input.UserID, input.BattleID, imageURL)
	if err != nil {
		fmt.Printf("Failed to create image entity: %v", err)
		return errors.New("failed to create image entity")
	}

	battle, err := uc.battleRepo.FindByID(ctx, input.BattleID)
	if err != nil {
		fmt.Printf("Failed to find battle: %v", err)
		return errors.New("failed to find battle")
	}
	if battle == nil {
		fmt.Printf("Battle not found: %s", input.BattleID)
		return errors.New("battle not found")
	}

	// battle内のユーザか確認
	if !slices.Contains(battle.UserIDs, input.UserID) {
		fmt.Printf("User %s is not in battle %s", input.UserID, input.BattleID)
		return errors.New("user is not in battle")
	}

	// ドメインイベントを記録
	battle.RecordImageSent(input.UserID, imageURL)

	// 画像情報の保存
	if err := uc.imageRepo.Save(ctx, image); err != nil {
		fmt.Printf("Failed to save image info: %v", err)
		return errors.New("failed to save image info")
	}

	// ドメインイベントを配信
	events := battle.PopEvents()
	if err := uc.dispatcher.Dispatch(ctx, events); err != nil {
		fmt.Println("failed to dispatch domain events", "error", err)
		return err
	}

	// 投稿した画像数を集計
	submittedCount, err := uc.imageRepo.CountDistinctUsersByBattleID(ctx, battle.ID)
	if err != nil {
		fmt.Printf("Failed to count submitted images: %v", err)
		return errors.New("failed to count submissions")
	}

	// デバッグ
	// fmt.Println("submittedCount", submittedCount)
	// fmt.Println("battle.UserIDs", len(battle.UserIDs))
	// fmt.Println("clapScheduler", uc.clapScheduler != nil)

	// 投稿した画像数がユーザー数と同じか、拍手スケジューラが設定されている場合は拍手フェーズへ移行
	if submittedCount >= len(battle.UserIDs) && uc.clapScheduler != nil {
		if err := uc.clapScheduler.TriggerNow(ctx, battle.ID); err != nil {
			fmt.Printf("Failed to start clap time: %v", err)
			return errors.New("failed to start clap time")
		}
	}

	return nil
}

func mimeToExtension(mime string) string {
	mime = strings.ToLower(strings.TrimSpace(mime))
	switch mime {
	case "image/jpeg", "image/jpg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/gif":
		return "gif"
	case "image/webp":
		return "webp"
	case "image/heic":
		return "heic"
	default:
		return "bin"
	}
}
