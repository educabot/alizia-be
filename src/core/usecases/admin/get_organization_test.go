package admin_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestGetOrganization_Success(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := admin.NewGetOrganization(orgs)

	orgID := uuid.New()
	ctx := context.Background()
	expected := &entities.Organization{
		ID:     orgID,
		Name:   "Test Org",
		Slug:   "test-org",
		Config: datatypes.JSON(`{"topic_max_levels":3}`),
	}

	orgs.On("FindByID", ctx, orgID).Return(expected, nil)

	result, err := uc.Execute(ctx, admin.GetOrganizationRequest{OrgID: orgID})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	orgs.AssertExpectations(t)
}

func TestGetOrganization_ValidationError(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := admin.NewGetOrganization(orgs)

	_, err := uc.Execute(context.Background(), admin.GetOrganizationRequest{
		OrgID: uuid.Nil,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestGetOrganization_NotFound(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := admin.NewGetOrganization(orgs)

	orgID := uuid.New()
	ctx := context.Background()

	orgs.On("FindByID", ctx, orgID).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, admin.GetOrganizationRequest{OrgID: orgID})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	orgs.AssertExpectations(t)
}
