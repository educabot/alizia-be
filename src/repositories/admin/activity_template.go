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

func (r *activityTemplateRepo) ListActivities(ctx context.Context, orgID uuid.UUID, moment *entities.ClassMoment) ([]entities.ActivityTemplate, error) {
	query := r.db.WithContext(ctx).Where("organization_id = ?", orgID)
	if moment != nil {
		query = query.Where("moment = ?", *moment)
	}
	var activities []entities.ActivityTemplate
	err := query.Order("moment, name").Limit(100).Find(&activities).Error
	return activities, err
}
