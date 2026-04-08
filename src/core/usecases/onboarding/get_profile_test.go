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

func TestGetProfile_Empty(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetProfile(users)

	orgID := uuid.New()
	ctx := context.Background()

	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{
		ID:          1,
		ProfileData: datatypes.JSON("{}"),
	}, nil)

	result, err := uc.Execute(ctx, onboarding.GetProfileRequest{OrgID: orgID, UserID: 1})

	assert.NoError(t, err)
	assert.Empty(t, result)
	users.AssertExpectations(t)
}

func TestGetProfile_WithData(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetProfile(users)

	orgID := uuid.New()
	ctx := context.Background()

	profileJSON, _ := json.Marshal(map[string]any{
		"disciplines":         []string{"Matemáticas", "Física"},
		"years_of_experience": 12,
	})

	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{
		ID:          1,
		ProfileData: datatypes.JSON(profileJSON),
	}, nil)

	result, err := uc.Execute(ctx, onboarding.GetProfileRequest{OrgID: orgID, UserID: 1})

	assert.NoError(t, err)
	assert.Contains(t, result, "disciplines")
	assert.Contains(t, result, "years_of_experience")
	users.AssertExpectations(t)
}

func TestGetProfile_UserNotFound(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetProfile(users)

	orgID := uuid.New()
	ctx := context.Background()

	users.On("FindByID", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, onboarding.GetProfileRequest{OrgID: orgID, UserID: 99})

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestGetProfile_ValidationErrors(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetProfile(users)

	tests := []struct {
		name string
		req  onboarding.GetProfileRequest
	}{
		{"missing org_id", onboarding.GetProfileRequest{UserID: 1}},
		{"missing user_id", onboarding.GetProfileRequest{OrgID: uuid.New()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}
}
