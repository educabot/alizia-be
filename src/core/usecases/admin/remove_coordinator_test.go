package admin_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestRemoveCoordinator_Success(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	coords := new(mockproviders.MockAreaCoordinatorProvider)
	uc := admin.NewRemoveCoordinator(areas, coords)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(1)).Return(&entities.Area{ID: 1, OrganizationID: orgID}, nil)
	coords.On("Remove", ctx, int64(1), int64(2)).Return(nil)

	err := uc.Execute(ctx, admin.RemoveCoordinatorRequest{
		OrgID: orgID, AreaID: 1, UserID: 2,
	})

	assert.NoError(t, err)
	areas.AssertExpectations(t)
	coords.AssertExpectations(t)
}

func TestRemoveCoordinator_AreaNotFound(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	coords := new(mockproviders.MockAreaCoordinatorProvider)
	uc := admin.NewRemoveCoordinator(areas, coords)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(99)).Return(nil, fmt.Errorf("%w: area 99", providers.ErrNotFound))

	err := uc.Execute(ctx, admin.RemoveCoordinatorRequest{
		OrgID: orgID, AreaID: 99, UserID: 2,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestRemoveCoordinator_AssignmentNotFound(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	coords := new(mockproviders.MockAreaCoordinatorProvider)
	uc := admin.NewRemoveCoordinator(areas, coords)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(1)).Return(&entities.Area{ID: 1, OrganizationID: orgID}, nil)
	coords.On("Remove", ctx, int64(1), int64(2)).Return(fmt.Errorf("%w: assignment not found", providers.ErrNotFound))

	err := uc.Execute(ctx, admin.RemoveCoordinatorRequest{
		OrgID: orgID, AreaID: 1, UserID: 2,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, providers.ErrNotFound)
}
