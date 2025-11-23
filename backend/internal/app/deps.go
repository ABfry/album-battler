package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	domainEvent "github.com/ABfry/album-battler/backend/internal/domain/event"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/ABfry/album-battler/backend/internal/domain/service/llm"
	"github.com/ABfry/album-battler/backend/internal/infra/ai"
	battleinfra "github.com/ABfry/album-battler/backend/internal/infra/battle"
	clapinfra "github.com/ABfry/album-battler/backend/internal/infra/clap"
	"github.com/ABfry/album-battler/backend/internal/infra/event"
	"github.com/ABfry/album-battler/backend/internal/infra/event/handlers"
	"github.com/ABfry/album-battler/backend/internal/infra/mysql"
	"github.com/ABfry/album-battler/backend/internal/infra/storage"
	"github.com/ABfry/album-battler/backend/internal/infra/validator"
	"github.com/ABfry/album-battler/backend/internal/infra/websocket"
	"github.com/ABfry/album-battler/backend/internal/usecase/battle"
	"github.com/ABfry/album-battler/backend/internal/usecase/clap"
	"github.com/ABfry/album-battler/backend/internal/usecase/room"
	"github.com/ABfry/album-battler/backend/internal/usecase/user"
)

// -- 依存関係の定義 --

type Dependencies struct {
	db *sql.DB

	// Repository
	RoomRepository       repository.RoomRepository
	BattleRepository     repository.BattleRepository
	BattleUserRepository repository.BattleUserRepository
	ImageRepository      repository.ImageRepository
	UserRepository       repository.UserRepository

	// WebSocket関連
	WebSocketHub   *websocket.Hub
	EventPublisher service.EventPublisher
	RoomManager    service.RoomManager

	// Event関連
	EventDispatcher service.EventDispatcher

	// AI関連
	LLMClient llm.LLMClient

	// Image関連
	ImageValidator service.ImageValidator
	ImageStorage   service.ImageStorage

	ImageSubmissionScheduler service.ImageSubmissionScheduler

	// Usecase
	CreateRoomUseCase        *room.CreateRoomUseCase
	JoinRoomUseCase          *room.JoinRoomUseCase
	LeaveRoomUseCase         *room.LeaveRoomUseCase
	StartGameUseCase         *room.StartGameUseCase
	GetRoomUseCase           *room.GetRoomUseCase
	UpdateRoomSettingsUseCase *room.UpdateRoomSettingsUseCase

	CreateBattleUseCase *battle.CreateBattleUseCase
	GetBattleUseCase    *battle.GetBattleUseCase
	GetBattleIDUseCase  *battle.GetBattleIDUseCase
	GetImageUseCase     *battle.GetImageUseCase
	ImageSendUseCase    *battle.ImageSendUseCase
	GetResultUseCase    *battle.GetResultUseCase

	ClapSendUseCase       *clap.ClapSendUseCase
	StartClapTimeUseCase  *clap.StartClapTimeUseCase
	ClapTimeManageUseCase *clap.ClapTimeManageUseCase
	StartResultUseCase    *battle.StartResultUseCase
	CreateUserUseCase     *user.CreateUserUseCase
}

// NewDependencies は依存関係を初期化する
func NewDependencies() (*Dependencies, error) {
	db, err := initDatabase()
	if err != nil {
		return nil, fmt.Errorf("initialize database: %w", err)
	}

	deps := &Dependencies{
		db: db,
	}

	// リポジトリの初期化
	if err := initRepositories(deps); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("initialize repositories: %w", err)
	}

	// WebSocket関連の初期化
	if err := initWebSocket(deps); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("initialize websocket: %w", err)
	}

	// Event関連の初期化
	if err := initEvents(deps); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("initialize events: %w", err)
	}

	// AI関連の初期化
	if err := initLLM(deps); err != nil {
		_ = deps.Close()
		return nil, fmt.Errorf("initialize llm client: %w", err)
	}

	// Image関連の初期化
	if err := initImage(deps); err != nil {
		_ = deps.Close()
		return nil, fmt.Errorf("initialize image services: %w", err)
	}

	// Usecase関連の初期化
	if err := initUseCases(deps); err != nil {
		_ = deps.Close()
		return nil, fmt.Errorf("initialize usecases: %w", err)
	}

	return deps, nil
}

// WebSocket関連の初期化
func initWebSocket(deps *Dependencies) error {
	// Hub作成（接続管理）
	deps.WebSocketHub = websocket.NewHub()

	// EventPublisher作成（依存性逆転）
	deps.EventPublisher = websocket.NewWebSocketEventPublisher(deps.WebSocketHub)

	// RoomManager作成（WebSocket部屋管理）
	deps.RoomManager = websocket.NewWebSocketRoomManager(deps.WebSocketHub)

	return nil
}

// Event関連の初期化
func initEvents(deps *Dependencies) error {
	// EventDispatcher作成
	dispatcher := event.NewEventDispatcher()
	deps.EventDispatcher = dispatcher

	// イベントハンドラーを登録
	dispatcherImpl := dispatcher.(*event.EventDispatcherImpl)

	userJoinedHandler := handlers.NewUserJoinedRoomHandler(
		deps.RoomManager,
		deps.EventPublisher,
	)
	dispatcherImpl.Register(domainEvent.UserJoinedRoomEvent{}.EventType(), userJoinedHandler)

	userLeftHandler := handlers.NewUserLeftRoomHandler(
		deps.RoomManager,
		deps.EventPublisher,
	)
	dispatcherImpl.Register(domainEvent.UserLeftRoomEvent{}.EventType(), userLeftHandler)

	startButtonPressedHandler := handlers.NewStartButtonPressedHandler(
		deps.EventPublisher,
	)
	dispatcherImpl.Register(domainEvent.GameStartButtonPressedEvent{}.EventType(), startButtonPressedHandler)

	gameStartedHandler := handlers.NewGameStartedHandler(
		deps.EventPublisher,
	)
	dispatcherImpl.Register(domainEvent.GameStartedEvent{}.EventType(), gameStartedHandler)

	imageSendHandler := handlers.NewImageSendHandler(
		deps.EventPublisher,
	)
	dispatcherImpl.Register(domainEvent.ImageSendEvent{}.EventType(), imageSendHandler)

	clapTimeStartedHandler := handlers.NewClapTimeStartedHandler(
		deps.EventPublisher,
	)
	dispatcherImpl.Register(domainEvent.StartClapTimeEvent{}.EventType(), clapTimeStartedHandler)

	changeClapUserHandler := handlers.NewClapChangeUserHandler(
		deps.EventPublisher,
	)
	dispatcherImpl.Register(domainEvent.ChangeClapUserEvent{}.EventType(), changeClapUserHandler)

	clapSendHandler := handlers.NewClapSendHandler(
		deps.EventPublisher,
	)
	dispatcherImpl.Register(domainEvent.ClapSendEvent{}.EventType(), clapSendHandler)

	resultStartedHandler := handlers.NewResultStartedHandler(
		deps.EventPublisher,
	)
	dispatcherImpl.Register(domainEvent.StartResultPhaseEvent{}.EventType(), resultStartedHandler)

	roomSettingsUpdatedHandler := handlers.NewRoomSettingsUpdatedHandler(
		deps.EventPublisher,
	)
	dispatcherImpl.Register(domainEvent.RoomSettingsUpdatedEvent{}.EventType(), roomSettingsUpdatedHandler)

	return nil
}

func initLLM(deps *Dependencies) error {
	apiKey := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	if apiKey == "" {
		return errors.New("GEMINI_API_KEY is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ai.NewGeminiClient(ctx, ai.GeminiConfig{
		APIKey:       apiKey,
		DefaultModel: strings.TrimSpace(os.Getenv("GEMINI_MODEL")),
	})
	if err != nil {
		return fmt.Errorf("create gemini client: %w", err)
	}

	deps.LLMClient = client
	return nil
}

func initImage(deps *Dependencies) error {
	// ImageValidator の初期化 (最大5MBまで許可)
	const maxImageSize = 5 * 1024 * 1024 // 5MB
	deps.ImageValidator = validator.NewImageValidator(maxImageSize)

	// ImageStorage の初期化 (S3)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bucket := strings.TrimSpace(os.Getenv("S3_BUCKET"))
	if bucket == "" {
		return errors.New("S3_BUCKET is not set")
	}

	region := strings.TrimSpace(os.Getenv("AWS_REGION"))
	if region == "" {
		region = "ap-northeast-1" // デフォルトリージョン
	}

	imageStorage, err := storage.NewS3ImageStorage(ctx, bucket, region)
	if err != nil {
		return fmt.Errorf("create s3 image storage: %w", err)
	}

	deps.ImageStorage = imageStorage
	return nil
}

// UseCase関連の初期化
func initUseCases(deps *Dependencies) error {
	// Domain Services
	roomNumberGenerator := service.NewRoomNumberGenerator()
	clapCounter := clapinfra.NewMemoryClapCounter()

	// StartResultUseCaseを先に初期化（依存が少ない）
	deps.StartResultUseCase = battle.NewStartResultUseCase(
		deps.BattleRepository,
		deps.RoomRepository,
		deps.ImageRepository,
		clapCounter,
		deps.EventDispatcher,
	)

	// ClapTimeManageUseCaseを初期化（StartResultUseCaseに依存）
	deps.ClapTimeManageUseCase = clap.NewClapTimeManageUseCase(
		deps.BattleRepository,
		deps.EventDispatcher,
		deps.StartResultUseCase,
	)

	// StartClapTimeUseCaseを初期化（ClapTimeManageUseCaseに依存）
	deps.StartClapTimeUseCase = clap.NewStartClapTimeUseCase(
		deps.BattleRepository,
		deps.RoomRepository,
		deps.ImageRepository,
		deps.EventDispatcher,
		deps.LLMClient,
		deps.ClapTimeManageUseCase,
	)

	deps.ImageSubmissionScheduler = battleinfra.NewImageSubmissionScheduler(
		deps.StartClapTimeUseCase,
		time.Minute,
	)

	deps.ClapSendUseCase = clap.NewClapSendUseCase(
		deps.RoomRepository,
		deps.BattleRepository,
		clapCounter,
		deps.EventDispatcher,
	)

	// Room Usecases
	deps.CreateRoomUseCase = room.NewCreateRoomUseCase(
		deps.RoomRepository,
		deps.EventDispatcher,
		roomNumberGenerator,
	)

	deps.JoinRoomUseCase = room.NewJoinRoomUseCase(
		deps.RoomRepository,
		deps.EventDispatcher,
	)

	deps.CreateBattleUseCase = battle.NewCreateBattleUseCase(
		deps.BattleRepository,
		deps.BattleUserRepository,
		deps.RoomRepository,
		deps.LLMClient,
		deps.ImageSubmissionScheduler,
	)

	deps.GetBattleUseCase = battle.NewGetBattleUseCase(
		deps.BattleRepository,
	)

	deps.GetBattleIDUseCase = battle.NewGetBattleIDUseCase(
		deps.BattleRepository,
	)

	deps.LeaveRoomUseCase = room.NewLeaveRoomUseCase(
		deps.RoomRepository,
		deps.EventDispatcher,
	)

	deps.StartGameUseCase = room.NewStartGameUseCase(
		deps.RoomRepository,
		deps.EventDispatcher,
		deps.CreateBattleUseCase,
	)

	deps.GetRoomUseCase = room.NewGetRoomUseCase(
		deps.RoomRepository,
		deps.UserRepository,
	)

	deps.UpdateRoomSettingsUseCase = room.NewUpdateRoomSettingsUseCase(
		deps.RoomRepository,
		deps.EventDispatcher,
	)

	deps.GetImageUseCase = battle.NewGetImageUseCase(
		deps.BattleRepository,
		deps.ImageRepository,
	)

	deps.ImageSendUseCase = battle.NewImageSendUseCase(
		deps.BattleRepository,
		deps.ImageRepository,
		deps.ImageValidator,
		deps.ImageStorage,
		deps.EventDispatcher,
		deps.ImageSubmissionScheduler,
	)

	deps.GetResultUseCase = battle.NewGetResultUseCase(
		deps.BattleRepository,
		deps.RoomRepository,
		deps.ImageRepository,
	)

	deps.CreateUserUseCase = user.NewCreateUserUseCase(
		deps.UserRepository,
	)

	return nil
}

// データベース接続の初期化する
func initDatabase() (*sql.DB, error) {
	cfg, err := loadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("load database config: %w", err)
	}

	if cfg.driver != "mysql" {
		return nil, fmt.Errorf("unsupported DB_DRIVER %q (only mysql is supported)", cfg.driver)
	}

	db, err := sql.Open(cfg.driver, cfg.dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// コネクションプール設定
	if cfg.connMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.connMaxLifetime)
	}
	if cfg.connMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.connMaxIdleTime)
	}
	if cfg.maxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.maxOpenConns)
	}
	if cfg.maxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.maxIdleConns)
	}

	// 接続確認（リトライロジック付き）
	maxRetries := 5
	retryInterval := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := db.PingContext(ctx)
		cancel()

		if err == nil {
			// 接続成功
			fmt.Printf("Database connection established successfully")
			return db, nil
		}

		fmt.Printf("Database connection attempt %d/%d failed: %v", i+1, maxRetries, err)

		if i < maxRetries-1 {
			fmt.Printf("Retrying in %v...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	_ = db.Close()
	return nil, fmt.Errorf("failed to connect to database after %d attempts", maxRetries)
}

func initRepositories(deps *Dependencies) error {
	deps.RoomRepository = mysql.NewRoomRepository(deps.db)
	deps.BattleRepository = mysql.NewBattleRepository(deps.db)
	deps.BattleUserRepository = mysql.NewBattleUserRepository(deps.db)
	deps.ImageRepository = mysql.NewImageRepository(deps.db)
	deps.UserRepository = mysql.NewUserRepository(deps.db)

	return nil
}

func (d *Dependencies) Close() error {
	var err error

	if d.ImageSubmissionScheduler != nil {
		d.ImageSubmissionScheduler.Close()
	}

	if d.LLMClient != nil {
		if closeErr := d.LLMClient.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}

	if d.db != nil {
		if dbErr := d.db.Close(); dbErr != nil && err == nil {
			err = dbErr
		}
	}

	return err
}

func (d *Dependencies) DB() *sql.DB {
	return d.db
}

type databaseConfig struct {
	driver          string
	dsn             string
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
}

func loadDatabaseConfig() (databaseConfig, error) {
	cfg := databaseConfig{
		driver: envOrDefault("DB_DRIVER", "mysql"),
	}

	dsn := os.Getenv("DSN")
	if strings.TrimSpace(dsn) == "" {
		var err error
		switch cfg.driver {
		case "mysql":
			dsn, err = buildMySQLDSN()
			if err != nil {
				return cfg, err
			}
		default:
			return cfg, fmt.Errorf("DSN environment variable is required for driver %q", cfg.driver)
		}
	}
	cfg.dsn = strings.TrimSpace(dsn)

	var err error
	if cfg.maxOpenConns, err = parseIntEnv("DB_MAX_OPEN_CONNS", 32); err != nil {
		return cfg, fmt.Errorf("parse DB_MAX_OPEN_CONNS: %w", err)
	}

	if cfg.maxIdleConns, err = parseIntEnv("DB_MAX_IDLE_CONNS", 16); err != nil {
		return cfg, fmt.Errorf("parse DB_MAX_IDLE_CONNS: %w", err)
	}

	if cfg.connMaxLifetime, err = parseDurationSecondsEnv("DB_CONN_MAX_LIFETIME", 90*time.Minute); err != nil {
		return cfg, fmt.Errorf("parse DB_CONN_MAX_LIFETIME: %w", err)
	}

	if cfg.connMaxIdleTime, err = parseDurationSecondsEnv("DB_CONN_MAX_IDLE_TIME", 15*time.Minute); err != nil {
		return cfg, fmt.Errorf("parse DB_CONN_MAX_IDLE_TIME: %w", err)
	}

	return cfg, nil
}

func buildMySQLDSN() (string, error) {
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	if user == "" {
		return "", errors.New("DB_USER environment variable is required when DSN is not set")
	}

	host := envOrDefault("DB_HOST", "127.0.0.1")
	port := envOrDefault("DB_PORT", "3306")
	name := strings.TrimSpace(os.Getenv("DB_NAME"))
	if name == "" {
		return "", errors.New("DB_NAME environment variable is required when DSN is not set")
	}

	password := os.Getenv("DB_PASSWORD")
	params := strings.TrimSpace(os.Getenv("DB_PARAMS"))
	if params == "" {
		params = "charset=utf8mb4&parseTime=true&loc=Local"
	} else {
		params = strings.TrimPrefix(params, "?")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, name)
	if params != "" {
		dsn = dsn + "?" + params
	}

	return dsn, nil
}

func envOrDefault(key, defaultVal string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return defaultVal
}

func parseIntEnv(key string, defaultVal int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultVal, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer value %q", value)
	}
	return parsed, nil
}

func parseDurationSecondsEnv(key string, defaultVal time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultVal, nil
	}

	seconds, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid duration (seconds) value %q", value)
	}

	if seconds <= 0 {
		return 0, nil
	}

	return time.Duration(seconds) * time.Second, nil
}
