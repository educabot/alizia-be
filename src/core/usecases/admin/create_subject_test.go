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

func TestCreateSubject_Success(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewCreateSubject(areas, subjects)

	orgID := uuid.New()
	ctx := context.Background()
	desc := "Álgebra, geometría y análisis"

	areas.On("GetArea", ctx, orgID, int64(1)).Return(&entities.Area{ID: 1, OrganizationID: orgID}, nil)
	subjects.On("CreateSubject", ctx, mock.AnythingOfType("*entities.Subject")).Return(int64(10), nil)

	result, err := uc.Execute(ctx, admin.CreateSubjectRequest{
		OrgID:       orgID,
		AreaID:      1,
		Name:        "Matemáticas",
		Description: &desc,
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(10), result.ID)
	assert.Equal(t, orgID, result.OrganizationID)
	assert.Equal(t, int64(1), result.AreaID)
	assert.Equal(t, "Matemáticas", result.Name)
	areas.AssertExpectations(t)
	subjects.AssertExpectations(t)
}

func TestCreateSubject_AreaNotFound(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewCreateSubject(areas, subjects)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(99)).Return(nil, fmt.Errorf("%w: area 99", providers.ErrNotFound))

	_, err := uc.Execute(ctx, admin.CreateSubjectRequest{
		OrgID:  orgID,
		AreaID: 99,
		Name:   "Física",
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, providers.ErrNotFound)
	subjects.AssertNotCalled(t, "CreateSubject", mock.Anything, mock.Anything)
}

func TestCreateSubject_ValidationErrors(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewCreateSubject(areas, subjects)

	tests := []struct {
		name string
		req  admin.CreateSubjectRequest
	}{
		{"missing org_id", admin.CreateSubjectRequest{AreaID: 1, Name: "Test"}},
		{"missing area_id", admin.CreateSubjectRequest{OrgID: uuid.New(), Name: "Test"}},
		{"missing name", admin.CreateSubjectRequest{OrgID: uuid.New(), AreaID: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}

	areas.AssertNotCalled(t, "GetArea", mock.Anything, mock.Anything, mock.Anything)
}
