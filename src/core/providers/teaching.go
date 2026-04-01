package providers

import (
	"context"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/google/uuid"
)

type TeachingProvider interface {
	CreateLessonPlan(ctx context.Context, plan *entities.LessonPlan) (int64, error)
	GetLessonPlan(ctx context.Context, orgID uuid.UUID, planID int64) (*entities.LessonPlan, error)
	UpdateLessonPlan(ctx context.Context, plan *entities.LessonPlan) error
	ListLessonPlans(ctx context.Context, orgID uuid.UUID, coordDocClassID int64) ([]entities.LessonPlan, error)

	CreateActivity(ctx context.Context, activity *entities.Activity) (int64, error)
	ListActivities(ctx context.Context, planID int64) ([]entities.Activity, error)
	UpdateActivity(ctx context.Context, activity *entities.Activity) error
	DeleteActivity(ctx context.Context, activityID int64) error
}
