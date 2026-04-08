package onboarding_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/onboarding"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestComplete_Success(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewComplete(users)

	orgID := uuid.New()
	ctx := context.Background()

	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{ID: 1}, nil)
	users.On("CompleteOnboarding", ctx, int64(1)).Return(nil)

	err := uc.Execute(ctx, onboarding.CompleteRequest{OrgID: orgID, UserID: 1})

	assert.NoError(t, err)
	users.AssertExpectations(t)
}

func TestComplete_UserNotFound(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewComplete(users)

	orgID := uuid.New()
	ctx := context.Background()

	users.On("FindByID", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	err := uc.Execute(ctx, onboarding.CompleteRequest{OrgID: orgID, UserID: 99})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	users.AssertNotCalled(t, "CompleteOnboarding")
}

func TestComplete_Idempotent(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewComplete(users)

	orgID := uuid.New()
	ctx := context.Background()

	users.On("FindByID", ctx, orgID, int64(1)).Return(&entities.User{ID: 1}, nil)
	users.On("CompleteOnboarding", ctx, int64(1)).Return(nil)

	// First call
	err := uc.Execute(ctx, onboarding.CompleteRequest{OrgID: orgID, UserID: 1})
	assert.NoError(t, err)

	// Second call — should still succeed (idempotent)
	err = uc.Execute(ctx, onboarding.CompleteRequest{OrgID: orgID, UserID: 1})
	assert.NoError(t, err)
}

func TestComplete_ValidationErrors(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := onboarding.NewComplete(users)

	tests := []struct {
		name string
		req  onboarding.CompleteRequest
	}{
		{"missing org_id", onboarding.CompleteRequest{UserID: 1}},
		{"missing user_id", onboarding.CompleteRequest{OrgID: uuid.New()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}
}
