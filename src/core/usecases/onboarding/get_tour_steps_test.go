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

func orgWithTourConfig(tourSteps []map[string]any, features map[string]bool) *entities.Organization {
	config := map[string]any{
		"onboarding": map[string]any{
			"tour_steps": tourSteps,
		},
	}
	if features != nil {
		config["features"] = features
	}
	configJSON, _ := json.Marshal(config)
	return &entities.Organization{
		ID:     uuid.New(),
		Config: datatypes.JSON(configJSON),
	}
}

func TestGetTourSteps_DefaultWhenNoConfig(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetTourSteps(orgs, users)

	orgID := uuid.New()
	ctx := context.Background()

	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID:     orgID,
		Config: datatypes.JSON("{}"),
	}, nil)
	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{
		ID: 1, Roles: []entities.UserRole{{Role: entities.RoleTeacher}},
	}, nil)

	steps, err := uc.Execute(ctx, onboarding.GetTourStepsRequest{OrgID: orgID, UserID: 1})

	assert.NoError(t, err)
	assert.Len(t, steps, 2)
	assert.Equal(t, "welcome", steps[0].Key)
	assert.Equal(t, "explore", steps[1].Key)
}

func TestGetTourSteps_FilterByRole(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetTourSteps(orgs, users)

	orgID := uuid.New()
	ctx := context.Background()

	org := orgWithTourConfig([]map[string]any{
		{"key": "welcome", "title": "Welcome", "description": "Intro", "order": 1},
		{"key": "areas", "title": "Your areas", "description": "For coordinators", "order": 2, "roles": []string{"coordinator"}},
		{"key": "courses", "title": "Your courses", "description": "For teachers", "order": 3, "roles": []string{"teacher"}},
	}, nil)
	org.ID = orgID

	orgs.On("FindByID", ctx, orgID).Return(org, nil)
	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{
		ID: 1, Roles: []entities.UserRole{{Role: entities.RoleTeacher}},
	}, nil)

	steps, err := uc.Execute(ctx, onboarding.GetTourStepsRequest{OrgID: orgID, UserID: 1})

	assert.NoError(t, err)
	assert.Len(t, steps, 2)
	assert.Equal(t, "welcome", steps[0].Key)
	assert.Equal(t, "courses", steps[1].Key)
}

func TestGetTourSteps_MultipleRolesDedup(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetTourSteps(orgs, users)

	orgID := uuid.New()
	ctx := context.Background()

	org := orgWithTourConfig([]map[string]any{
		{"key": "welcome", "title": "Welcome", "description": "Intro", "order": 1},
		{"key": "welcome", "title": "Welcome dup", "description": "Dup", "order": 1, "roles": []string{"teacher"}},
		{"key": "areas", "title": "Areas", "description": "Coord", "order": 2, "roles": []string{"coordinator"}},
		{"key": "courses", "title": "Courses", "description": "Teacher", "order": 3, "roles": []string{"teacher"}},
	}, nil)
	org.ID = orgID

	orgs.On("FindByID", ctx, orgID).Return(org, nil)
	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{
		ID: 1, Roles: []entities.UserRole{
			{Role: entities.RoleCoordinator},
			{Role: entities.RoleTeacher},
		},
	}, nil)

	steps, err := uc.Execute(ctx, onboarding.GetTourStepsRequest{OrgID: orgID, UserID: 1})

	assert.NoError(t, err)
	assert.Len(t, steps, 3)
	keys := make([]string, len(steps))
	for i, s := range steps {
		keys[i] = s.Key
	}
	assert.Equal(t, []string{"welcome", "areas", "courses"}, keys)
}

func TestGetTourSteps_FilterByFeatureFlag(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetTourSteps(orgs, users)

	orgID := uuid.New()
	ctx := context.Background()

	org := orgWithTourConfig([]map[string]any{
		{"key": "welcome", "title": "Welcome", "description": "Intro", "order": 1},
		{"key": "shared", "title": "Shared classes", "description": "Shared", "order": 2, "requires_feature": "shared_classes_enabled"},
	}, map[string]bool{"shared_classes_enabled": false})
	org.ID = orgID

	orgs.On("FindByID", ctx, orgID).Return(org, nil)
	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{
		ID: 1, Roles: []entities.UserRole{{Role: entities.RoleTeacher}},
	}, nil)

	steps, err := uc.Execute(ctx, onboarding.GetTourStepsRequest{OrgID: orgID, UserID: 1})

	assert.NoError(t, err)
	assert.Len(t, steps, 1)
	assert.Equal(t, "welcome", steps[0].Key)
}

func TestGetTourSteps_FeatureEnabled(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetTourSteps(orgs, users)

	orgID := uuid.New()
	ctx := context.Background()

	org := orgWithTourConfig([]map[string]any{
		{"key": "welcome", "title": "Welcome", "description": "Intro", "order": 1},
		{"key": "shared", "title": "Shared classes", "description": "Shared", "order": 2, "requires_feature": "shared_classes_enabled"},
	}, map[string]bool{"shared_classes_enabled": true})
	org.ID = orgID

	orgs.On("FindByID", ctx, orgID).Return(org, nil)
	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{
		ID: 1, Roles: []entities.UserRole{{Role: entities.RoleTeacher}},
	}, nil)

	steps, err := uc.Execute(ctx, onboarding.GetTourStepsRequest{OrgID: orgID, UserID: 1})

	assert.NoError(t, err)
	assert.Len(t, steps, 2)
}

func TestGetTourSteps_SortedByOrder(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetTourSteps(orgs, users)

	orgID := uuid.New()
	ctx := context.Background()

	org := orgWithTourConfig([]map[string]any{
		{"key": "third", "title": "Third", "description": "3", "order": 3},
		{"key": "first", "title": "First", "description": "1", "order": 1},
		{"key": "second", "title": "Second", "description": "2", "order": 2},
	}, nil)
	org.ID = orgID

	orgs.On("FindByID", ctx, orgID).Return(org, nil)
	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{
		ID: 1, Roles: []entities.UserRole{{Role: entities.RoleTeacher}},
	}, nil)

	steps, err := uc.Execute(ctx, onboarding.GetTourStepsRequest{OrgID: orgID, UserID: 1})

	assert.NoError(t, err)
	assert.Equal(t, "first", steps[0].Key)
	assert.Equal(t, "second", steps[1].Key)
	assert.Equal(t, "third", steps[2].Key)
}

func TestGetTourSteps_ValidationErrors(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetTourSteps(orgs, users)

	tests := []struct {
		name string
		req  onboarding.GetTourStepsRequest
	}{
		{"missing org_id", onboarding.GetTourStepsRequest{UserID: 1}},
		{"missing user_id", onboarding.GetTourStepsRequest{OrgID: uuid.New()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}
}

func TestGetTourSteps_OrgNotFound(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetTourSteps(orgs, users)

	orgID := uuid.New()
	ctx := context.Background()

	orgs.On("FindByID", ctx, orgID).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, onboarding.GetTourStepsRequest{OrgID: orgID, UserID: 1})

	assert.ErrorIs(t, err, providers.ErrNotFound)
}
