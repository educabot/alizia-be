package onboarding_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/onboarding"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestGetConfig_WithFullConfig(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := onboarding.NewGetConfig(orgs)

	orgID := uuid.New()
	ctx := context.Background()

	configJSON, _ := json.Marshal(map[string]any{
		"onboarding": map[string]any{
			"skip_allowed": false,
			"profile_fields": []map[string]any{
				{"key": "disciplines", "label": "Disciplinas", "type": "multiselect", "options": []string{"Matemáticas"}, "required": true},
			},
			"tour_steps": []map[string]any{
				{"key": "welcome", "title": "Bienvenido", "description": "Intro", "order": 1, "roles": []string{"coordinator", "teacher"}},
			},
		},
	})

	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID:     orgID,
		Config: datatypes.JSON(configJSON),
	}, nil)

	result, err := uc.Execute(ctx, onboarding.GetConfigRequest{OrgID: orgID})

	assert.NoError(t, err)
	assert.False(t, result.SkipAllowed)
	assert.Len(t, result.ProfileFields, 1)
	assert.Equal(t, "disciplines", result.ProfileFields[0].Key)
	assert.Len(t, result.TourSteps, 1)
	assert.Equal(t, "welcome", result.TourSteps[0].Key)
}

func TestGetConfig_EmptyConfig(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := onboarding.NewGetConfig(orgs)

	orgID := uuid.New()
	ctx := context.Background()

	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID:     orgID,
		Config: datatypes.JSON("{}"),
	}, nil)

	result, err := uc.Execute(ctx, onboarding.GetConfigRequest{OrgID: orgID})

	assert.NoError(t, err)
	assert.False(t, result.SkipAllowed)
	assert.Empty(t, result.ProfileFields)
	assert.Empty(t, result.TourSteps)
}

func TestGetConfig_SkipAllowed(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := onboarding.NewGetConfig(orgs)

	orgID := uuid.New()
	ctx := context.Background()

	configJSON, _ := json.Marshal(map[string]any{
		"onboarding": map[string]any{
			"skip_allowed": true,
		},
	})

	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID:     orgID,
		Config: datatypes.JSON(configJSON),
	}, nil)

	result, err := uc.Execute(ctx, onboarding.GetConfigRequest{OrgID: orgID})

	assert.NoError(t, err)
	assert.True(t, result.SkipAllowed)
}

func TestGetConfig_OrgNotFound(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := onboarding.NewGetConfig(orgs)

	orgID := uuid.New()
	ctx := context.Background()

	orgs.On("FindByID", ctx, orgID).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, onboarding.GetConfigRequest{OrgID: orgID})

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestGetConfig_ValidationError(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := onboarding.NewGetConfig(orgs)

	_, err := uc.Execute(context.Background(), onboarding.GetConfigRequest{})

	assert.ErrorIs(t, err, providers.ErrValidation)
}
