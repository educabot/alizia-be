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

func TestUpdateOrgConfig_Success(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := admin.NewUpdateOrgConfig(orgs)

	orgID := uuid.New()
	ctx := context.Background()
	patch := map[string]any{"shared_classes_enabled": false, "desarrollo_max_activities": 5}

	existing := &entities.Organization{
		ID:     orgID,
		Name:   "Test Org",
		Slug:   "test-org",
		Config: datatypes.JSON(`{"topic_max_levels":3,"shared_classes_enabled":true}`),
	}
	updated := &entities.Organization{
		ID:     orgID,
		Name:   "Test Org",
		Slug:   "test-org",
		Config: datatypes.JSON(`{"topic_max_levels":3,"shared_classes_enabled":false,"desarrollo_max_activities":5}`),
	}

	orgs.On("FindByID", ctx, orgID).Return(existing, nil)
	orgs.On("UpdateConfig", ctx, orgID, patch).Return(updated, nil)

	result, err := uc.Execute(ctx, admin.UpdateOrgConfigRequest{
		OrgID:       orgID,
		ConfigPatch: patch,
	})

	assert.NoError(t, err)
	assert.Equal(t, updated, result)
	orgs.AssertExpectations(t)
}

func TestUpdateOrgConfig_ValidationError_NilOrgID(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := admin.NewUpdateOrgConfig(orgs)

	_, err := uc.Execute(context.Background(), admin.UpdateOrgConfigRequest{
		OrgID:       uuid.Nil,
		ConfigPatch: map[string]any{"key": "val"},
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestUpdateOrgConfig_ValidationError_EmptyPatch(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := admin.NewUpdateOrgConfig(orgs)

	_, err := uc.Execute(context.Background(), admin.UpdateOrgConfigRequest{
		OrgID:       uuid.New(),
		ConfigPatch: map[string]any{},
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestUpdateOrgConfig_OrgNotFound(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	uc := admin.NewUpdateOrgConfig(orgs)

	orgID := uuid.New()
	ctx := context.Background()

	orgs.On("FindByID", ctx, orgID).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, admin.UpdateOrgConfigRequest{
		OrgID:       orgID,
		ConfigPatch: map[string]any{"key": "val"},
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	orgs.AssertExpectations(t)
}
