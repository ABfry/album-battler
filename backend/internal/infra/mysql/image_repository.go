// image_repository.go は MySQL 上の images テーブルを操作するリポジトリ実装を提供する。
// 責務: images テーブルの CRUD（現在は保存と各種検索のみ）を集約し、ドメイン層に対して
// 一貫したエンティティ変換を行う。依存: database/sql による接続オブジェクトと DAO ヘルパー。
// 使用例: usecase 層から ImageRepository を注入して呼び出す。
package mysql

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/google/uuid"
)

var _ repository.ImageRepository = (*mysqlImageRepository)(nil)

type mysqlImageRepository struct {
	db *sql.DB
}

// NewImageRepository は MySQL 接続を受け取り ImageRepository 実装を返す。
// why: 他レイヤーから具体型を直接扱わせないことでテスト容易性を確保する。
func NewImageRepository(db *sql.DB) repository.ImageRepository {
	return &mysqlImageRepository{db: db}
}

// Save は images テーブルへ upsert を行う。
// why: save ヘルパーを使うことでカラム順序の重複定義を防ぎ、将来の変更コストを下げる。
func (r *mysqlImageRepository) Save(ctx context.Context, image *entity.Image) error {
	return save(
		ctx, r.db, ImagesTable,
		image.ID.String(),
		image.UserID.String(),
		image.BattleID.String(),
		image.ImageURL,
		image.UploadedAt,
		image.AIScore,
		image.UserScore,
	)
}

// FindByID は主キー検索で 1 件だけ取り出す。
// why: エンティティ単位の検証で厳密に 1 件を想定するケースが多いため。
func (r *mysqlImageRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Image, error) {
	return r.findBy(ctx, "id", id.String())
}

// FindImagesByUserID はユーザーが投稿した画像をすべて返す。
// why: バトル画面でユーザー単位の履歴を一覧表示するユースケースがあるため。
func (r *mysqlImageRepository) FindImagesByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Image, error) {
	return r.findImagesBy(ctx, "user_id", userID.String())
}

func (r *mysqlImageRepository) FindImagesByBattleAndUserID(ctx context.Context, battleID, userID uuid.UUID) (*entity.Image, error) {
	const query = `SELECT id, user_id, battle_id, image_url, uploaded_at, ai_score, user_score FROM images WHERE battle_id = ? AND user_id = ?`

	rows, err := r.db.QueryContext(ctx, query, battleID.String(), userID.String())
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}()

	if rows.Next() {
		return r.createImage(rows)
	}
	return nil, sql.ErrNoRows
}

// FindImagesByBattleID は指定バトルに紐づく画像一覧を返す。
// why: バトル集計時に一括で読み込む必要があるため。
func (r *mysqlImageRepository) FindImagesByBattleID(ctx context.Context, battleID uuid.UUID) ([]*entity.Image, error) {
	return r.findImagesBy(ctx, "battle_id", battleID.String())
}

func (r *mysqlImageRepository) CountDistinctUsersByBattleID(ctx context.Context, battleID uuid.UUID) (int, error) {
	const query = `SELECT COUNT(DISTINCT user_id) FROM images WHERE battle_id = ?`
	row := r.db.QueryRowContext(ctx, query, battleID.String())

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// --- private ---

// findBy は単一レコード取得専用のヘルパー。
// why: WHERE 句のカラムだけを差し替えたいパターンが多いため。
func (r *mysqlImageRepository) findBy(ctx context.Context, key string, v any) (*entity.Image, error) {
	row, err := findRowByKey(ctx, r.db, ImagesTable, key, v)
	if err != nil {
		return nil, err
	}

	return r.createImage(row)
}

// findImagesBy は複数レコード取得を共通化するヘルパー。
// createImage を使い回し、カラムスキャンの重複を避ける。
func (r *mysqlImageRepository) findImagesBy(ctx context.Context, key string, v any) ([]*entity.Image, error) {
	rows, err := findRowsByKey(ctx, r.db, ImagesTable, key, v)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}()

	var images []*entity.Image
	for rows.Next() {
		img, err := r.createImage(rows)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return images, nil
}

// createImage は rows/row から entity.Image を生成する唯一の場所。
// why: カラム順を 1 箇所に閉じ込めてスキーマ変更時の齟齬を防ぐ。
func (r *mysqlImageRepository) createImage(scanner rowScanner) (*entity.Image, error) {
	var (
		idStr, userIDStr, battleIDStr string
		imageURL                      string
		uploadedAt                    time.Time
		aiScore                       sql.NullFloat64
		userScore                     sql.NullInt64
	)

	if err := scanner.Scan(&idStr, &userIDStr, &battleIDStr, &imageURL, &uploadedAt, &aiScore, &userScore); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}
	bid, err := uuid.Parse(battleIDStr)
	if err != nil {
		return nil, err
	}

	img := &entity.Image{
		ID:         id,
		UserID:     uid,
		BattleID:   bid,
		ImageURL:   imageURL,
		UploadedAt: uploadedAt,
	}
	if aiScore.Valid {
		img.AIScore = aiScore.Float64
	}
	if userScore.Valid {
		img.UserScore = int(userScore.Int64)
	}
	return img, nil
}
