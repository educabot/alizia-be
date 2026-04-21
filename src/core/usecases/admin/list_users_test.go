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

func TestListUsers_ValidationMissingOrg(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := admin.NewListUsers(users)

	_, err := uc.Execute(context.Background(), admin.ListUsersRequest{})
	assert.ErrorIs(t, err, providers.ErrValidation)
	users.AssertNotCalled(t, "ListUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestListUsers_ValidationInvalidRole(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := admin.NewListUsers(users)

	role := "boss"
	_, err := uc.Execute(context.Background(), admin.ListUsersRequest{OrgID: uuid.New(), Role: &role})
	assert.ErrorIs(t, err, providers.ErrValidation)
	users.AssertNotCalled(t, "ListUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestListUsers_ValidationInvalidAreaID(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := admin.NewListUsers(users)

	bad := int64(0)
	_, err := uc.Execute(context.Background(), admin.ListUsersRequest{OrgID: uuid.New(), AreaID: &bad})
	assert.ErrorIs(t, err, providers.ErrValidation)
}

// TestListUsers_NoFilters verifies the zero-filter path passes an empty
// UserFilter and threads the pagination verbatim.
func TestListUsers_NoFilters(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := admin.NewListUsers(users)

	orgID := uuid.New()
	ctx := context.Background()

	expected := []entities.User{{ID: 1, Email: "a@b.c"}, {ID: 2, Email: "c@d.e"}}
	users.On("ListUsers", ctx, orgID, providers.UserFilter{}, providers.Pagination{}).
		Return(expected, false, nil)

	result, err := uc.Execute(ctx, admin.ListUsersRequest{OrgID: orgID})
	assert.NoError(t, err)
	assert.Len(t, result.Items, 2)
	assert.False(t, result.More)
	users.AssertExpectations(t)
}

// TestListUsers_AllFiltersAndPagination ensures every optional field is
// forwarded into the provider call untouched.
func TestListUsers_AllFiltersAndPagination(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := admin.NewListUsers(users)

	orgID := uuid.New()
	ctx := context.Background()
	role := string(entities.RoleTeacher)
	areaID := int64(5)
	search := "juan"
	page := providers.Pagination{Limit: 25, Offset: 50}

	expectedRole := entities.RoleTeacher
	expectedFilter := providers.UserFilter{
		Role:   &expectedRole,
		AreaID: &areaID,
		Search: &search,
	}

	expected := []entities.User{{ID: 7, FirstName: "Juan"}}
	users.On("ListUsers", ctx, orgID, expectedFilter, page).Return(expected, true, nil)

	result, err := uc.Execute(ctx, admin.ListUsersRequest{
		OrgID:      orgID,
		Role:       &role,
		AreaID:     &areaID,
		Search:     &search,
		Pagination: page,
	})
	assert.NoError(t, err)
	assert.True(t, result.More, "more flag must bubble up from provider")
	assert.Equal(t, "Juan", result.Items[0].FirstName)
	users.AssertExpectations(t)
}

func TestListUsers_ProviderError(t *testing.T) {
	users := new(mockproviders.MockUserProvider)
	uc := admin.NewListUsers(users)

	orgID := uuid.New()
	boom := errors.New("db down")
	users.On("ListUsers", mock.Anything, orgID, mock.Anything, mock.Anything).
		Return(nil, false, boom)

	_, err := uc.Execute(context.Background(), admin.ListUsersRequest{OrgID: orgID})
	assert.ErrorIs(t, err, boom)
}
