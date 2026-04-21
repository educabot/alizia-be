package admin_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestListAllSubjects_SuccessNoAreaFilter(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewListAllSubjects(areas, subjects)

	orgID := uuid.New()
	ctx := context.Background()

	expected := []entities.Subject{
		{ID: 1, OrganizationID: orgID, AreaID: 1, Name: "Matemáticas"},
		{ID: 2, OrganizationID: orgID, AreaID: 2, Name: "Historia"},
	}
	subjects.On("ListSubjectsByOrg", ctx, orgID, (*int64)(nil), mock.AnythingOfType("providers.Pagination")).Return(expected, false, nil)

	result, err := uc.Execute(ctx, admin.ListAllSubjectsRequest{OrgID: orgID})

	assert.NoError(t, err)
	assert.Len(t, result.Items, 2)
	assert.False(t, result.More)
	subjects.AssertExpectations(t)
	areas.AssertNotCalled(t, "GetArea", mock.Anything, mock.Anything, mock.Anything)
}

func TestListAllSubjects_SuccessFilteredByArea(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewListAllSubjects(areas, subjects)

	orgID := uuid.New()
	ctx := context.Background()
	areaID := int64(1)

	areas.On("GetArea", ctx, orgID, areaID).Return(&entities.Area{ID: 1, OrganizationID: orgID}, nil)
	expected := []entities.Subject{
		{ID: 1, OrganizationID: orgID, AreaID: 1, Name: "Matemáticas"},
	}
	subjects.On("ListSubjectsByOrg", ctx, orgID, &areaID, mock.AnythingOfType("providers.Pagination")).Return(expected, false, nil)

	result, err := uc.Execute(ctx, admin.ListAllSubjectsRequest{OrgID: orgID, AreaID: &areaID})

	assert.NoError(t, err)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, "Matemáticas", result.Items[0].Name)
	areas.AssertExpectations(t)
	subjects.AssertExpectations(t)
}

func TestListAllSubjects_WithPagination(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewListAllSubjects(areas, subjects)

	orgID := uuid.New()
	ctx := context.Background()

	expected := []entities.Subject{
		{ID: 1, OrganizationID: orgID, AreaID: 1, Name: "Matemáticas"},
	}
	subjects.On("ListSubjectsByOrg", ctx, orgID, (*int64)(nil), mock.AnythingOfType("providers.Pagination")).Return(expected, true, nil)

	result, err := uc.Execute(ctx, admin.ListAllSubjectsRequest{
		OrgID: orgID,
		Page:  providers.Pagination{Limit: 1, Offset: 0},
	})

	assert.NoError(t, err)
	assert.Len(t, result.Items, 1)
	assert.True(t, result.More)
	subjects.AssertExpectations(t)
}

func TestListAllSubjects_ValidationMissingOrg(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewListAllSubjects(areas, subjects)

	_, err := uc.Execute(context.Background(), admin.ListAllSubjectsRequest{})

	assert.ErrorIs(t, err, providers.ErrValidation)
	areas.AssertNotCalled(t, "GetArea", mock.Anything, mock.Anything, mock.Anything)
	subjects.AssertNotCalled(t, "ListSubjectsByOrg", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestListAllSubjects_AreaNotFound(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewListAllSubjects(areas, subjects)

	orgID := uuid.New()
	ctx := context.Background()
	areaID := int64(99)

	areas.On("GetArea", ctx, orgID, areaID).Return(nil, fmt.Errorf("%w: area 99", providers.ErrNotFound))

	_, err := uc.Execute(ctx, admin.ListAllSubjectsRequest{OrgID: orgID, AreaID: &areaID})

	assert.Error(t, err)
	assert.ErrorIs(t, err, providers.ErrNotFound)
	subjects.AssertNotCalled(t, "ListSubjectsByOrg", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}
