package admin

import (
	"context"

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
