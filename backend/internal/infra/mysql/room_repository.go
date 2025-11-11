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
	return r.createRoom(row)
}

// FindByRoomNumber は表示用のルーム番号をキーに 1 件取得する。
// why: ユーザーが数字で参加する UX のため、room_number をユニークキーとして扱っている。
func (r *mysqlRoomRepository) FindByRoomNumber(ctx context.Context, roomNumber int) (*entity.Room, error) {
	row, err := findRowByKey(ctx, r.db, RoomsTable, "room_number", roomNumber)
	if err != nil {
		return nil, err
	}
	return r.createRoom(row)
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

	return save(
		ctx,
		r.db,
		RoomsTable,
		room.ID.String(),
		room.RoomNumber,
		hostUserID,
		room.CreatedAt,
		room.ExpiredAt,
		statusStr,
	)
}

// --- private helpers ---

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
	)

	if err := scanner.Scan(&idStr, &roomNumber, &hostUserID, &createdAt, &expiredAt, &statusString); err != nil {
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
		"waiting":  entity.WaitJoin,
		"full":     entity.FullyJoined,
		"battling": entity.InBattle,
		"result":   entity.Result,
		"closed":   entity.Closed,
	}
	entityToDBRoomStatus = map[entity.RoomStatus]string{
		entity.WaitJoin:    "waiting",
		entity.FullyJoined: "full",
		entity.InBattle:    "battling",
		entity.Result:      "result",
		entity.Closed:      "closed",
	}
)
