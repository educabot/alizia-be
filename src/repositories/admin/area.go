package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type areaRepo struct {
	db *gorm.DB
}

func NewAreaRepo(db *gorm.DB) providers.AreaProvider {
	return &areaRepo{db: db}
}

func (r *areaRepo) CreateArea(ctx context.Context, area *entities.Area) (int64, error) {
	if err := r.db.WithContext(ctx).Create(area).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%w: area name already exists in this organization", providers.ErrConflict)
		}
		return 0, err
	}
	return area.ID, nil
}

func (r *areaRepo) GetArea(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Area, error) {
	var area entities.Area
	err := r.db.WithContext(ctx).
		Preload("Subjects").
		Preload("Coordinators.User").
		Where("organization_id = ? AND id = ?", orgID, id).
		First(&area).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: area %d", providers.ErrNotFound, id)
		}
		return nil, err
	}
	return &area, nil
}

func (r *areaRepo) ListAreas(ctx context.Context, orgID uuid.UUID) ([]entities.Area, error) {
	var areas []entities.Area
	err := r.db.WithContext(ctx).
		Preload("Subjects").
		Preload("Coordinators.User").
		Where("organization_id = ?", orgID).
		Order("name ASC").Limit(100).
		Find(&areas).Error
	return areas, err
}
