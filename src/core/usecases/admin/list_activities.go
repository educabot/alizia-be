package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type ListActivitiesRequest struct {
	OrgID      uuid.UUID
	Moment     *string // optional filter
	Pagination providers.Pagination
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

type ListActivitiesResponse struct {
	Items []entities.ActivityTemplate
	More  bool
}

type ListActivities interface {
	Execute(ctx context.Context, req ListActivitiesRequest) (ListActivitiesResponse, error)
}

type listActivitiesImpl struct {
	activities providers.ActivityTemplateProvider
}

func NewListActivities(activities providers.ActivityTemplateProvider) ListActivities {
	return &listActivitiesImpl{activities: activities}
}

func (uc *listActivitiesImpl) Execute(ctx context.Context, req ListActivitiesRequest) (ListActivitiesResponse, error) {
	if err := req.Validate(); err != nil {
		return ListActivitiesResponse{}, err
	}

	var moment *entities.ClassMoment
	if req.Moment != nil {
		m := entities.ClassMoment(*req.Moment)
		moment = &m
	}

	items, more, err := uc.activities.ListActivities(ctx, req.OrgID, moment, req.Pagination)
	if err != nil {
		return ListActivitiesResponse{}, err
	}
	return ListActivitiesResponse{Items: items, More: more}, nil
}
