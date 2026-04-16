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

func TestUpdateArea_UpdateNameOnly(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewUpdateArea(areas)

	orgID := uuid.New()
	ctx := context.Background()
	originalDesc := "original description"
	existing := &entities.Area{
		ID:             10,
		OrganizationID: orgID,
		Name:           "Old Name",
		Description:    &originalDesc,
	}

	areas.On("GetArea", ctx, orgID, int64(10)).Return(existing, nil)
	areas.On("UpdateArea", ctx, mock.MatchedBy(func(a *entities.Area) bool {
		return a.ID == 10 && a.Name == "New Name" && a.Description == &originalDesc
	})).Return(nil)

	newName := "New Name"
	result, err := uc.Execute(ctx, admin.UpdateAreaRequest{
		OrgID:  orgID,
		AreaID: 10,
		Name:   &newName,
	})

	assert.NoError(t, err)
	assert.Equal(t, "New Name", result.Name)
	assert.Equal(t, &originalDesc, result.Description)
	areas.AssertExpectations(t)
}

func TestUpdateArea_ClearDescription(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewUpdateArea(areas)

	orgID := uuid.New()
	ctx := context.Background()
	originalDesc := "to be cleared"
	existing := &entities.Area{
		ID:             5,
		OrganizationID: orgID,
		Name:           "Sciences",
		Description:    &originalDesc,
	}

	areas.On("GetArea", ctx, orgID, int64(5)).Return(existing, nil)
	areas.On("UpdateArea", ctx, mock.MatchedBy(func(a *entities.Area) bool {
		return a.ID == 5 && a.Description == nil
	})).Return(nil)

	result, err := uc.Execute(ctx, admin.UpdateAreaRequest{
		OrgID:          orgID,
		AreaID:         5,
		SetDescription: true,
		Description:    nil,
	})

	assert.NoError(t, err)
	assert.Nil(t, result.Description)
	areas.AssertExpectations(t)
}

func TestUpdateArea_UpdateBothFields(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewUpdateArea(areas)

	orgID := uuid.New()
	ctx := context.Background()
	existing := &entities.Area{
		ID:             7,
		OrganizationID: orgID,
		Name:           "Old",
		Description:    nil,
	}

	areas.On("GetArea", ctx, orgID, int64(7)).Return(existing, nil)
	areas.On("UpdateArea", ctx, mock.AnythingOfType("*entities.Area")).Return(nil)

	newName := "Brand New"
	newDesc := "Brand new description"
	result, err := uc.Execute(ctx, admin.UpdateAreaRequest{
		OrgID:          orgID,
		AreaID:         7,
		Name:           &newName,
		SetDescription: true,
		Description:    &newDesc,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Brand New", result.Name)
	assert.Equal(t, &newDesc, result.Description)
	areas.AssertExpectations(t)
}

func TestUpdateArea_NoChanges(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewUpdateArea(areas)

	orgID := uuid.New()
	ctx := context.Background()
	existing := &entities.Area{
		ID:             3,
		OrganizationID: orgID,
		Name:           "Untouched",
	}

	areas.On("GetArea", ctx, orgID, int64(3)).Return(existing, nil)
	areas.On("UpdateArea", ctx, mock.MatchedBy(func(a *entities.Area) bool {
		return a.ID == 3 && a.Name == "Untouched"
	})).Return(nil)

	result, err := uc.Execute(ctx, admin.UpdateAreaRequest{
		OrgID:  orgID,
		AreaID: 3,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Untouched", result.Name)
	areas.AssertExpectations(t)
}

func TestUpdateArea_ValidationErrors(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewUpdateArea(areas)

	emptyName := ""
	tests := []struct {
		name string
		req  admin.UpdateAreaRequest
	}{
		{"missing org_id", admin.UpdateAreaRequest{AreaID: 1}},
		{"missing area_id", admin.UpdateAreaRequest{OrgID: uuid.New()}},
		{"empty name", admin.UpdateAreaRequest{OrgID: uuid.New(), AreaID: 1, Name: &emptyName}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}

	areas.AssertNotCalled(t, "GetArea", mock.Anything, mock.Anything, mock.Anything)
	areas.AssertNotCalled(t, "UpdateArea", mock.Anything, mock.Anything)
}

func TestUpdateArea_NotFound(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewUpdateArea(areas)

	orgID := uuid.New()
	ctx := context.Background()

	areas.On("GetArea", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	newName := "x"
	_, err := uc.Execute(ctx, admin.UpdateAreaRequest{
		OrgID:  orgID,
		AreaID: 99,
		Name:   &newName,
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	areas.AssertNotCalled(t, "UpdateArea", mock.Anything, mock.Anything)
}

func TestUpdateArea_RepoError(t *testing.T) {
	areas := new(mockproviders.MockAreaProvider)
	uc := admin.NewUpdateArea(areas)

	orgID := uuid.New()
	ctx := context.Background()
	existing := &entities.Area{ID: 1, OrganizationID: orgID, Name: "x"}

	areas.On("GetArea", ctx, orgID, int64(1)).Return(existing, nil)
	areas.On("UpdateArea", ctx, mock.AnythingOfType("*entities.Area")).Return(errors.New("db down"))

	newName := "y"
	_, err := uc.Execute(ctx, admin.UpdateAreaRequest{
		OrgID:  orgID,
		AreaID: 1,
		Name:   &newName,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db down")
}
