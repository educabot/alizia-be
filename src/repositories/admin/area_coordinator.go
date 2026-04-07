package admin

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type areaCoordinatorRepo struct {
	db *gorm.DB
}

func NewAreaCoordinatorRepo(db *gorm.DB) providers.AreaCoordinatorProvider {
	return &areaCoordinatorRepo{db: db}
}

func (r *areaCoordinatorRepo) Assign(ctx context.Context, areaID, userID int64) (*entities.AreaCoordinator, error) {
	ac := &entities.AreaCoordinator{
		AreaID: areaID,
		UserID: userID,
	}
	if err := r.db.WithContext(ctx).Create(ac).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fmt.Errorf("%w: coordinator already assigned to area", providers.ErrConflict)
		}
		return nil, err
	}
	return ac, nil
}

func (r *areaCoordinatorRepo) Remove(ctx context.Context, areaID, userID int64) error {
	result := r.db.WithContext(ctx).
		Where("area_id = ? AND user_id = ?", areaID, userID).
		Delete(&entities.AreaCoordinator{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: assignment not found", providers.ErrNotFound)
	}
	return nil
}

func (r *areaCoordinatorRepo) FindByAreaID(ctx context.Context, areaID int64) ([]entities.AreaCoordinator, error) {
	var results []entities.AreaCoordinator
	err := r.db.WithContext(ctx).Where("area_id = ?", areaID).Find(&results).Error
	return results, err
}

func (r *areaCoordinatorRepo) FindByUserID(ctx context.Context, userID int64) ([]entities.AreaCoordinator, error) {
	var results []entities.AreaCoordinator
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&results).Error
	return results, err
}

func (r *areaCoordinatorRepo) IsCoordinator(ctx context.Context, areaID, userID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.AreaCoordinator{}).
		Where("area_id = ? AND user_id = ?", areaID, userID).
		Count(&count).Error
	return count > 0, err
}
