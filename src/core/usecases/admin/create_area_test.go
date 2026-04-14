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

func TestCreateArea_Success(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewCreateArea(areas)

	orgID := uuid.New()
	ctx := context.Background()
	desc := "Mathematics and Physics"

	areas.On("CreateArea", ctx, mock.AnythingOfType("*entities.Area")).Return(int64(1), nil)

	result, err := uc.Execute(ctx, admin.CreateAreaRequest{
		OrgID:       orgID,
		Name:        "Ciencias Exactas",
		Description: &desc,
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, orgID, result.OrganizationID)
	assert.Equal(t, "Ciencias Exactas", result.Name)
	assert.Equal(t, &desc, result.Description)
	areas.AssertExpectations(t)
}

func TestCreateArea_SuccessWithoutDescription(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewCreateArea(areas)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("CreateArea", ctx, mock.AnythingOfType("*entities.Area")).Return(int64(2), nil)

	result, err := uc.Execute(ctx, admin.CreateAreaRequest{
		OrgID: orgID,
		Name:  "Humanidades",
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(2), result.ID)
	assert.Nil(t, result.Description)
	areas.AssertExpectations(t)
}

func TestCreateArea_ValidationErrors(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewCreateArea(areas)

	tests := []struct {
		name string
		req  admin.CreateAreaRequest
	}{
		{"missing org_id", admin.CreateAreaRequest{Name: "Test"}},
		{"missing name", admin.CreateAreaRequest{OrgID: uuid.New()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}

	areas.AssertNotCalled(t, "CreateArea", mock.Anything, mock.Anything)
}

func TestCreateArea_RepoError(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewCreateArea(areas)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("CreateArea", ctx, mock.AnythingOfType("*entities.Area")).Return(int64(0), errors.New("db error"))

	_, err := uc.Execute(ctx, admin.CreateAreaRequest{
		OrgID: orgID,
		Name:  "Test",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}
