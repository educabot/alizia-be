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

func TestAssignCoordinator_Success(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	users := new(mockproviders.MockUserProvider)
	coords := new(mockproviders.MockAreaCoordinatorProvider)
	uc := admin.NewAssignCoordinator(areas, users, coords)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(1)).Return(&entities.Area{ID: 1, OrganizationID: orgID}, nil)
	users.On("FindByID", ctx, orgID, int64(2)).Return(&entities.User{
		ID: 2, Roles: []entities.UserRole{{Role: entities.RoleCoordinator}},
	}, nil)
	coords.On("Assign", ctx, int64(1), int64(2)).Return(&entities.AreaCoordinator{
		ID: 1, AreaID: 1, UserID: 2,
	}, nil)

	result, err := uc.Execute(ctx, admin.AssignCoordinatorRequest{
		OrgID: orgID, AreaID: 1, UserID: 2,
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.AreaID)
	assert.Equal(t, int64(2), result.UserID)
	areas.AssertExpectations(t)
	users.AssertExpectations(t)
	coords.AssertExpectations(t)
}

func TestAssignCoordinator_AreaNotFound(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	users := new(mockproviders.MockUserProvider)
	coords := new(mockproviders.MockAreaCoordinatorProvider)
	uc := admin.NewAssignCoordinator(areas, users, coords)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(99)).Return(nil, fmt.Errorf("%w: area 99", providers.ErrNotFound))

	_, err := uc.Execute(ctx, admin.AssignCoordinatorRequest{
		OrgID: orgID, AreaID: 99, UserID: 2,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestAssignCoordinator_UserNotInOrg(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	users := new(mockproviders.MockUserProvider)
	coords := new(mockproviders.MockAreaCoordinatorProvider)
	uc := admin.NewAssignCoordinator(areas, users, coords)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(1)).Return(&entities.Area{ID: 1, OrganizationID: orgID}, nil)
	users.On("FindByID", ctx, orgID, int64(99)).Return(nil, fmt.Errorf("%w: user 99", providers.ErrNotFound))

	_, err := uc.Execute(ctx, admin.AssignCoordinatorRequest{
		OrgID: orgID, AreaID: 1, UserID: 99,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestAssignCoordinator_UserNotCoordinator(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	users := new(mockproviders.MockUserProvider)
	coords := new(mockproviders.MockAreaCoordinatorProvider)
	uc := admin.NewAssignCoordinator(areas, users, coords)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(1)).Return(&entities.Area{ID: 1, OrganizationID: orgID}, nil)
	users.On("FindByID", ctx, orgID, int64(2)).Return(&entities.User{
		ID: 2, Roles: []entities.UserRole{{Role: entities.RoleTeacher}},
	}, nil)

	_, err := uc.Execute(ctx, admin.AssignCoordinatorRequest{
		OrgID: orgID, AreaID: 1, UserID: 2,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestAssignCoordinator_AlreadyAssigned(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	users := new(mockproviders.MockUserProvider)
	coords := new(mockproviders.MockAreaCoordinatorProvider)
	uc := admin.NewAssignCoordinator(areas, users, coords)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(1)).Return(&entities.Area{ID: 1, OrganizationID: orgID}, nil)
	users.On("FindByID", ctx, orgID, int64(2)).Return(&entities.User{
		ID: 2, Roles: []entities.UserRole{{Role: entities.RoleCoordinator}},
	}, nil)
	coords.On("Assign", ctx, int64(1), int64(2)).Return(nil, fmt.Errorf("%w: coordinator already assigned to area", providers.ErrConflict))

	_, err := uc.Execute(ctx, admin.AssignCoordinatorRequest{
		OrgID: orgID, AreaID: 1, UserID: 2,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, providers.ErrConflict)
}

func TestAssignCoordinator_ValidationErrors(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	users := new(mockproviders.MockUserProvider)
	coords := new(mockproviders.MockAreaCoordinatorProvider)
	uc := admin.NewAssignCoordinator(areas, users, coords)

	tests := []struct {
		name string
		req  admin.AssignCoordinatorRequest
	}{
		{"missing org_id", admin.AssignCoordinatorRequest{AreaID: 1, UserID: 2}},
		{"missing area_id", admin.AssignCoordinatorRequest{OrgID: uuid.New(), UserID: 2}},
		{"missing user_id", admin.AssignCoordinatorRequest{OrgID: uuid.New(), AreaID: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}

	// No calls should have been made
	areas.AssertNotCalled(t, "GetArea", mock.Anything, mock.Anything, mock.Anything)
}
