package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type activityTemplateRepo struct {
	db *gorm.DB
}

func NewActivityTemplateRepo(db *gorm.DB) providers.ActivityTemplateProvider {
	return &activityTemplateRepo{db: db}
}

func (r *activityTemplateRepo) CreateActivity(ctx context.Context, activity *entities.ActivityTemplate) (int64, error) {
	if err := r.db.WithContext(ctx).Create(activity).Error; err != nil {
		return 0, err
	}
	return activity.ID, nil
}

func (r *activityTemplateRepo) ListActivities(ctx context.Context, orgID uuid.UUID, moment *entities.ClassMoment, p providers.Pagination) ([]entities.ActivityTemplate, bool, error) {
	p = p.Normalize()
	query := r.db.WithContext(ctx).Where("organization_id = ?", orgID)
	if moment != nil {
		query = query.Where("moment = ?", *moment)
	}
	var rows []entities.ActivityTemplate
	err := query.Order("moment, name").
		Offset(p.Offset).
		Limit(p.Limit + 1).
		Find(&rows).Error
	if err != nil {
		return nil, false, err
	}
	if len(rows) > p.Limit {
		return rows[:p.Limit], true, nil
	}
	return rows, false, nil
}

func (r *activityTemplateRepo) CountByMoment(ctx context.Context, orgID uuid.UUID, moment entities.ClassMoment) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.ActivityTemplate{}).
		Where("organization_id = ? AND moment = ?", orgID, moment).
		Count(&count).Error
	return count, err
}

func (r *activityTemplateRepo) GetActivity(ctx context.Context, orgID uuid.UUID, id int64) (*entities.ActivityTemplate, error) {
	var activity entities.ActivityTemplate
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", orgID, id).
		First(&activity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: activity %d", providers.ErrNotFound, id)
		}
		return nil, err
	}
	return &activity, nil
}

// UpdateActivity writes the activity's mutable columns scoped to (org, id).
// Returns ErrNotFound if the row doesn't belong to the org.
func (r *activityTemplateRepo) UpdateActivity(ctx context.Context, activity *entities.ActivityTemplate) error {
	result := r.db.WithContext(ctx).
		Model(&entities.ActivityTemplate{}).
		Where("organization_id = ? AND id = ?", activity.OrganizationID, activity.ID).
		Updates(map[string]any{
			"moment":           activity.Moment,
			"name":             activity.Name,
			"description":      activity.Description,
			"duration_minutes": activity.DurationMinutes,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: activity %d", providers.ErrNotFound, activity.ID)
	}
	return nil
}

// DeleteActivity removes an activity scoped to (org, id). No dependency check
// today — lesson_plans and coordination_documents tables ship in later epics.
func (r *activityTemplateRepo) DeleteActivity(ctx context.Context, orgID uuid.UUID, id int64) error {
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", orgID, id).
		Delete(&entities.ActivityTemplate{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: activity %d", providers.ErrNotFound, id)
	}
	return nil
}
