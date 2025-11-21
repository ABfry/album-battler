// room_repository.go は rooms テーブルを扱う MySQL 実装を提供する。
// 責務: 永続層の行と entity.Room の相互変換、及び単純な検索/保存を集約する。
// 依存: database/sql（接続）、dao.go 内の共通クエリビルダー。
// 使用例: usecase 層が repository.RoomRepository を通じて部屋情報を取得・保存する。
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/google/uuid"
)

var _ repository.RoomRepository = (*mysqlRoomRepository)(nil)

// mysqlRoomRepository は RoomRepository を MySQL 上で実装する具体型。
type mysqlRoomRepository struct {
	db *sql.DB
}

// NewRoomRepository は MySQL 接続を受け取り、RoomRepository 実装を返す。
// why: インフラ層の具象型を外部へ漏らさず、テストで差し替えられるようにするため。
func NewRoomRepository(db *sql.DB) repository.RoomRepository {
	return &mysqlRoomRepository{db: db}
}

// FindByID は主キー検索で部屋情報を 1 件取得する。
// why: ルーム詳細画面などで ID 直接参照が必要なための共通処理。
func (r *mysqlRoomRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Room, error) {
	row, err := findRowByKey(ctx, r.db, RoomsTable, "id", id.String())
	if err != nil {
		return nil, err
	}
	room, err := r.createRoom(row)
	if err != nil {
		return nil, err
	}

	// 部屋が見つからなかった場合はnilを返す
	if room == nil {
		return nil, nil
	}

	// room_users から UserIDs を読み込む
	if err := r.loadUserIDs(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

// FindByRoomNumber は表示用のルーム番号をキーに 1 件取得する。
// why: ユーザーが数字で参加する UX のため、room_number をユニークキーとして扱っている。
func (r *mysqlRoomRepository) FindByRoomNumber(ctx context.Context, roomNumber int) (*entity.Room, error) {
	row, err := findRowByKey(ctx, r.db, RoomsTable, "room_number", roomNumber)
	if err != nil {
		return nil, err
	}
	room, err := r.createRoom(row)
	if err != nil {
		return nil, err
	}

	// 部屋が見つからなかった場合はnilを返す
	if room == nil {
		return nil, nil
	}

	// room_users から UserIDs を読み込む
	if err := r.loadUserIDs(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

func (r *mysqlRoomRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Room, error) {
	query := "SELECT room_id FROM room_users WHERE user_id = ?"
	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("failed to close rows: %v\n", err)
		}
	}()

	var rooms []*entity.Room
	for rows.Next() {
		var roomIDStr string
		if err := rows.Scan(&roomIDStr); err != nil {
			return nil, err
		}

		roomRow, err := findRowByKey(ctx, r.db, RoomsTable, "id", roomIDStr)
		if err != nil {
			return nil, err
		}

		room, err := r.createRoom(roomRow)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}

// FindAll は rooms テーブル全件を読み出す。
// why: 管理画面やバッチでの一括処理向け。件数増大時は呼び出し側でページネーションを検討する。
func (r *mysqlRoomRepository) FindAll(ctx context.Context) ([]*entity.Room, error) {
	rows, err := findAllRows(ctx, r.db, RoomsTable)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("failed to close rows: %v\n", err)
		}
	}()

	var rooms []*entity.Room
	for rows.Next() {
		room, err := r.createRoom(rows)
		if err != nil {
			return nil, err
		}
		// room_users から UserIDs を読み込む
		if err := r.loadUserIDs(ctx, room); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}

// Save は rooms テーブルへ upsert を行う。
// why: INSERT/UPDATE の両方で同じカラムセットを維持するため。
func (r *mysqlRoomRepository) Save(ctx context.Context, room *entity.Room) error {
	statusStr, err := toDBRoomStatus(room.Status)
	if err != nil {
		return err
	}

	var hostUserID interface{}
	if room.HostUserID != nil {
		hostUserID = room.HostUserID.String()
	} else {
		hostUserID = nil // 明示的にnilを設定
	}

	// rooms テーブルを保存
	if err := save(
		ctx,
		r.db,
		RoomsTable,
		room.ID.String(),
		room.RoomNumber,
		hostUserID,
		room.CreatedAt,
		room.ExpiredAt,
		statusStr,
		room.MaxUsers,
	); err != nil {
		return err
	}

	// room_users テーブルを更新（全削除→再挿入）
	return r.saveUserIDs(ctx, room)
}

// --- private helpers ---

// loadUserIDs は room_users テーブルから UserIDs を読み込んで room に設定する。
// why: 部屋に所属するユーザー一覧を取得し、ドメインエンティティに反映するため。
func (r *mysqlRoomRepository) loadUserIDs(ctx context.Context, room *entity.Room) error {
	query := "SELECT user_id FROM room_users WHERE room_id = ? ORDER BY joined_at"
	rows, err := r.db.QueryContext(ctx, query, room.ID.String())
	if err != nil {
		return fmt.Errorf("failed to load user IDs: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			fmt.Printf("failed to close rows: %v\n", closeErr)
		}
	}()

	var userIDs []uuid.UUID
	for rows.Next() {
		var userIDStr string
		if err := rows.Scan(&userIDStr); err != nil {
			return fmt.Errorf("failed to scan user ID: %w", err)
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return fmt.Errorf("failed to parse user ID: %w", err)
		}
		userIDs = append(userIDs, userID)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows iteration error: %w", err)
	}

	room.UserIDs = userIDs
	return nil
}

// saveUserIDs は room_users テーブルに UserIDs を保存する。
// why: 既存のユーザー関連付けを全削除してから再挿入することで、整合性を保つ。
func (r *mysqlRoomRepository) saveUserIDs(ctx context.Context, room *entity.Room) error {
	// トランザクション開始
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			fmt.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	// 既存のroom_usersを削除
	deleteQuery := "DELETE FROM room_users WHERE room_id = ?"
	if _, err := tx.ExecContext(ctx, deleteQuery, room.ID.String()); err != nil {
		return fmt.Errorf("failed to delete room_users: %w", err)
	}

	// UserIDsを再挿入
	if len(room.UserIDs) > 0 {
		insertQuery := "INSERT INTO room_users (room_id, user_id, joined_at) VALUES (?, ?, NOW())"
		for _, userID := range room.UserIDs {
			if _, err := tx.ExecContext(ctx, insertQuery, room.ID.String(), userID.String()); err != nil {
				return fmt.Errorf("failed to insert room_user: %w", err)
			}
		}
	}

	// コミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// createRoom は 1 行分のスキャン結果から entity.Room を生成する。
// why: カラム順序と変換ロジックを集中管理し、スキーマ変更時の漏れを防ぐ。
func (r *mysqlRoomRepository) createRoom(scanner rowScanner) (*entity.Room, error) {
	var (
		idStr        string
		roomNumber   int
		hostUserID   string
		createdAt    time.Time
		expiredAt    sql.NullTime
		statusString string
		maxUsers     int
	)

	if err := scanner.Scan(&idStr, &roomNumber, &hostUserID, &createdAt, &expiredAt, &statusString, &maxUsers); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	roomID, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse room id: %w", err)
	}
	hostID, err := uuid.Parse(hostUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse host user id: %w", err)
	}
	status, err := toEntityRoomStatus(statusString)
	if err != nil {
		return nil, err
	}

	room := &entity.Room{
		ID:         roomID,
		RoomNumber: roomNumber,
		CreatedAt:  createdAt,
		ExpiredAt:  expiredAt.Time,
		Status:     status,
		MaxUsers:   maxUsers,
	}
	room.HostUserID = &hostID
	return room, nil
}

// toEntityRoomStatus は DB 上の文字列をドメイン列挙型へ変換する。
func toEntityRoomStatus(dbValue string) (entity.RoomStatus, error) {
	if status, ok := dbToEntityRoomStatus[dbValue]; ok {
		return status, nil
	}
	return 0, fmt.Errorf("unknown room status string: %s", dbValue)
}

// toDBRoomStatus はドメイン列挙を DB の文字列表現へ変換する。
func toDBRoomStatus(status entity.RoomStatus) (string, error) {
	if dbValue, ok := entityToDBRoomStatus[status]; ok {
		return dbValue, nil
	}
	return "", fmt.Errorf("unknown room status enum: %d", status)
}

var (
	dbToEntityRoomStatus = map[string]entity.RoomStatus{
		"waiting":   entity.WaitJoin,
		"full":      entity.FullyJoined,
		"battling":  entity.InBattle,
		"clap_time": entity.ClapTime,
		"result":    entity.Result,
		"closed":    entity.Closed,
	}
	entityToDBRoomStatus = map[entity.RoomStatus]string{
		entity.WaitJoin:    "waiting",
		entity.FullyJoined: "full",
		entity.InBattle:    "battling",
		entity.ClapTime:    "clap_time",
		entity.Result:      "result",
		entity.Closed:      "closed",
	}
)
