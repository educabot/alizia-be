package admin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

const defaultDesarrolloMaxActivities = 0 // 0 = unlimited

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
	orgs       providers.OrganizationProvider
	activities providers.ActivityTemplateProvider
}

func NewCreateActivity(orgs providers.OrganizationProvider, activities providers.ActivityTemplateProvider) CreateActivity {
	return &createActivityImpl{orgs: orgs, activities: activities}
}

func (uc *createActivityImpl) Execute(ctx context.Context, req CreateActivityRequest) (*entities.ActivityTemplate, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	moment := entities.ClassMoment(req.Moment)

	if moment == entities.MomentDesarrollo {
		org, err := uc.orgs.FindByID(ctx, req.OrgID)
		if err != nil {
			return nil, err
		}
		if maxAllowed := desarrolloMaxActivities(org); maxAllowed > 0 {
			count, err := uc.activities.CountByMoment(ctx, req.OrgID, moment)
			if err != nil {
				return nil, err
			}
			if count >= int64(maxAllowed) {
				return nil, fmt.Errorf("%w: desarrollo activities limit %d reached", providers.ErrActivityMaxReached, maxAllowed)
			}
		}
	}

	activity := &entities.ActivityTemplate{
		OrganizationID:  req.OrgID,
		Moment:          moment,
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

// desarrolloMaxActivities reads `desarrollo_max_activities` from org config.
// A value <= 0 means unlimited.
func desarrolloMaxActivities(org *entities.Organization) int {
	var cfg map[string]any
	if err := json.Unmarshal(org.Config, &cfg); err != nil {
		return defaultDesarrolloMaxActivities
	}
	if v, ok := cfg["desarrollo_max_activities"]; ok {
		if f, ok := v.(float64); ok && f > 0 {
			return int(f)
		}
	}
	return defaultDesarrolloMaxActivities
}
