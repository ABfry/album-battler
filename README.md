# Album Battler

## 技術スタック

### Backend

- **言語**: Go 1.25.1
- **アーキテクチャ**: クリーンアーキテクチャ（DDD）
- **データベース**: MySQL 8.0
- **ストレージ**: AWS S3
- **主要ライブラリ**:
  - AWS SDK for Go v2
  - database/sql (MySQL driver)

### Frontend

- **フレームワーク**: Next.js 16.0.1
- **言語**: TypeScript 5
- **UI ライブラリ**: React 19.2.0
- **スタイリング**: Tailwind CSS 4
- **コード品質**: ESLint, Prettier

### Infrastructure

- **コンテナ**: Docker / Docker Compose
- **IaC**: Terraform
- **クラウド**: AWS (VPC, S3)

## プロジェクト構造

```
album-battler/
├── backend/              # Go APIサーバー
│   ├── cmd/             # エントリーポイント
│   ├── internal/        # アプリケーションコード
│   │   ├── app/        # HTTPサーバー、ルーティング
│   │   ├── domain/     # ドメインロジック（エンティティ、リポジトリ、サービス）
│   │   ├── infra/      # インフラ実装（MySQL、S3、バリデーター）
│   │   └── usecase/    # ビジネスロジック
│   ├── migrations/     # DBマイグレーション
│   └── Dockerfile
│
├── frontend/            # Next.jsフロントエンド
│   ├── app/            # App Router
│   ├── public/         # 静的ファイル
│   └── Dockerfile
│
├── terraform/          # インフラ定義
│   ├── modules/        # Terraformモジュール
│   └── README.md       # インフラドキュメント
│
├── docker-compose.yml  # 開発環境定義
└── Makefile           # 開発コマンド
```

## 前提条件

### 必須ツール

- [Docker](https://www.docker.com/products/docker-desktop) & Docker Compose
- [Make](https://www.gnu.org/software/make/) (オプション)
- [Terraform](https://www.terraform.io/downloads) >= 1.0 (インフラ管理用)
- [AWS CLI](https://aws.amazon.com/cli/) (インフラ管理用)

### 推奨ツール

- [Go](https://golang.org/dl/) 1.25.1+ (ローカル開発用)
- [Node.js](https://nodejs.org/) 20+ (ローカル開発用)

## セットアップ

### 1. リポジトリのクローン

```bash
git clone https://github.com/ABfry/album-battler.git
cd album-battler
```

### 2. 環境変数の設定

Backend 環境変数ファイルを作成：

```bash
cp backend/.env.example backend/.env
```

`backend/.env` を編集して、必要な値を設定：

```env
PORT=8080

DB_DRIVER=mysql
DB_ENV=local
DB_HOST=db
DB_PORT=3306
DB_NAME=album_battler
DB_USER=album_user
DB_PASSWORD=album_pass
DB_PARAMS=charset=utf8mb4&parseTime=true&loc=Local
DB_MAX_OPEN_CONNS=32
DB_MAX_IDLE_CONNS=16
DB_CONN_MAX_LIFETIME=5400
DB_CONN_MAX_IDLE_TIME=900

AWS_REGION=ap-northeast-1
S3_BUCKET_NAME=album-battler-images
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```

Frontend 環境変数ファイルを作成：

```bash
cp frontend/.env.example frontend/.env.local
```

`frontend/.env.local` を編集して、必要な値を設定：

```env
# WebSocket Server URL
NEXT_PUBLIC_WEBSOCKET_URL=ws://localhost:8080/ws
```

### 3. Docker Compose で起動

```bash
# すべてのサービスを起動（ビルド含む）
make up

# または
docker compose up --build
```

起動するサービス：

- **Frontend**: http://localhost:3000
- **Backend**: http://localhost:8080
- **phpMyAdmin**: http://localhost:8081 (DB 管理)
- **MySQL**: localhost:3306

### 4. 動作確認

```bash
# Backendヘルスチェック
curl http://localhost:8080/health

# Frontendアクセス
open http://localhost:3000
```

## 開発ワークフロー

### Make コマンド

```bash
# すべてのサービスを起動
make up

# Goコードのリント（自動修正）
make lint
```

### Backend 開発

```bash
cd backend

# 依存関係のインストール
go mod download

# テスト実行
go test ./...

# 特定パッケージのテスト
go test ./internal/infra/storage

# リント実行
make lint
```

### Frontend 開発

```bash
cd frontend

# 依存関係のインストール
npm install

# コードフォーマット
npm run format

# リント実行
npm run lint
```

## アーキテクチャ

### Backend: クリーンアーキテクチャ

```
┌─────────────────────────────────────────┐
│           app/http (Interface)          │
│  HTTPハンドラー、ルーティング、レスポンス  │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            usecase (Use Case)           │
│      ビジネスロジック、オーケストレーション  │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│          domain (Domain Layer)          │
│  エンティティ、リポジトリIF、ドメインサービス │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      infra (Infrastructure Layer)       │
│   MySQL実装、S3実装、バリデーター実装    │
└─────────────────────────────────────────┘
```

**依存関係の方向**: 外側 → 内側（Domain 層は外部に依存しない）

### 主要コンポーネント

#### Domain Layer

- **Entity**: `User`, `Room`, `Battle`, `Image`
- **Repository**: データ永続化のインターフェース
- **Service**: ドメイン固有のロジック（画像検証、ストレージ）

#### Use Case Layer

- `UploadAlbumImageUseCase`: 画像アップロードビジネスロジック

#### Infrastructure Layer

- **MySQL**: リポジトリ実装
- **S3**: 画像ストレージ実装
- **Validator**: 画像バリデーション実装

## インフラ管理

### Terraform でのインフラ構築

詳細は [terraform/README.md](./terraform/README.md) を参照。

```bash
cd terraform

# 初期化
terraform init

# 実行計画の確認
terraform plan

# インフラ作成
terraform apply
```

## データベース管理

### phpMyAdmin

http://localhost:8081 でアクセス

- **サーバー**: db
- **ユーザー**: root
- **パスワード**: rootpass

### マイグレーション

初回起動時に `backend/migrations/init_db.sql` が自動実行されます。

## トラブルシューティング

### DB コンテナ起動失敗

```bash
# ボリュームを削除して再作成
docker compose down -v
docker compose up --build
```

### AWS 認証エラー

1. AWS 認証情報が正しく設定されているか確認
2. `.env` ファイルの `AWS_ACCESS_KEY_ID` と `AWS_SECRET_ACCESS_KEY` を確認
3. IAM ポリシーで S3 へのアクセス権限があるか確認
