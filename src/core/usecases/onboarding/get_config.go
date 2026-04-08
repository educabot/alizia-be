package onboarding

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type GetConfigRequest struct {
	OrgID uuid.UUID
}

func (r GetConfigRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	return nil
}

type GetConfig interface {
	Execute(ctx context.Context, req GetConfigRequest) (*entities.OnboardingConfig, error)
}

type getConfigImpl struct {
	orgs providers.OrganizationProvider
}

func NewGetConfig(orgs providers.OrganizationProvider) GetConfig {
	return &getConfigImpl{orgs: orgs}
}

func (uc *getConfigImpl) Execute(ctx context.Context, req GetConfigRequest) (*entities.OnboardingConfig, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	org, err := uc.orgs.FindByID(ctx, req.OrgID)
	if err != nil {
		return nil, err
	}

	return parseOnboardingConfig(org.Config), nil
}

func parseOnboardingConfig(configJSON []byte) *entities.OnboardingConfig {
	var config struct {
		Onboarding entities.OnboardingConfig `json:"onboarding"`
	}
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return &entities.OnboardingConfig{SkipAllowed: true}
	}
	return &config.Onboarding
}
