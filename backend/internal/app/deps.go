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

	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/ABfry/album-battler/backend/internal/domain/service"
	"github.com/ABfry/album-battler/backend/internal/infra/event"
	"github.com/ABfry/album-battler/backend/internal/infra/event/handlers"
	"github.com/ABfry/album-battler/backend/internal/infra/mysql"
	"github.com/ABfry/album-battler/backend/internal/infra/websocket"
	"github.com/ABfry/album-battler/backend/internal/usecase/room"
)

// -- 依存関係の定義 --

type Dependencies struct {
	db *sql.DB

	// Repository
	RoomRepository   repository.RoomRepository
	BattleRepository repository.BattleRepository
	ImageRepository  repository.ImageRepository
	UserRepository   repository.UserRepository

	// WebSocket関連
	WebSocketHub   *websocket.Hub
	EventPublisher service.EventPublisher
	RoomManager    service.RoomManager

	// Event関連
	EventDispatcher service.EventDispatcher

	// Usecase
	CreateRoomUseCase *room.CreateRoomUseCase
	JoinRoomUseCase   *room.JoinRoomUseCase
	LeaveRoomUseCase  *room.LeaveRoomUseCase
	StartGameUseCase  *room.StartGameUseCase
	GetRoomUseCase    *room.GetRoomUseCase
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

	// Usecase関連の初期化
	if err := initUseCases(deps); err != nil {
		_ = db.Close()
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
	dispatcherImpl.Register("user_joined_room", userJoinedHandler)

	userLeftHandler := handlers.NewUserLeftRoomHandler(
		deps.RoomManager,
		deps.EventPublisher,
	)
	dispatcherImpl.Register("user_left_room", userLeftHandler)

	gameStartedHandler := handlers.NewGameStartedHandler(
		deps.EventPublisher,
	)
	dispatcherImpl.Register("game_started", gameStartedHandler)

	return nil
}

// UseCase関連の初期化
func initUseCases(deps *Dependencies) error {
	// Room Usecases
	deps.CreateRoomUseCase = room.NewCreateRoomUseCase(
		deps.RoomRepository,
		deps.EventDispatcher,
	)

	deps.JoinRoomUseCase = room.NewJoinRoomUseCase(
		deps.RoomRepository,
		deps.EventDispatcher,
	)

	deps.LeaveRoomUseCase = room.NewLeaveRoomUseCase(
		deps.RoomRepository,
		deps.EventDispatcher,
	)

	deps.StartGameUseCase = room.NewStartGameUseCase(
		deps.RoomRepository,
		deps.EventDispatcher,
	)

	deps.GetRoomUseCase = room.NewGetRoomUseCase(
		deps.RoomRepository,
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
	deps.ImageRepository = mysql.NewImageRepository(deps.db)
	deps.UserRepository = mysql.NewUserRepository(deps.db)

	return nil
}

func (d *Dependencies) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
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
