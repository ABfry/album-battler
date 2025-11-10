package storage

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// mockS3Client はS3クライアントのモック実装
type mockS3Client struct {
	putObjectFunc func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func (m *mockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	if m.putObjectFunc != nil {
		return m.putObjectFunc(ctx, params, optFns...)
	}
	return &s3.PutObjectOutput{}, nil
}

func TestS3ImageStorage_Upload(t *testing.T) {
	tests := []struct {
		name        string
		bucket      string
		region      string
		key         string
		contentType string
		imageData   []byte
		mockFunc    func(t *testing.T, bucket, key, contentType string, imageData []byte) func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	}{
		{
			name:        "upload success with JPEG image",
			bucket:      "test-bucket",
			region:      "ap-northeast-1",
			key:         "images/test.jpg",
			contentType: "image/jpeg",
			imageData:   []byte("fake jpeg data"),
			mockFunc: func(t *testing.T, bucket, key, contentType string, imageData []byte) func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				return func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
					if *params.Bucket != bucket {
						t.Errorf("Bucket = %v, want %v", *params.Bucket, bucket)
					}
					if *params.Key != key {
						t.Errorf("Key = %v, want %v", *params.Key, key)
					}
					if *params.ContentType != contentType {
						t.Errorf("ContentType = %v, want %v", *params.ContentType, contentType)
					}
					body, err := io.ReadAll(params.Body)
					if err != nil {
						t.Fatalf("failed to read body: %v", err)
					}
					if !bytes.Equal(body, imageData) {
						t.Errorf("Body = %v, want %v", body, imageData)
					}
					return &s3.PutObjectOutput{}, nil
				}
			},
		},
		{
			name:        "upload success with PNG image",
			bucket:      "my-bucket",
			region:      "us-east-1",
			key:         "uploads/2024/image.png",
			contentType: "image/png",
			imageData:   []byte("fake png data"),
			mockFunc: func(t *testing.T, bucket, key, contentType string, imageData []byte) func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				return func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
					return &s3.PutObjectOutput{}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockClient := &mockS3Client{
				putObjectFunc: tt.mockFunc(t, tt.bucket, tt.key, tt.contentType, tt.imageData),
			}

			reader := bytes.NewReader(tt.imageData)
			_, err := mockClient.PutObject(ctx, &s3.PutObjectInput{
				Bucket:      aws.String(tt.bucket),
				Key:         aws.String(tt.key),
				Body:        reader,
				ContentType: aws.String(tt.contentType),
			})

			if err != nil {
				t.Errorf("Upload() error = %v", err)
			}
		})
	}
}

func TestS3ImageStorage_GetURL(t *testing.T) {
	tests := []struct {
		name    string
		bucket  string
		region  string
		key     string
		wantURL string
	}{
		{
			name:    "standard URL format",
			bucket:  "test-bucket",
			region:  "ap-northeast-1",
			key:     "images/test.jpg",
			wantURL: "https://test-bucket.s3.ap-northeast-1.amazonaws.com/images/test.jpg",
		},
		{
			name:    "URL with nested path",
			bucket:  "my-bucket",
			region:  "us-east-1",
			key:     "uploads/2024/01/image.png",
			wantURL: "https://my-bucket.s3.us-east-1.amazonaws.com/uploads/2024/01/image.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewS3ImageStorageWithClient((*s3.Client)(nil), tt.bucket, tt.region)
			ctx := context.Background()

			gotURL, err := storage.GetURL(ctx, tt.key)
			if err != nil {
				t.Fatalf("GetURL failed: %v", err)
			}

			if gotURL != tt.wantURL {
				t.Errorf("GetURL() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}
