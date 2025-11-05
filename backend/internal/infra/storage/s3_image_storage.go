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

// S3ImageStorage はAWS S3を使用した画像ストレージの実装
type S3ImageStorage struct {
	client *s3.Client
	bucket string
	region string
}

// NewS3ImageStorage は新しいS3ImageStorageインスタンスを作成する
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

// NewS3ImageStorageWithClient はカスタムS3クライアントを使用してインスタンスを作成する
// テストやDI時に使用
func NewS3ImageStorageWithClient(client *s3.Client, bucket, region string) *S3ImageStorage {
	return &S3ImageStorage{
		client: client,
		bucket: bucket,
		region: region,
	}
}

// Upload は画像をS3にアップロードし、URLを返す
func (s *S3ImageStorage) Upload(ctx context.Context, key string, data io.Reader, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        data,
		ContentType: aws.String(contentType),
		// ACLをpublic-readに設定する場合はコメントアウトを外す
		// ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}

	// S3オブジェクトのURLを生成
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)
	return url, nil
}

// Delete はS3から画像を削除する
func (s *S3ImageStorage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete image from S3: %w", err)
	}

	return nil
}

// GetURL は画像のURLを取得する
// 署名付きURLが必要な場合は、s3.NewPresignClientを使用する
func (s *S3ImageStorage) GetURL(ctx context.Context, key string) (string, error) {
	// 通常のURLを返す場合
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)
	return url, nil

	// 署名付きURL（有効期限付き）が必要な場合は以下のようにする:
	// presignClient := s3.NewPresignClient(s.client)
	// presignedReq, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
	// 	Bucket: aws.String(s.bucket),
	// 	Key:    aws.String(key),
	// }, s3.WithPresignExpires(15*time.Minute))
	// if err != nil {
	// 	return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	// }
	// return presignedReq.URL, nil
}

// インターフェースの実装を確認するためのコンパイル時チェック
var _ service.ImageStorage = (*S3ImageStorage)(nil)
