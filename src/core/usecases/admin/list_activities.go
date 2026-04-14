package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type ListActivitiesRequest struct {
	OrgID  uuid.UUID
	Moment *string // optional filter
}

func (r ListActivitiesRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.Moment != nil && !entities.ValidMoment(*r.Moment) {
		return fmt.Errorf("%w: moment must be apertura, desarrollo, or cierre", providers.ErrValidation)
	}
	return nil
}

type ListActivities interface {
	Execute(ctx context.Context, req ListActivitiesRequest) ([]entities.ActivityTemplate, error)
}

type listActivitiesImpl struct {
	activities providers.ActivityTemplateProvider
}

func NewListActivities(activities providers.ActivityTemplateProvider) ListActivities {
	return &listActivitiesImpl{activities: activities}
}

func (uc *listActivitiesImpl) Execute(ctx context.Context, req ListActivitiesRequest) ([]entities.ActivityTemplate, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var moment *entities.ClassMoment
	if req.Moment != nil {
		m := entities.ClassMoment(*req.Moment)
		moment = &m
	}

	return uc.activities.ListActivities(ctx, req.OrgID, moment)
}
