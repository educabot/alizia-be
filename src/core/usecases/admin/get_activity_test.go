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

func TestGetActivity_ValidationMissingOrg(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewGetActivity(activities)

	_, err := uc.Execute(context.Background(), admin.GetActivityRequest{ActivityID: 1})
	assert.ErrorIs(t, err, providers.ErrValidation)
	activities.AssertNotCalled(t, "GetActivity", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetActivity_ValidationMissingID(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewGetActivity(activities)

	_, err := uc.Execute(context.Background(), admin.GetActivityRequest{OrgID: uuid.New()})
	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestGetActivity_NotFound(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewGetActivity(activities)

	orgID := uuid.New()
	activities.On("GetActivity", mock.Anything, orgID, int64(99)).
		Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(context.Background(), admin.GetActivityRequest{
		OrgID:      orgID,
		ActivityID: 99,
	})
	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestGetActivity_HappyPath(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewGetActivity(activities)

	orgID := uuid.New()
	expected := &entities.ActivityTemplate{
		ID:             1,
		OrganizationID: orgID,
		Moment:         entities.MomentApertura,
		Name:           "Warm-up",
	}
	activities.On("GetActivity", mock.Anything, orgID, int64(1)).Return(expected, nil)

	result, err := uc.Execute(context.Background(), admin.GetActivityRequest{
		OrgID:      orgID,
		ActivityID: 1,
	})
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
