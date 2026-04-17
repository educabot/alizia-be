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

func (r *areaRepo) UpdateArea(ctx context.Context, area *entities.Area) error {
	err := r.db.WithContext(ctx).
		Model(&entities.Area{}).
		Where("organization_id = ? AND id = ?", area.OrganizationID, area.ID).
		Updates(map[string]any{
			"name":        area.Name,
			"description": area.Description,
		}).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("%w: area name already exists in this organization", providers.ErrConflict)
		}
		return err
	}
	return nil
}

// CountDependencies counts subjects that reference the area. Course-subjects
// are not counted separately because they reference subjects (not areas
// directly); a non-zero Subjects count already implies any downstream
// course-subjects would also block deletion. Coordination documents are not
// counted yet: the table ships in Épica 4.
func (r *areaRepo) CountDependencies(ctx context.Context, orgID uuid.UUID, id int64) (providers.AreaDependencies, error) {
	var deps providers.AreaDependencies
	if err := r.db.WithContext(ctx).
		Model(&entities.Subject{}).
		Where("organization_id = ? AND area_id = ?", orgID, id).
		Count(&deps.Subjects).Error; err != nil {
		return deps, err
	}
	return deps, nil
}

// DeleteArea removes the area and its coordinator assignments in a single
// transaction. Caller must verify with CountDependencies that no blocking
// dependencies exist. Coordinator rows are role assignments (no domain data)
// so they're cascade-removed here rather than forcing admins to remove them
// individually via DELETE /areas/:id/coordinators/:user_id.
func (r *areaRepo) DeleteArea(ctx context.Context, orgID uuid.UUID, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("area_id = ?", id).Delete(&entities.AreaCoordinator{}).Error; err != nil {
			return err
		}
		result := tx.Where("organization_id = ? AND id = ?", orgID, id).Delete(&entities.Area{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("%w: area %d", providers.ErrNotFound, id)
		}
		return nil
	})
}
