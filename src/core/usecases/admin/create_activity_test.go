package admin_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestCreateActivity_Success(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewCreateActivity(orgs, activities)

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
	orgs := new(mockproviders.MockOrganizationProvider)
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewCreateActivity(orgs, activities)

	_, err := uc.Execute(context.Background(), admin.CreateActivityRequest{
		OrgID: uuid.New(), Moment: "invalid", Name: "Test",
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	assert.Contains(t, err.Error(), "moment")
}

func TestCreateActivity_ValidationErrors(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewCreateActivity(orgs, activities)

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

func TestCreateActivity_DesarrolloLimitReached(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewCreateActivity(orgs, activities)

	orgID := uuid.New()
	ctx := context.Background()

	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID:     orgID,
		Config: datatypes.JSON([]byte(`{"desarrollo_max_activities": 3}`)),
	}, nil)
	activities.On("CountByMoment", ctx, orgID, entities.MomentDesarrollo).Return(int64(3), nil)

	_, err := uc.Execute(ctx, admin.CreateActivityRequest{
		OrgID: orgID, Moment: "desarrollo", Name: "Cuarta",
	})

	assert.ErrorIs(t, err, providers.ErrActivityMaxReached)
	activities.AssertNotCalled(t, "CreateActivity", mock.Anything, mock.Anything)
}

func TestCreateActivity_DesarrolloUnderLimit(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewCreateActivity(orgs, activities)

	orgID := uuid.New()
	ctx := context.Background()

	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID:     orgID,
		Config: datatypes.JSON([]byte(`{"desarrollo_max_activities": 3}`)),
	}, nil)
	activities.On("CountByMoment", ctx, orgID, entities.MomentDesarrollo).Return(int64(2), nil)
	activities.On("CreateActivity", ctx, mock.AnythingOfType("*entities.ActivityTemplate")).Return(int64(7), nil)

	result, err := uc.Execute(ctx, admin.CreateActivityRequest{
		OrgID: orgID, Moment: "desarrollo", Name: "Tercera",
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(7), result.ID)
	activities.AssertExpectations(t)
}
