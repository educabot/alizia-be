package admin_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestListActivities_All(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewListActivities(activities)

	orgID := uuid.New()
	ctx := context.Background()

	expected := []entities.ActivityTemplate{
		{ID: 1, Moment: entities.MomentApertura, Name: "Lluvia de ideas"},
		{ID: 2, Moment: entities.MomentDesarrollo, Name: "Trabajo en grupo"},
	}
	var nilMoment *entities.ClassMoment
	activities.On("ListActivities", ctx, orgID, nilMoment).Return(expected, nil)

	result, err := uc.Execute(ctx, admin.ListActivitiesRequest{OrgID: orgID})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	activities.AssertExpectations(t)
}

func TestListActivities_ByMoment(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewListActivities(activities)

	orgID := uuid.New()
	ctx := context.Background()
	moment := "apertura"
	classMoment := entities.MomentApertura

	expected := []entities.ActivityTemplate{
		{ID: 1, Moment: entities.MomentApertura, Name: "Lluvia de ideas"},
	}
	activities.On("ListActivities", ctx, orgID, &classMoment).Return(expected, nil)

	result, err := uc.Execute(ctx, admin.ListActivitiesRequest{OrgID: orgID, Moment: &moment})

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	activities.AssertExpectations(t)
}

func TestListActivities_InvalidMoment(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewListActivities(activities)

	invalid := "invalid"
	_, err := uc.Execute(context.Background(), admin.ListActivitiesRequest{
		OrgID: uuid.New(), Moment: &invalid,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestListActivities_ValidationError(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewListActivities(activities)

	_, err := uc.Execute(context.Background(), admin.ListActivitiesRequest{})
	assert.ErrorIs(t, err, providers.ErrValidation)
	activities.AssertNotCalled(t, "ListActivities", mock.Anything, mock.Anything, mock.Anything)
}
