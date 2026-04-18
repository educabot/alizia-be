package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

// UpdateActivityRequest patches an activity template. Pointers opt in to a
// patch; SetDescription / SetDurationMinutes gate the null-clearing case that
// a plain pointer can't express.
type UpdateActivityRequest struct {
	OrgID              uuid.UUID
	ActivityID         int64
	Moment             *string
	Name               *string
	Description        *string
	SetDescription     bool
	DurationMinutes    *int
	SetDurationMinutes bool
}

func (r UpdateActivityRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.ActivityID <= 0 {
		return fmt.Errorf("%w: activity_id is required", providers.ErrValidation)
	}
	if r.Moment == nil && r.Name == nil && !r.SetDescription && !r.SetDurationMinutes {
		return fmt.Errorf("%w: at least one field must be provided", providers.ErrValidation)
	}
	if r.Name != nil && strings.TrimSpace(*r.Name) == "" {
		return fmt.Errorf("%w: name must not be blank", providers.ErrValidation)
	}
	if r.Moment != nil && !entities.ValidMoment(*r.Moment) {
		return fmt.Errorf("%w: moment must be apertura, desarrollo, or cierre", providers.ErrValidation)
	}
	if r.SetDurationMinutes && r.DurationMinutes != nil && *r.DurationMinutes <= 0 {
		return fmt.Errorf("%w: duration_minutes must be positive", providers.ErrValidation)
	}
	return nil
}

type UpdateActivity interface {
	Execute(ctx context.Context, req UpdateActivityRequest) (*entities.ActivityTemplate, error)
}

type updateActivityImpl struct {
	activities providers.ActivityTemplateProvider
}

func NewUpdateActivity(activities providers.ActivityTemplateProvider) UpdateActivity {
	return &updateActivityImpl{activities: activities}
}

func (uc *updateActivityImpl) Execute(ctx context.Context, req UpdateActivityRequest) (*entities.ActivityTemplate, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	current, err := uc.activities.GetActivity(ctx, req.OrgID, req.ActivityID)
	if err != nil {
		return nil, err
	}

	if req.Moment != nil {
		current.Moment = entities.ClassMoment(*req.Moment)
	}
	if req.Name != nil {
		current.Name = strings.TrimSpace(*req.Name)
	}
	if req.SetDescription {
		current.Description = req.Description
	}
	if req.SetDurationMinutes {
		current.DurationMinutes = req.DurationMinutes
	}

	if err := uc.activities.UpdateActivity(ctx, current); err != nil {
		return nil, err
	}
	return current, nil
}
