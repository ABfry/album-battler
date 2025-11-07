# 画像ストレージアーキテクチャ

## 概要

このプロジェクトでは、画像ファイルの保存にS3を使用しており、クリーンアーキテクチャ/DDDの原則に従って設計されています。

## アーキテクチャの特徴

### 1. 責務の分離

**Repository (DB) と ImageStorage (S3) を分離**

- `Repository`: ドメインエンティティの永続化（CRUD操作）に責務を限定
- `ImageStorage`: ファイルストレージ操作に責務を限定

```
❌ 避けるべき設計:
AlbumRepository に SaveImage メソッドを追加
→ 責務が混在し、Single Responsibility Principle に反する

✅ 推奨する設計:
AlbumRepository (DB用) と ImageStorage (S3用) を独立させる
→ 各インターフェースが単一の責務を持つ
```

### 2. ディレクトリ構造

```
backend/
├── internal/
│   ├── domain/                    # ドメイン層（ビジネスロジック）
│   │   ├── entity/
│   │   │   ├── user.go
│   │   │   └── album.go           # Albumエンティティ
│   │   ├── repository/
│   │   │   ├── user_repository.go
│   │   │   └── album_repository.go # DB操作のインターフェース
│   │   └── service/
│   │       └── image_storage.go    # ストレージ操作のインターフェース
│   │
│   ├── usecase/                   # ユースケース層（アプリケーションロジック）
│   │   └── album_usecase.go       # RepositoryとStorageを組み合わせて使う
│   │
│   └── infra/                     # インフラ層（具体的な実装）
│       ├── persistence/
│       │   └── postgres_album_repo.go  # AlbumRepositoryのPostgreSQL実装
│       └── storage/
│           └── s3_image_storage.go     # ImageStorageのS3実装
```

## インターフェース定義

### ImageStorage インターフェース

`domain/service/image_storage.go` に定義されています。

```go
type ImageStorage interface {
    // 画像をアップロードし、URLを返す
    Upload(ctx context.Context, key string, data io.Reader, contentType string) (string, error)

    // 画像を削除する
    Delete(ctx context.Context, key string) error

    // 画像のURLを取得する（署名付きURLなど）
    GetURL(ctx context.Context, key string) (string, error)
}
```

### AlbumRepository インターフェース

`domain/repository/album_repository.go` に定義されています。

```go
type AlbumRepository interface {
    Create(ctx context.Context, album *entity.Album) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Album, error)
    Update(ctx context.Context, album *entity.Album) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context) ([]*entity.Album, error)
}
```

## 実装例

### S3ImageStorage の実装

`infra/storage/s3_image_storage.go`

- AWS SDK v2 を使用
- 通常のURL生成と署名付きURL生成の両方に対応
- テスト用にクライアントを注入可能

### AlbumUsecase での使用例

`usecase/album_usecase.go`

```go
type AlbumUsecase struct {
    albumRepo    repository.AlbumRepository  // DB操作
    imageStorage service.ImageStorage         // ファイル操作
}

func (u *AlbumUsecase) CreateAlbum(ctx context.Context, input CreateAlbumInput) (*entity.Album, error) {
    // 1. エンティティ作成
    album := entity.NewAlbum(input.Name, input.Description)

    // 2. 画像をS3にアップロード
    imageKey := fmt.Sprintf("albums/%s/cover.jpg", album.ID.String())
    imageURL, err := u.imageStorage.Upload(ctx, imageKey, input.ImageData, input.ContentType)
    if err != nil {
        return nil, err
    }

    // 3. エンティティに画像URLをセット
    album.UpdateImage(imageURL)

    // 4. DBに保存
    if err := u.albumRepo.Create(ctx, album); err != nil {
        // ロールバック: アップロードした画像を削除
        _ = u.imageStorage.Delete(ctx, imageKey)
        return nil, err
    }

    return album, nil
}
```

## メリット

### 1. テスタビリティ

- モックに簡単に差し替え可能
- ユニットテストで実際のS3を使わなくて済む

```go
// テスト用のモック実装
type MockImageStorage struct {
    uploadFunc func(ctx context.Context, key string, data io.Reader, contentType string) (string, error)
}

func TestCreateAlbum(t *testing.T) {
    mockStorage := &MockImageStorage{
        uploadFunc: func(...) (string, error) {
            return "https://example.com/image.jpg", nil
        },
    }

    usecase := NewAlbumUsecase(mockRepo, mockStorage)
    // テスト実行
}
```

### 2. 柔軟性

将来、S3からCloudinaryや他のサービスに変更する場合：

1. `infra/storage/cloudinary_image_storage.go` を実装
2. `ImageStorage` インターフェースを満たすようにする
3. DIコンテナの設定を変更するだけ

→ ドメイン層やユースケース層の変更は不要

### 3. ドメインの純粋性

- ドメイン層は「画像をどこかに保存する」という抽象的な概念だけを知っている
- S3、GCS、Azure Blobなどの具体的な実装には依存しない
- ビジネスロジックが技術的な詳細から分離される

## 依存関係

```
go.mod に追加されたパッケージ:
- github.com/aws/aws-sdk-go-v2
- github.com/aws/aws-sdk-go-v2/config
- github.com/aws/aws-sdk-go-v2/service/s3
```

## 環境設定

S3ImageStorageを使用する場合、以下の環境変数またはAWS認証情報が必要です：

```bash
AWS_REGION=ap-northeast-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
S3_BUCKET_NAME=your-bucket-name
```

または、IAMロールを使用する場合は環境変数不要です。

## トランザクション管理

画像アップロードとDB保存は別々の操作なので、適切なエラーハンドリングが重要です：

1. **画像アップロード失敗**: 処理を中断し、エラーを返す
2. **DB保存失敗**: アップロードした画像を削除（ロールバック）

```go
// DB保存に失敗した場合、アップロードした画像を削除
if err := u.albumRepo.Create(ctx, album); err != nil {
    _ = u.imageStorage.Delete(ctx, imageKey)
    return nil, fmt.Errorf("failed to create album: %w", err)
}
```

## 参考資料

- Clean Architecture by Robert C. Martin
- Domain-Driven Design by Eric Evans
- AWS SDK for Go v2 Documentation
