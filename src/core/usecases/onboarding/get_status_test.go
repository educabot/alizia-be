package onboarding_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/onboarding"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestGetStatus_NotCompleted(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetStatus(users)

	orgID := uuid.New()
	ctx := context.Background()

	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{
		ID: 1, OnboardingCompletedAt: nil,
	}, nil)

	result, err := uc.Execute(ctx, onboarding.GetStatusRequest{OrgID: orgID, UserID: 1})

	assert.NoError(t, err)
	assert.False(t, result.Completed)
	assert.Nil(t, result.CompletedAt)
	users.AssertExpectations(t)
}

func TestGetStatus_Completed(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetStatus(users)

	orgID := uuid.New()
	ctx := context.Background()
	completedAt := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)

	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{
		ID: 1, OnboardingCompletedAt: &completedAt,
	}, nil)

	result, err := uc.Execute(ctx, onboarding.GetStatusRequest{OrgID: orgID, UserID: 1})

	assert.NoError(t, err)
	assert.True(t, result.Completed)
	assert.NotNil(t, result.CompletedAt)
	users.AssertExpectations(t)
}

func TestGetStatus_UserNotFound(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetStatus(users)

	orgID := uuid.New()
	ctx := context.Background()

	users.On("FindByID", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, onboarding.GetStatusRequest{OrgID: orgID, UserID: 99})

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestGetStatus_ValidationErrors(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewGetStatus(users)

	tests := []struct {
		name string
		req  onboarding.GetStatusRequest
	}{
		{"missing org_id", onboarding.GetStatusRequest{UserID: 1}},
		{"missing user_id", onboarding.GetStatusRequest{OrgID: uuid.New()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}
}
