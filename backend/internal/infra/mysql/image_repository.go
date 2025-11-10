package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/ABfry/album-battler/backend/internal/domain/repository"
	"github.com/google/uuid"
)

var _ repository.ImageRepository = (*mysqlImageRepository)(nil)

type mysqlImageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) repository.ImageRepository {
	return &mysqlImageRepository{db: db}
}

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

func (r *mysqlImageRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Image, error) {
	return r.findBy(ctx, "id", id.String())
}

func (r *mysqlImageRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.Image, error) {
	return r.findBy(ctx, "user_id", userID.String())
}

func (r *mysqlImageRepository) FindByBattleID(ctx context.Context, battleID uuid.UUID) (*entity.Image, error) {
	return r.findBy(ctx, "battle_id", battleID.String())
}

// --- private ---

func (r *mysqlImageRepository) findBy(ctx context.Context, key string, v any) (*entity.Image, error) {
	row, err := findByKey(ctx, r.db, ImagesTable, key, v)
	if err != nil {
		return nil, err
	}

	var (
		idStr, userIDStr, battleIDStr string
		imageURL                      string
		uploadedAt                    time.Time
		aiScore                       sql.NullFloat64
		userScore                     sql.NullInt64
	)

	if err := row.Scan(&idStr, &userIDStr, &battleIDStr, &imageURL, &uploadedAt, &aiScore, &userScore); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
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
