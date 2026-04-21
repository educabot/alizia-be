package admin_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestDeleteActivity_ValidationMissingOrg(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewDeleteActivity(activities)

	err := uc.Execute(context.Background(), admin.DeleteActivityRequest{ActivityID: 1})
	assert.ErrorIs(t, err, providers.ErrValidation)
	activities.AssertNotCalled(t, "DeleteActivity", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteActivity_ValidationMissingID(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewDeleteActivity(activities)

	err := uc.Execute(context.Background(), admin.DeleteActivityRequest{OrgID: uuid.New()})
	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestDeleteActivity_NotFound(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewDeleteActivity(activities)

	orgID := uuid.New()
	activities.On("DeleteActivity", mock.Anything, orgID, int64(99)).
		Return(providers.ErrNotFound)

	err := uc.Execute(context.Background(), admin.DeleteActivityRequest{
		OrgID:      orgID,
		ActivityID: 99,
	})
	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestDeleteActivity_RepoError(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewDeleteActivity(activities)

	orgID := uuid.New()
	boom := errors.New("db down")
	activities.On("DeleteActivity", mock.Anything, orgID, int64(1)).Return(boom)

	err := uc.Execute(context.Background(), admin.DeleteActivityRequest{
		OrgID:      orgID,
		ActivityID: 1,
	})
	assert.ErrorIs(t, err, boom)
}

func TestDeleteActivity_HappyPath(t *testing.T) {
	activities := new(mockproviders.MockActivityTemplateProvider)
	uc := admin.NewDeleteActivity(activities)

	orgID := uuid.New()
	activities.On("DeleteActivity", mock.Anything, orgID, int64(1)).Return(nil)

	err := uc.Execute(context.Background(), admin.DeleteActivityRequest{
		OrgID:      orgID,
		ActivityID: 1,
	})
	assert.NoError(t, err)
	activities.AssertExpectations(t)
}
