package onboarding

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type TourStep struct {
	Key         string `json:"key"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

var defaultTourSteps = []TourStep{
	{Key: "welcome", Title: "Welcome to Alizia", Description: "Alizia helps you plan the school year collaboratively.", Order: 1},
	{Key: "explore", Title: "Explore the platform", Description: "Browse the sections to discover available tools.", Order: 2},
}

type GetTourStepsRequest struct {
	OrgID  uuid.UUID
	UserID int64
}

func (r GetTourStepsRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.UserID == 0 {
		return fmt.Errorf("%w: user_id is required", providers.ErrValidation)
	}
	return nil
}

type GetTourSteps interface {
	Execute(ctx context.Context, req GetTourStepsRequest) ([]TourStep, error)
}

type getTourStepsImpl struct {
	orgs  providers.OrganizationProvider
	users providers.UserProvider
}

func NewGetTourSteps(orgs providers.OrganizationProvider, users providers.UserProvider) GetTourSteps {
	return &getTourStepsImpl{orgs: orgs, users: users}
}

func (uc *getTourStepsImpl) Execute(ctx context.Context, req GetTourStepsRequest) ([]TourStep, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	org, err := uc.orgs.FindByID(ctx, req.OrgID)
	if err != nil {
		return nil, err
	}

	user, err := uc.users.FindByID(ctx, req.OrgID, req.UserID)
	if err != nil {
		return nil, err
	}

	cfg := entities.ParseOrgConfig(org.Config)
	if len(cfg.Onboarding.TourSteps) == 0 {
		return defaultTourSteps, nil
	}

	userRoles := user.RoleNames()

	var steps []TourStep
	seen := make(map[string]bool)

	for _, sc := range cfg.Onboarding.TourSteps {
		if seen[sc.Key] {
			continue
		}
		if !matchesRoles(sc.Roles, userRoles) {
			continue
		}
		if sc.RequiresFeature != "" && !cfg.IsFeatureActive(sc.RequiresFeature) {
			continue
		}
		seen[sc.Key] = true
		steps = append(steps, TourStep{
			Key:         sc.Key,
			Title:       sc.Title,
			Description: sc.Description,
			Order:       sc.Order,
		})
	}

	sort.Slice(steps, func(i, j int) bool {
		return steps[i].Order < steps[j].Order
	})

	return steps, nil
}

func matchesRoles(stepRoles []string, userRoles []string) bool {
	if len(stepRoles) == 0 {
		return true
	}
	for _, sr := range stepRoles {
		for _, ur := range userRoles {
			if sr == ur {
				return true
			}
		}
	}
	return false
}
