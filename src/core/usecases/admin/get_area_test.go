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

func TestGetArea_Success(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewGetArea(areas)

	orgID := uuid.New()
	ctx := context.Background()

	expected := &entities.Area{ID: 42, OrganizationID: orgID, Name: "Ciencias Exactas"}
	areas.On("GetArea", ctx, orgID, int64(42)).Return(expected, nil)

	result, err := uc.Execute(ctx, admin.GetAreaRequest{OrgID: orgID, AreaID: 42})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	areas.AssertExpectations(t)
}

func TestGetArea_NotFound(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewGetArea(areas)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, admin.GetAreaRequest{OrgID: orgID, AreaID: 99})

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestGetArea_ValidationErrors(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewGetArea(areas)

	tests := []struct {
		name string
		req  admin.GetAreaRequest
	}{
		{"missing org_id", admin.GetAreaRequest{AreaID: 1}},
		{"missing area_id", admin.GetAreaRequest{OrgID: uuid.New()}},
		{"negative area_id", admin.GetAreaRequest{OrgID: uuid.New(), AreaID: -1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}

	areas.AssertNotCalled(t, "GetArea", mock.Anything, mock.Anything, mock.Anything)
}
