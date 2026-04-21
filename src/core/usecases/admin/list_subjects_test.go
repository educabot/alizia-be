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

func TestListSubjects_Success(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewListSubjects(areas, subjects)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(1)).Return(&entities.Area{ID: 1, OrganizationID: orgID}, nil)
	expected := []entities.Subject{
		{ID: 1, OrganizationID: orgID, AreaID: 1, Name: "Matemáticas"},
		{ID: 2, OrganizationID: orgID, AreaID: 1, Name: "Física"},
	}
	subjects.On("ListSubjectsByArea", ctx, orgID, int64(1)).Return(expected, nil)

	result, err := uc.Execute(ctx, admin.ListSubjectsRequest{OrgID: orgID, AreaID: 1})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Matemáticas", result[0].Name)
	areas.AssertExpectations(t)
	subjects.AssertExpectations(t)
}

func TestListSubjects_AreaNotFound(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewListSubjects(areas, subjects)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(99)).Return(nil, fmt.Errorf("%w: area 99", providers.ErrNotFound))

	_, err := uc.Execute(ctx, admin.ListSubjectsRequest{OrgID: orgID, AreaID: 99})

	assert.Error(t, err)
	assert.ErrorIs(t, err, providers.ErrNotFound)
	subjects.AssertNotCalled(t, "ListSubjectsByArea", mock.Anything, mock.Anything, mock.Anything)
}

func TestListSubjects_ValidationErrors(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewListSubjects(areas, subjects)

	tests := []struct {
		name string
		req  admin.ListSubjectsRequest
	}{
		{"missing org_id", admin.ListSubjectsRequest{AreaID: 1}},
		{"missing area_id", admin.ListSubjectsRequest{OrgID: uuid.New()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}

	areas.AssertNotCalled(t, "GetArea", mock.Anything, mock.Anything, mock.Anything)
}
