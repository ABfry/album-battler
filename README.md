# アルバムバトラー

アルバムの写真で戦う！AI× リアルタイム即興フォトバトル

![image](https://ptera-publish.topaz.dev/project/01KAQEKFJ9FD2HRYWHR54CD75T.png)

# 遊び方

1. ルームを作る
   ![image](https://ptera-publish.topaz.dev/project/01KAQEJVZ117A0F3E5FSNVQ0Z0.png)
2. ルームナンバーを共有して、対戦相手を募集
   ![image](https://ptera-publish.topaz.dev/project/01KAQEMGQ0SZGHT5KA2JV7SBDY.png)
3. メンバーが揃ったら、ゲームを開始
   ![image](https://ptera-publish.topaz.dev/project/01KAQEPBBN64T8KZVPB1BQTW5Q.png)
4. お題が出題されるので、自身のアルバムから探索
   ![image](https://ptera-publish.topaz.dev/project/01KAQEQ0T3N9RR3QEMRDC5BWVQ.png)
5. いい写真が見つかったら、確定
   ![image](https://ptera-publish.topaz.dev/project/01KAQEQKVNXG8XJGFNFNS5BCD4.png)
6. 他のユーザの写真を拍手で評価
   ![image](https://ptera-publish.topaz.dev/project/01KAQES7XRKQJ9H621PKFPBG47.png)
7. 結果発表
   ![image](https://ptera-publish.topaz.dev/project/01KAQESJTWBN4QKW31NC7EQP1W.png)
8. 点数の詳細や AI の評価が見れる
   ![image](https://ptera-publish.topaz.dev/project/01KAQESWWYH949YDCW3TW189JE.png)

## 技術スタック

### Backend

#### 言語・フレームワーク

- **言語**: Go 1.25.1
- **アーキテクチャ**: クリーンアーキテクチャ（DDD）+ イベント駆動アーキテクチャ
- **データベース**: MySQL 8.0
- **ストレージ**: AWS S3

#### リアルタイム通信

- **WebSocket**: 独自実装（Gorilla WebSocket ベース）
- **Hub パターン**: ルーム単位のメッセージブロードキャスト管理
- **イベント駆動**: ドメインイベント → WebSocket 配信の自動化

#### AI/ML

- **LLM**: Google Gemini 2.5 Flash
- **機能**:
  - バトルテーマ自動生成（Function Calling）
  - 画像評価（マルチモーダル + 並列処理）
- **統合**: google.golang.org/genai

#### 主要ライブラリ

- AWS SDK for Go v2
- database/sql (MySQL driver)
- google.golang.org/genai

### Frontend

#### 言語・フレームワーク

- **フレームワーク**: Next.js 16.0.1
- **言語**: TypeScript 5
- **UI ライブラリ**: React 19.2.0
- **ルーティング**: App Router

#### UI/UX

- **スタイリング**: Tailwind CSS 4
- **コンポーネント**: Radix UI
- **デザインシステム**: Storybook 8.6.14
- **フォント**: Google Fonts (Geist, Noto Sans JP), Adobe Fonts (Typekit)
- **アニメーション**: Framer Motion
- **トースト通知**: Sonner
- **テーマ管理**: next-themes
- **画像処理**: browser-image-compression, heic2any
- **Cookie 管理**: js-cookie

#### リアルタイム通信

- **WebSocket Client**: ネイティブ WebSocket API
- **状態管理**: React Context API
- **型安全**: TypeScript 完全対応

#### コード品質

- **Linter**: ESLint
- **Formatter**: Prettier
- **ビルド**: Vite 6.4.1

### Infrastructure

- **コンテナ**: Docker / Docker Compose
- **IaC**: Terraform
- **クラウド**: AWS (VPC, S3)
- **開発 DB 管理**: phpMyAdmin

![image](https://ptera-publish.topaz.dev/project/01KAP35AE82JW42Y443Q7NY4AJ.jpeg)

## 技術的工夫点

### フロントエンド

#### WebSocket 接続の共有

リアルタイムでの通信に使用する 1 つの WebSocket 接続を、アプリ全体で共有できるようにしました。Context API を用いて WebSocket コンテキストとプロバイダーを設計し、ページ遷移における接続の維持を実現しています。

```tsx
// WebSocketContext - 接続状態と操作を提供
type WebSocketContextType = {
  status: "idle" | "connecting" | "connected" | "disconnected" | "error";
  sendMessage: (data: string) => void;
  connect: () => void;
  disconnect: () => void;
  isReconnecting: boolean;
  // ...
};
```

WebSocketProvider がルートレイアウトでラップするので、ルーム画面からバトル画面へ遷移しても接続が維持されます。

#### 用途別のカスタムフック

1 つの WebSocket 接続を複数の用途で容易に使えるよう、専用のフックを用意しました。

- useWebSocketEvents: ゲーム進行のイベント購読
- useWebSocketClap: 拍手メッセージの送信

#### バトル画面のフェーズ管理

バトルの進行状態を「フェーズ」として管理し、各フェーズの開始・終了・タイムアップ時にコールバックを発火させる設計にしました。

```tsx
type BattlePhase =
  | "waiting" | "selecting"
  | "clap_time_1" | "clap_time_2" | "clap_time_3" | "clap_time_4" | "clap_time_5"
  | "result" | "finished";

// フェーズごとのイベントハンドラ
onPhaseStart?: (phase) => void;  // フェーズ開始時（演出開始など）
onPhaseEnd?: (phase) => void;    // フェーズ終了時
onTimeUp?: (phase) => void;      // 制限時間終了時
```

これにより、フェーズ遷移に伴う演出（タイマーリセット、画像切り替えなど）のトリガーが容易になりました。

#### API テストページによるデバッグ支援

開発効率を上げるため、WebSocket イベントと REST API を個別にテストできるページを用意しました。ルーム作成からゲーム開始，拍手送信まで一通りの操作をテストでき、WebSocket メッセージや API レスポンスをリアルタイムで確認できます。
これにより、バトル途中の状態からのデバッグや、特定のイベントの動作確認が容易になりました。
![image](https://ptera-publish.topaz.dev/project/01KAQAM8MFH9181423MY3HQTJD.png)

### バックエンド

通信は用途に応じて WebSocket と REST API を使い分けました。
イベント通知などリアルタイム性が必要なデータ送受信は WebSocket を使用し、
その通知を受けてクライアントは REST API を呼び出してデータを取得することで、リアルタイム性とシンプルさの両立を目指しました。
もう少し具体的に言うと、WebSocket はイベントタイプのみの「何かあったよ」という通知のみを送信し、実際のデータはこのイベントを受けてクライアントが REST によってデータの取得を行うといった感じです。

また、イベント通知などで少し複雑になりそうだったため、シーケンス図やクラス図を作成・設計によって認識の共有を行ってから実装にとりかかりました。

#### ルーム作成

![image](https://ptera-publish.topaz.dev/project/01KAN3KPCMJSBDYY3198XYJXNR.png)

#### バトル

![image](https://ptera-publish.topaz.dev/project/01KAN3M6VN2FW93K29XJJD2ZPE.png)

#### リザルト

![image](https://ptera-publish.topaz.dev/project/01KAN3MARQGQYRNTT4P5J291QB.png)

### イベント駆動 + DDD

クリーンアーキテクチャに沿って、ドメイン層を中心に据えた設計にしました。
ビジネスルールをドメインに閉じ込めることで、他の層に振り回されない構造を目指しました。

WebSocket のイベント通知周りはイベント駆動で設計しました。
EventDispatcher であるイベントタイプに関連したハンドラーを登録しておき、イベントの発火を検知してハンドラーを呼び出すようにしています。

#### クラス図(一部抜粋)

![image](https://ptera-publish.topaz.dev/project/01KAN3JVGC12TQGSX0920Y9411.jpeg)

## インフラ

0 ベースで書くのは初めてでしたが、Terraform にチャレンジしてみました。
AWS のリソースを定義して、コンテナのビルドとデプロイのみは Github Actions での自動化を行いました。

### 自動レビュー

プルリクエストを出すと AI が自動でレビューしてくれるシステムを構築しました。
見落としがちなエラーにも気づけてとても助かりました。(↓ 寝ぼけてた中助かりました
![image](https://ptera-publish.topaz.dev/project/01KAP4T5PRZ9K1W173G4RKSDVD.jpeg)

### 自動デプロイ

main ブランチに変更がマージされると、AWS へ自動でビルド・デプロイされます。![image](https://ptera-publish.topaz.dev/project/01KAQET1G8TPKR3AKNE13ZT34J.png)

## プロジェクト構造

```
album-battler/
├── backend/                      # Go APIサーバー
│   ├── cmd/                     # エントリーポイント
│   │   ├── main.go             # サーバー起動
│   │   └── test-ai/            # AI機能テストツール
│   ├── internal/                # アプリケーションコード
│   │   ├── app/                # インターフェース層
│   │   │   ├── http/          # HTTPハンドラー、ルーティング
│   │   │   └── deps.go        # 依存性注入
│   │   ├── domain/             # ドメイン層
│   │   │   ├── entity/        # エンティティ（User, Room, Battle, Image）
│   │   │   ├── event/         # ドメインイベント定義
│   │   │   ├── repository/    # リポジトリインターフェース
│   │   │   └── service/       # ドメインサービス（LLM, EventDispatcher, etc.）
│   │   ├── infra/              # インフラ層
│   │   │   ├── mysql/         # MySQLリポジトリ実装
│   │   │   ├── storage/       # S3ストレージ実装
│   │   │   ├── websocket/     # WebSocket Hub実装
│   │   │   ├── event/         # イベントハンドラー実装
│   │   │   ├── ai/            # Gemini統合
│   │   │   ├── clap/          # 拍手スケジューラー
│   │   │   └── validator/     # バリデーター実装
│   │   └── usecase/            # ユースケース層
│   │       ├── room/          # ルーム管理（作成、参加、退出、開始）
│   │       ├── battle/        # バトル管理（作成、画像提出、取得）
│   │       └── clap/          # 拍手フェーズ管理
│   ├── migrations/             # DBマイグレーション
│   └── Dockerfile
│
├── frontend/                    # Next.jsフロントエンド
│   ├── app/                    # App Router
│   │   ├── layout.tsx         # WebSocketProvider統合
│   │   ├── title/             # タイトル画面
│   │   ├── room/              # ルーム関連画面
│   │   ├── battle/            # バトル画面
│   │   ├── search/            # 検索画面
│   │   ├── api-test/          # API統合テスト画面
│   │   └── ui-demo/           # UIデモ画面
│   ├── src/
│   │   ├── components/        # 共有UIコンポーネント
│   │   │   ├── ui/           # Radix UI + Storybook対応
│   │   │   └── ApiTest/      # テスト用コンポーネント
│   │   ├── features/          # フィーチャーベース設計
│   │   │   ├── room/
│   │   │   └── api-test/
│   │   ├── hooks/             # カスタムフック
│   │   └── lib/
│   │       ├── websocket/     # WebSocket統合
│   │       └── api/           # API Client
│   ├── .storybook/            # Storybook設定
│   ├── public/                # 静的ファイル
│   └── Dockerfile
│
├── terraform/                  # インフラ定義
│   ├── modules/               # Terraformモジュール
│   └── README.md              # インフラドキュメント
│
├── docker-compose.yml         # 開発環境定義
└── Makefile                   # 開発コマンド
```

## ゲームフロー

### 基本的な流れ

```
1. ルーム作成・参加
   ↓
2. ゲーム開始（ホストのみ）
   ↓ [Gemini: テーマ自動生成]
3. 画像提出フェーズ（全プレイヤー）
   ↓ [S3: 画像保存]
4. 拍手フェーズ（評価）
   ↓ [Gemini: 画像評価]
5. 結果表示
```

### 詳細フロー

#### 1. ルーム管理フェーズ

```
[ユーザー] → POST /room
         → WebSocket接続（user_id付き）
         → POST /room/join
         ← WebSocket: player_join_room イベント（全メンバーに配信）
```

**実装機能:**

- ルーム番号自動生成（0000-9999）
- ホスト自動選定（最初の参加者）
- 最大 5 人まで参加可能
- 退出時の自動ホスト交代

#### 2. ゲーム開始フェーズ

```
[ホスト] → POST /room/{id}/start
       ↓
[Backend] → CreateBattleUseCase
         → Gemini: テーマ生成（Function Calling）
         → BattleRepository: 保存
         → ClapScheduler: 1分後のタイマー設定
         → EventDispatcher: GameStartedEvent発火
         ← WebSocket: start_game イベント（ルーム全員に配信）
```

**Gemini テーマ生成例:**

- "夏の思い出"
- "動物の可愛い瞬間"
- "都市の夜景"

#### 3. 画像提出フェーズ

```
[プレイヤー] → POST /battle/{id}/image
           ↓
[Backend] → 画像バリデーション
         → S3: 画像アップロード
         → ImageRepository: 保存
         → EventDispatcher: ImageSendEvent発火
         ← WebSocket: image_send イベント（ルーム全員に配信）
         → 全員提出判定
         → [全員提出時] ClapScheduler.TriggerNow()
```

**画像バリデーション:**

- MIME type: image/jpeg, image/png
- サイズ上限: 5MB
- 形式チェック

#### 4. 拍手フェーズ

```
[ClapScheduler] → StartClapTimeUseCase
               → RoomRepository: status更新（battling → clap_time）
               → EventDispatcher: StartClapTimeEvent発火
               ← WebSocket: start_clap_time イベント（ルーム全員に配信）
```

**トリガー条件:**

- タイマー: 画像提出開始から 1 分後
- 即時実行: 全員が画像を提出した場合

#### 5. 評価・結果フェーズ（実装予定）

```
[Gemini] → 並列画像評価（最大3並列）
        → スコア算出（0-100点）
        → 評価理由生成
        ← WebSocket: 結果配信
```

## 開発ワークフロー

### Make コマンド

```bash
# すべてのサービスを起動（フォアグラウンド）
make up

# すべてのサービスをバックグラウンドで起動
make upd

# コンテナを停止
make down

# コンテナとボリュームを停止・削除
make downv

# 各サービスのログを確認
make log

# Goコードのリント（自動修正あり）
make lint

# Goコードのリント（CI 相当で自動修正なし）
make lint-ci

# Frontend Storybook を起動
make storybook

# Storybook を停止
make storybook-stop
```

`make lint` / `make lint-ci` は backend ディレクトリで golangci-lint を実行します。`make storybook` 系コマンドは frontend 配下で Storybook を起動・停止します。

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

## WebSocket API 仕様

### 接続方法

```
WebSocket URL: ws://localhost:8080/ws?user_id={userId}
```

**パラメータ:**

- `user_id`: ユーザー ID（UUID 形式）

### メッセージタイプ

WebSocket は 5 種類のメッセージチャネルをサポート:

| チャネル              | 送信元   | 配信先             | 用途                 |
| --------------------- | -------- | ------------------ | -------------------- |
| `broadcast`           | ユーザー | 全接続クライアント | 全体ブロードキャスト |
| `systemBroadcast`     | システム | 全接続クライアント | システム通知         |
| `roomBroadcast`       | ユーザー | ルーム内メンバー   | ルーム内通信         |
| `systemRoomBroadcast` | システム | ルーム内メンバー   | システムイベント通知 |
| `userMessage`         | システム | 特定ユーザー       | 個別通知             |

### イベントタイプ例

#### 1. player_join_room

ユーザーがルームに参加した時に発火

```json
{
  "type": "player_join_room",
  "data": {
    "room_id": "uuid",
    "user_id": "uuid",
    "user_name": "string",
    "room_members": ["uuid1", "uuid2"]
  }
}
```

#### 2. player_leave_room

ユーザーがルームから退出した時に発火

```json
{
  "type": "player_leave_room",
  "data": {
    "room_id": "uuid",
    "user_id": "uuid",
    "new_host_id": "uuid or null"
  }
}
```

#### 3. start_game

ゲームが開始された時に発火

```json
{
  "type": "start_game",
  "data": {
    "room_id": "uuid",
    "battle_id": "uuid",
    "theme": "string",
    "started_at": "timestamp"
  }
}
```

#### 4. image_send

プレイヤーが画像を提出した時に発火

```json
{
  "type": "image_send",
  "data": {
    "battle_id": "uuid",
    "user_id": "uuid",
    "image_url": "string",
    "submitted_count": 3,
    "total_players": 5
  }
}
```

#### 5. start_clap_time

拍手フェーズが開始された時に発火

```json
{
  "type": "start_clap_time",
  "data": {
    "room_id": "uuid",
    "battle_id": "uuid",
    "clap_started_at": "timestamp"
  }
}
```

## AI 機能の詳細

#### 1. テーマ自動生成

**機能:** バトル開始時にテーマを自動生成

**実装方法:**

- Function Calling 使用
- 構造化出力（JSON Schema）
- 自動リトライ（最大 5 回）

**生成プロンプト例:**

```
写真対戦ゲームのテーマを1つ生成してください。
条件:
- 30文字以内
- 創造性を刺激する具体的なテーマ
- 季節や感情、シチュエーションを含む
```

**出力例:**

```json
{
  "theme": "雨上がりの虹と笑顔"
}
```

**コード位置:** `backend/internal/infra/ai/gemini.go`

#### 2. 画像評価（実装予定）

**機能:** 提出された画像をテーマに基づいて評価

**実装方法:**

- マルチモーダル対応（画像 + テキスト）
- 並列処理（最大 3 画像同時評価）
- スコアリング（0-100 点）
- 評価理由生成

**評価プロンプト例:**

```
テーマ「{theme}」に対して、この画像を評価してください。

評価基準:
- テーマとの一致度 (40点)
- 創造性 (30点)
- 写真の質 (30点)

出力形式:
- score: 0-100
- reason: 評価理由（100文字以内）
```

**出力例:**

```json
{
  "score": 85,
  "reason": "虹と笑顔が見事に捉えられており、テーマを完璧に表現しています。構図も美しい。"
}
```

**コード位置:** `backend/internal/infra/ai/gemini.go:JudgeImage()`

#### 3. エラーハンドリング

**リトライ戦略:**

- テーマ生成: 最大 5 回リトライ
- 画像評価: エラー時はスコア 0 で記録

**タイムアウト:**

- デフォルト: 30 秒
- Function Calling: 60 秒

**ログ記録:**

- トークン使用量追跡
- エラー詳細ログ

## データベーススキーマ

### ER 図

```
┌─────────────┐
│    users    │
├─────────────┤
│ id (PK)     │───┐
│ name        │   │
│ icon_url    │   │
│ hashed_pass │   │
│ created_at  │   │
└─────────────┘   │
                  │
                  │  ┌──────────────┐
                  ├──│  room_users  │ (多対多)
                  │  ├──────────────┤
                  │  │ room_id (FK) │
                  │  │ user_id (FK) │
                  │  │ joined_at    │
                  │  └──────────────┘
                  │         │
┌─────────────┐   │         │
│    rooms    │───┘         │
├─────────────┤             │
│ id (PK)     │─────────────┘
│ room_number │
│ host_user_id│ (FK → users)
│ created_at  │
│ expired_at  │
│ status      │ (waiting/full/battling/clap_time/result/closed)
│ max_users   │
└─────────────┘
      │
      │
      │  ┌─────────────┐
      └──│   battles   │
         ├─────────────┤
         │ id (PK)     │───┐
         │ room_id (FK)│   │
         │ started_at  │   │
         │ theme       │   │
         └─────────────┘   │
                           │
                           │  ┌────────────────┐
                           ├──│ battle_users   │ (多対多)
                           │  ├────────────────┤
                           │  │ battle_id (FK) │
                           │  │ user_id (FK)   │
                           │  └────────────────┘
                           │
                           │  ┌─────────────┐
                           └──│   images    │
                              ├─────────────┤
                              │ id (PK)     │
                              │ user_id (FK)│
                              │ battle_id   │
                              │ image_url   │
                              │ uploaded_at │
                              │ ai_score    │ (0-100 or NULL)
                              │ user_score  │ (0+ or NULL)
                              └─────────────┘
```

### テーブル詳細

#### users

ユーザー情報

| カラム          | 型           | 制約     | 説明             |
| --------------- | ------------ | -------- | ---------------- |
| id              | CHAR(36)     | PK       | UUID             |
| name            | VARCHAR(50)  | NOT NULL | ユーザー名       |
| icon_url        | TEXT         | NOT NULL | アイコン画像 URL |
| hashed_password | VARCHAR(255) | NOT NULL | bcrypt ハッシュ  |
| created_at      | DATETIME     | NOT NULL | 作成日時         |

#### rooms

ゲームルーム

| カラム       | 型       | 制約                    | 説明           |
| ------------ | -------- | ----------------------- | -------------- |
| id           | CHAR(36) | PK                      | UUID           |
| room_number  | INT      | NOT NULL, CHECK(0-9999) | 4 桁ルーム番号 |
| host_user_id | CHAR(36) | FK(users)               | ホストユーザー |
| status       | ENUM     | NOT NULL                | ルーム状態     |
| max_users    | INT      | NOT NULL, DEFAULT 5     | 最大参加者数   |
| created_at   | DATETIME | NOT NULL                | 作成日時       |
| expired_at   | DATETIME | NULL                    | 有効期限       |

**status 値:**

- `waiting`: 待機中
- `full`: 満員
- `battling`: バトル中（画像提出フェーズ）
- `clap_time`: 拍手フェーズ
- `result`: 結果表示中
- `closed`: 終了

#### battles

バトルセッション

| カラム     | 型           | 制約      | 説明              |
| ---------- | ------------ | --------- | ----------------- |
| id         | CHAR(36)     | PK        | UUID              |
| room_id    | CHAR(36)     | FK(rooms) | ルーム ID         |
| started_at | DATETIME     | NOT NULL  | 開始日時          |
| theme      | VARCHAR(100) | NOT NULL  | Gemini 生成テーマ |

#### images

提出画像

| カラム      | 型       | 制約               | 説明                     |
| ----------- | -------- | ------------------ | ------------------------ |
| id          | CHAR(36) | PK                 | UUID                     |
| user_id     | CHAR(36) | FK(users)          | 提出者                   |
| battle_id   | CHAR(36) | FK(battles)        | バトル ID                |
| image_url   | TEXT     | NOT NULL           | S3 URL                   |
| uploaded_at | DATETIME | NOT NULL           | アップロード日時         |
| ai_score    | FLOAT    | NULL, CHECK(0-100) | AI スコア                |
| user_score  | INT      | NULL, CHECK(0+)    | ユーザースコア（拍手数） |

#### room_users / battle_users

多対多中間テーブル（参加者管理）

## アーキテクチャ

### Backend: クリーンアーキテクチャ + イベント駆動

```
┌─────────────────────────────────────────┐
│       app/http (Interface Layer)        │
│  HTTPハンドラー、ルーティング、WebSocket  │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│        usecase (Use Case Layer)         │
│      ビジネスロジック、オーケストレーション  │
│   • Room管理 • Battle管理 • Clap管理    │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         domain (Domain Layer)           │
│  • エンティティ（User, Room, Battle）    │
│  • リポジトリIF                          │
│  • ドメインサービス                       │
│  • ドメインイベント                       │
│  • イベントディスパッチャー               │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│    infra (Infrastructure Layer)         │
│  • MySQL（リポジトリ実装）               │
│  • S3（ストレージ実装）                  │
│  • WebSocket Hub                        │
│  • EventDispatcher実装                  │
│  • Gemini統合                           │
│  • ClapScheduler                        │
└─────────────────────────────────────────┘
```

**依存関係の方向**: 外側 → 内側（Domain 層は外部に依存しない）

### イベント駆動アーキテクチャ

#### 仕組み

```
[UseCase] → Entity.RecordEvent(event)
         → UseCase完了後
         → EventDispatcher.Dispatch(events)
         → EventHandler実行
         → WebSocket配信 / DB更新 / 外部API呼び出し
```

#### 登録済みイベントハンドラー

| イベント              | ハンドラー             | 処理内容                           |
| --------------------- | ---------------------- | ---------------------------------- |
| `GameStartedEvent`    | GameStartedHandler     | WebSocket: ルーム全員に配信        |
| `UserJoinedRoomEvent` | UserJoinedRoomHandler  | WebSocket: ルームメンバーに配信    |
| `UserLeftRoomEvent`   | UserLeftRoomHandler    | WebSocket: 退出通知 + 新ホスト選定 |
| `ImageSendEvent`      | ImageSendHandler       | WebSocket: 画像提出通知            |
| `StartClapTimeEvent`  | ClapTimeStartedHandler | WebSocket: 拍手フェーズ開始通知    |

#### 利点

- **疎結合**: UseCase とインフラ層の分離
- **拡張性**: 新規イベント・ハンドラーを簡単に追加
- **テスタビリティ**: イベント記録をモック可能
- **監査**: 全ドメインイベントをログ記録可能

### WebSocket Hub パターン

#### 構造

```
                ┌─────────────┐
                │  WebSocket  │
                │     Hub     │
                └──────┬──────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
   ┌────▼────┐    ┌───▼────┐    ┌───▼────┐
   │ Client1 │    │Client2 │    │Client3 │
   │(Room A) │    │(Room A)│    │(Room B)│
   └─────────┘    └────────┘    └────────┘
```

**機能:**

- クライアント登録・解除
- ルーム単位のメッセージルーティング
- ブロードキャスト（全体/ルーム/個別）
- 接続状態管理

**コード位置:** `backend/internal/infra/websocket/hub.go`

### 非同期スケジューリング（ClapScheduler）

#### 仕組み

```
[ゲーム開始] → Schedule(1分後, battleID)
           → Goroutine起動
           → time.After(1分) 待機
           → StartClapTimeUseCase実行

[全員提出完了] → TriggerNow(battleID)
             → 即座にタイマーキャンセル
             → StartClapTimeUseCase実行
```

**特徴:**

- 時間ベーストリガー: 1 分後自動実行
- 即時トリガー: 全員提出時は待機せず実行
- スレッドセーフ: sync.Mutex で排他制御

**コード位置:** `backend/internal/infra/clap/clap_scheduler.go`

### 主要コンポーネント

#### Domain Layer

- **Entity**: `User`, `Room`, `Battle`, `Image`
- **Repository**: データ永続化のインターフェース
- **Service**:
  - `ImageValidator`: 画像検証
  - `ImageStorage`: S3 ストレージ IF
  - `LLMService`: Gemini 統合 IF
  - `EventDispatcher`: イベント配信 IF
  - `ClapScheduler`: 拍手スケジューラー IF

#### Use Case Layer

**Room 管理:**

- `CreateRoomUseCase`: ルーム作成
- `JoinRoomUseCase`: ルーム参加
- `LeaveRoomUseCase`: ルーム退出
- `GetRoomUseCase`: ルーム情報取得
- `StartGameUseCase`: ゲーム開始

**Battle 管理:**

- `CreateBattleUseCase`: バトル作成（Gemini 統合）
- `GetBattleUseCase`: バトル情報取得
- `ImageSendUseCase`: 画像提出

**Clap 管理:**

- `StartClapTimeUseCase`: 拍手フェーズ開始

#### Infrastructure Layer

- **MySQL**: リポジトリ実装（5 種類）
- **S3**: 画像ストレージ実装
- **WebSocket**: Hub + EventPublisher 実装
- **Gemini**: LLMService 実装
- **Validator**: 画像バリデーション実装
- **ClapScheduler**: 非同期タイマー実装

## インフラ管理

### Terraform でのインフラ構築（AWS VPC/S3）

基本的な VPC・S3 バケットの構築手順。

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
