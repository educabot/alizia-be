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

func TestListAreas_Success(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewListAreas(areas)

	orgID := uuid.New()
	ctx := context.Background()

	expected := []entities.Area{
		{ID: 1, OrganizationID: orgID, Name: "Ciencias Exactas"},
		{ID: 2, OrganizationID: orgID, Name: "Humanidades"},
	}
	areas.On("ListAreas", ctx, orgID, mock.AnythingOfType("providers.Pagination")).Return(expected, false, nil)

	result, err := uc.Execute(ctx, admin.ListAreasRequest{OrgID: orgID})

	assert.NoError(t, err)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, "Ciencias Exactas", result.Items[0].Name)
	assert.Equal(t, "Humanidades", result.Items[1].Name)
	assert.False(t, result.More)
	areas.AssertExpectations(t)
}

func TestListAreas_WithPagination(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewListAreas(areas)

	orgID := uuid.New()
	ctx := context.Background()

	expected := []entities.Area{
		{ID: 1, OrganizationID: orgID, Name: "Ciencias Exactas"},
	}
	areas.On("ListAreas", ctx, orgID, mock.AnythingOfType("providers.Pagination")).Return(expected, true, nil)

	result, err := uc.Execute(ctx, admin.ListAreasRequest{
		OrgID: orgID,
		Page:  providers.Pagination{Limit: 1, Offset: 0},
	})

	assert.NoError(t, err)
	assert.Len(t, result.Items, 1)
	assert.True(t, result.More)
	areas.AssertExpectations(t)
}

func TestListAreas_Empty(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewListAreas(areas)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("ListAreas", ctx, orgID, mock.AnythingOfType("providers.Pagination")).Return([]entities.Area{}, false, nil)

	result, err := uc.Execute(ctx, admin.ListAreasRequest{OrgID: orgID})

	assert.NoError(t, err)
	assert.Empty(t, result.Items)
	assert.False(t, result.More)
	areas.AssertExpectations(t)
}

func TestListAreas_ValidationError(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewListAreas(areas)

	_, err := uc.Execute(context.Background(), admin.ListAreasRequest{})

	assert.ErrorIs(t, err, providers.ErrValidation)
	areas.AssertNotCalled(t, "ListAreas", mock.Anything, mock.Anything, mock.Anything)
}
