package admin_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestCreateActivity_Success(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewCreateActivity(activities)

	orgID := uuid.New()
	ctx := context.Background()
	dur := 15

	activities.On("CreateActivity", ctx, mock.AnythingOfType("*entities.ActivityTemplate")).Return(int64(1), nil)

	result, err := uc.Execute(ctx, admin.CreateActivityRequest{
		OrgID:           orgID,
		Moment:          "apertura",
		Name:            "Lluvia de ideas",
		DurationMinutes: &dur,
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Lluvia de ideas", result.Name)
	activities.AssertExpectations(t)
}

func TestCreateActivity_InvalidMoment(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewCreateActivity(activities)

	_, err := uc.Execute(context.Background(), admin.CreateActivityRequest{
		OrgID: uuid.New(), Moment: "invalid", Name: "Test",
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	assert.Contains(t, err.Error(), "moment")
}

func TestCreateActivity_ValidationErrors(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewCreateActivity(activities)

	tests := []struct {
		name string
		req  admin.CreateActivityRequest
	}{
		{"missing org_id", admin.CreateActivityRequest{Moment: "apertura", Name: "Test"}},
		{"missing name", admin.CreateActivityRequest{OrgID: uuid.New(), Moment: "apertura"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}
	activities.AssertNotCalled(t, "CreateActivity", mock.Anything, mock.Anything)
}

// TestCreateActivity_DoesNotCapCatalog guards against regressing to enforcing
// config.desarrollo_max_activities at template creation time. Per RFC HU-3.6 /
// HU-5.3 that limit is runtime-only (lesson plan per-class pick), not a cap on
// how many templates an admin may register.
func TestCreateActivity_DoesNotCapCatalog(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewCreateActivity(activities)
	ctx := context.Background()

	activities.On("CreateActivity", ctx, mock.AnythingOfType("*entities.ActivityTemplate")).
		Return(int64(42), nil)

	_, err := uc.Execute(ctx, admin.CreateActivityRequest{
		OrgID: uuid.New(), Moment: "desarrollo", Name: "Cuarta desarrollo",
	})

	assert.NoError(t, err)
	activities.AssertNotCalled(t, "CountByMoment", mock.Anything, mock.Anything, mock.Anything)
}
