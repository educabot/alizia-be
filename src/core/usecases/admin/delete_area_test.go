package admin_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestDeleteArea_Success(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewDeleteArea(areas)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(10)).Return(&entities.Area{ID: 10, OrganizationID: orgID}, nil)
	areas.On("CountDependencies", ctx, orgID, int64(10)).Return(providers.AreaDependencies{}, nil)
	areas.On("DeleteArea", ctx, orgID, int64(10)).Return(nil)

	err := uc.Execute(ctx, admin.DeleteAreaRequest{OrgID: orgID, AreaID: 10})

	assert.NoError(t, err)
	areas.AssertExpectations(t)
}

func TestDeleteArea_NotFound(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewDeleteArea(areas)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	err := uc.Execute(ctx, admin.DeleteAreaRequest{OrgID: orgID, AreaID: 99})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	areas.AssertNotCalled(t, "CountDependencies", mock.Anything, mock.Anything, mock.Anything)
	areas.AssertNotCalled(t, "DeleteArea", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteArea_BlockedBySubjects(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewDeleteArea(areas)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(10)).Return(&entities.Area{ID: 10}, nil)
	areas.On("CountDependencies", ctx, orgID, int64(10)).
		Return(providers.AreaDependencies{Subjects: 3}, nil)

	err := uc.Execute(ctx, admin.DeleteAreaRequest{OrgID: orgID, AreaID: 10})

	assert.ErrorIs(t, err, providers.ErrConflict)
	assert.Contains(t, err.Error(), "3 subjects")
	areas.AssertNotCalled(t, "DeleteArea", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteArea_ValidationErrors(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewDeleteArea(areas)

	tests := []struct {
		name string
		req  admin.DeleteAreaRequest
	}{
		{"missing org_id", admin.DeleteAreaRequest{AreaID: 1}},
		{"missing area_id", admin.DeleteAreaRequest{OrgID: uuid.New()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}

	areas.AssertNotCalled(t, "GetArea", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteArea_RepoError(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewDeleteArea(areas)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(10)).Return(&entities.Area{ID: 10}, nil)
	areas.On("CountDependencies", ctx, orgID, int64(10)).Return(providers.AreaDependencies{}, nil)
	areas.On("DeleteArea", ctx, orgID, int64(10)).Return(errors.New("db down"))

	err := uc.Execute(ctx, admin.DeleteAreaRequest{OrgID: orgID, AreaID: 10})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db down")
}
