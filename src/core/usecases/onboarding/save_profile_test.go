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

func orgWithProfileFields(fields []onboarding.ProfileFieldConfig) *entities.Organization {
	config := map[string]any{
		"onboarding": map[string]any{
			"profile_fields": fields,
		},
	}
	configJSON, _ := json.Marshal(config)
	return &entities.Organization{
		ID:     uuid.New(),
		Config: datatypes.JSON(configJSON),
	}
}

func TestSaveProfile_Success(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := onboarding.NewSaveProfile(users, orgs)

	orgID := uuid.New()
	ctx := context.Background()
	data := map[string]any{
		"disciplines":         []any{"Matemáticas"},
		"years_of_experience": float64(5),
	}

	org := orgWithProfileFields([]onboarding.ProfileFieldConfig{
		{Key: "disciplines", Type: "multiselect", Required: true, Options: []string{"Matemáticas", "Física"}},
		{Key: "years_of_experience", Type: "number", Required: false},
	})
	org.ID = orgID

	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{ID: 1}, nil)
	orgs.On("FindByID", ctx, orgID).Return(org, nil)
	users.On("UpdateProfileData", ctx, int64(1), data).Return(nil)

	err := uc.Execute(ctx, onboarding.SaveProfileRequest{OrgID: orgID, UserID: 1, Data: data})

	assert.NoError(t, err)
	users.AssertExpectations(t)
	orgs.AssertExpectations(t)
}

func TestSaveProfile_MissingRequiredField(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := onboarding.NewSaveProfile(users, orgs)

	orgID := uuid.New()
	ctx := context.Background()
	data := map[string]any{
		"years_of_experience": float64(5),
	}

	org := orgWithProfileFields([]onboarding.ProfileFieldConfig{
		{Key: "disciplines", Type: "multiselect", Required: true, Options: []string{"Matemáticas"}},
	})
	org.ID = orgID

	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{ID: 1}, nil)
	orgs.On("FindByID", ctx, orgID).Return(org, nil)

	err := uc.Execute(ctx, onboarding.SaveProfileRequest{OrgID: orgID, UserID: 1, Data: data})

	assert.ErrorIs(t, err, providers.ErrValidation)
	assert.Contains(t, err.Error(), "disciplines")
}

func TestSaveProfile_InvalidSelectOption(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := onboarding.NewSaveProfile(users, orgs)

	orgID := uuid.New()
	ctx := context.Background()
	data := map[string]any{
		"level": "Posgrado",
	}

	org := orgWithProfileFields([]onboarding.ProfileFieldConfig{
		{Key: "level", Type: "select", Required: true, Options: []string{"Primario", "Secundario"}},
	})
	org.ID = orgID

	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{ID: 1}, nil)
	orgs.On("FindByID", ctx, orgID).Return(org, nil)

	err := uc.Execute(ctx, onboarding.SaveProfileRequest{OrgID: orgID, UserID: 1, Data: data})

	assert.ErrorIs(t, err, providers.ErrValidation)
	assert.Contains(t, err.Error(), "invalid option")
}

func TestSaveProfile_InvalidFieldType(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := onboarding.NewSaveProfile(users, orgs)

	orgID := uuid.New()
	ctx := context.Background()
	data := map[string]any{
		"years": "not a number",
	}

	org := orgWithProfileFields([]onboarding.ProfileFieldConfig{
		{Key: "years", Type: "number", Required: true},
	})
	org.ID = orgID

	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{ID: 1}, nil)
	orgs.On("FindByID", ctx, orgID).Return(org, nil)

	err := uc.Execute(ctx, onboarding.SaveProfileRequest{OrgID: orgID, UserID: 1, Data: data})

	assert.ErrorIs(t, err, providers.ErrValidation)
	assert.Contains(t, err.Error(), "must be a number")
}

func TestSaveProfile_NoFieldsConfigured(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := onboarding.NewSaveProfile(users, orgs)

	orgID := uuid.New()
	ctx := context.Background()
	data := map[string]any{"anything": "goes"}

	org := &entities.Organization{
		ID:     orgID,
		Config: datatypes.JSON("{}"),
	}

	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{ID: 1}, nil)
	orgs.On("FindByID", ctx, orgID).Return(org, nil)
	users.On("UpdateProfileData", ctx, int64(1), data).Return(nil)

	err := uc.Execute(ctx, onboarding.SaveProfileRequest{OrgID: orgID, UserID: 1, Data: data})

	assert.NoError(t, err)
}

func TestSaveProfile_ValidationErrors(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := onboarding.NewSaveProfile(users, orgs)

	tests := []struct {
		name string
		req  onboarding.SaveProfileRequest
	}{
		{"missing org_id", onboarding.SaveProfileRequest{UserID: 1, Data: map[string]any{}}},
		{"missing user_id", onboarding.SaveProfileRequest{OrgID: uuid.New(), Data: map[string]any{}}},
		{"missing data", onboarding.SaveProfileRequest{OrgID: uuid.New(), UserID: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}
}
