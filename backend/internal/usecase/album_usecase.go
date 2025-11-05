package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/google/uuid"
)

// AlbumUsecase はアルバムに関するユースケースを実装する
type AlbumUsecase struct {
	albumRepo    repository.AlbumRepository
	imageStorage service.ImageStorage
}

// NewAlbumUsecase は新しいAlbumUsecaseを作成する
func NewAlbumUsecase(
	albumRepo repository.AlbumRepository,
	imageStorage service.ImageStorage,
) *AlbumUsecase {
	return &AlbumUsecase{
		albumRepo:    albumRepo,
		imageStorage: imageStorage,
	}
}

// CreateAlbumInput はアルバム作成の入力パラメータ
type CreateAlbumInput struct {
	Name        string
	Description string
	ImageData   io.Reader
	ContentType string // "image/jpeg", "image/png" など
}

// CreateAlbum は画像付きアルバムを作成する
// このメソッドが、DB保存(Repository)とファイル保存(ImageStorage)を組み合わせている例
func (u *AlbumUsecase) CreateAlbum(ctx context.Context, input CreateAlbumInput) (*entity.Album, error) {
	// 1. ドメインエンティティを作成
	album := entity.NewAlbum(input.Name, input.Description)

	// 2. 画像をS3にアップロード
	imageKey := fmt.Sprintf("albums/%s/cover.jpg", album.ID.String())
	imageURL, err := u.imageStorage.Upload(ctx, imageKey, input.ImageData, input.ContentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	// 3. アルバムエンティティに画像URLをセット
	if err := album.UpdateImage(imageURL); err != nil {
		// アップロードは成功したが、エンティティの更新に失敗した場合、画像を削除
		_ = u.imageStorage.Delete(ctx, imageKey)
		return nil, fmt.Errorf("failed to update album image: %w", err)
	}

	// 4. アルバム情報をDBに保存
	if err := u.albumRepo.Create(ctx, album); err != nil {
		// DB保存に失敗した場合、アップロードした画像を削除
		_ = u.imageStorage.Delete(ctx, imageKey)
		return nil, fmt.Errorf("failed to create album: %w", err)
	}

	return album, nil
}

// UpdateAlbumImage はアルバムの画像を更新する
func (u *AlbumUsecase) UpdateAlbumImage(
	ctx context.Context,
	albumID uuid.UUID,
	imageData io.Reader,
	contentType string,
) error {
	// 1. アルバムを取得
	album, err := u.albumRepo.FindByID(ctx, albumID)
	if err != nil {
		return fmt.Errorf("failed to find album: %w", err)
	}

	// 2. 古い画像のキーを保存（ロールバック用）
	oldImageURL := album.ImageURL

	// 3. 新しい画像をアップロード
	imageKey := fmt.Sprintf("albums/%s/cover.jpg", albumID.String())
	newImageURL, err := u.imageStorage.Upload(ctx, imageKey, imageData, contentType)
	if err != nil {
		return fmt.Errorf("failed to upload new image: %w", err)
	}

	// 4. アルバムの画像URLを更新
	if err := album.UpdateImage(newImageURL); err != nil {
		return fmt.Errorf("failed to update album image: %w", err)
	}

	// 5. DBに保存
	if err := u.albumRepo.Update(ctx, album); err != nil {
		// 失敗した場合、古い画像URLに戻す必要があるかもしれない
		return fmt.Errorf("failed to update album: %w", err)
	}

	// 6. （オプション）古い画像を削除
	// 注意: 複数のアルバムが同じ画像を参照している可能性がある場合は削除しないこと
	if oldImageURL != "" {
		// エラーは無視（削除失敗しても処理は成功とみなす）
		_ = u.imageStorage.Delete(ctx, imageKey)
	}

	return nil
}

// DeleteAlbum はアルバムと関連する画像を削除する
func (u *AlbumUsecase) DeleteAlbum(ctx context.Context, albumID uuid.UUID) error {
	// 1. アルバムを取得
	album, err := u.albumRepo.FindByID(ctx, albumID)
	if err != nil {
		return fmt.Errorf("failed to find album: %w", err)
	}

	// 2. DBからアルバムを削除
	if err := u.albumRepo.Delete(ctx, albumID); err != nil {
		return fmt.Errorf("failed to delete album: %w", err)
	}

	// 3. 画像を削除（失敗しても処理は成功とする）
	if album.ImageURL != "" {
		imageKey := fmt.Sprintf("albums/%s/cover.jpg", albumID.String())
		_ = u.imageStorage.Delete(ctx, imageKey)
	}

	return nil
}

// GetAlbum はアルバムを取得する
func (u *AlbumUsecase) GetAlbum(ctx context.Context, albumID uuid.UUID) (*entity.Album, error) {
	album, err := u.albumRepo.FindByID(ctx, albumID)
	if err != nil {
		return nil, fmt.Errorf("failed to find album: %w", err)
	}
	return album, nil
}

// ListAlbums は全てのアルバムを取得する
func (u *AlbumUsecase) ListAlbums(ctx context.Context) ([]*entity.Album, error) {
	albums, err := u.albumRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list albums: %w", err)
	}
	return albums, nil
}
