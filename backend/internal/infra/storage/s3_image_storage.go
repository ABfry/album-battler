package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var _ service.ImageStorage = (*S3ImageStorage)(nil)

// AWS S3を使用した画像ストレージの実装
type S3ImageStorage struct {
	client *s3.Client
	bucket string
	region string
}

// S3ImageStorageインスタンスを作成する
func NewS3ImageStorage(ctx context.Context, bucket, region string) (*S3ImageStorage, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3ImageStorage{
		client: client,
		bucket: bucket,
		region: region,
	}, nil
}

// カスタムS3クライアントを使用してインスタンスを作成する
// テストやDI用
func NewS3ImageStorageWithClient(client *s3.Client, bucket, region string) *S3ImageStorage {
	return &S3ImageStorage{
		client: client,
		bucket: bucket,
		region: region,
	}
}

// 画像をS3にアップロードし、URLを返す
func (s *S3ImageStorage) Upload(ctx context.Context, key string, data io.Reader, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        data,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}

	// S3オブジェクトのURLを生成
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)
	return url, nil
}

// 画像のURLを取得する
func (s *S3ImageStorage) GetURL(ctx context.Context, key string) (string, error) {
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)
	return url, nil
}
