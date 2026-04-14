package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type CreateActivityRequest struct {
	OrgID           uuid.UUID
	Moment          string
	Name            string
	Description     *string
	DurationMinutes *int
}

func (r CreateActivityRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.Name == "" {
		return fmt.Errorf("%w: name is required", providers.ErrValidation)
	}
	if !entities.ValidMoment(r.Moment) {
		return fmt.Errorf("%w: moment must be apertura, desarrollo, or cierre", providers.ErrValidation)
	}
	return nil
}

type CreateActivity interface {
	Execute(ctx context.Context, req CreateActivityRequest) (*entities.ActivityTemplate, error)
}

type createActivityImpl struct {
	activities providers.ActivityTemplateProvider
}

func NewCreateActivity(activities providers.ActivityTemplateProvider) CreateActivity {
	return &createActivityImpl{activities: activities}
}

func (uc *createActivityImpl) Execute(ctx context.Context, req CreateActivityRequest) (*entities.ActivityTemplate, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	activity := &entities.ActivityTemplate{
		OrganizationID:  req.OrgID,
		Moment:          entities.ClassMoment(req.Moment),
		Name:            req.Name,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
	}

	id, err := uc.activities.CreateActivity(ctx, activity)
	if err != nil {
		return nil, err
	}

	activity.ID = id
	return activity, nil
}
