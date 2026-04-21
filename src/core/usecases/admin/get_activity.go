package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type GetActivityRequest struct {
	OrgID      uuid.UUID
	ActivityID int64
}

func (r GetActivityRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.ActivityID <= 0 {
		return fmt.Errorf("%w: activity_id is required", providers.ErrValidation)
	}
	return nil
}

type GetActivity interface {
	Execute(ctx context.Context, req GetActivityRequest) (*entities.ActivityTemplate, error)
}

type getActivityImpl struct {
	activities providers.ActivityTemplateProvider
}

func NewGetActivity(activities providers.ActivityTemplateProvider) GetActivity {
	return &getActivityImpl{activities: activities}
}

func (uc *getActivityImpl) Execute(ctx context.Context, req GetActivityRequest) (*entities.ActivityTemplate, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return uc.activities.GetActivity(ctx, req.OrgID, req.ActivityID)
}
