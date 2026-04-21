package admin_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestCreateTopic_RootSuccess(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewCreateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()

	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID:     orgID,
		Config: datatypes.JSON(`{"topic_max_levels": 3}`),
	}, nil)
	topics.On("CreateTopic", ctx, mock.AnythingOfType("*entities.Topic")).Return(int64(1), nil)

	result, err := uc.Execute(ctx, admin.CreateTopicRequest{
		OrgID: orgID,
		Name:  "Pensamiento Lógico-Matemático",
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, 1, result.Level)
	assert.Nil(t, result.ParentID)
	orgs.AssertExpectations(t)
	topics.AssertExpectations(t)
}

func TestCreateTopic_ChildSuccess(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewCreateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()
	parentID := int64(1)

	topics.On("GetTopicByID", ctx, orgID, int64(1)).Return(&entities.Topic{
		ID: 1, OrganizationID: orgID, Level: 1,
	}, nil)
	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID:     orgID,
		Config: datatypes.JSON(`{"topic_max_levels": 3}`),
	}, nil)
	topics.On("CreateTopic", ctx, mock.AnythingOfType("*entities.Topic")).Return(int64(5), nil)

	result, err := uc.Execute(ctx, admin.CreateTopicRequest{
		OrgID:    orgID,
		ParentID: &parentID,
		Name:     "Aritmética Básica",
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(5), result.ID)
	assert.Equal(t, 2, result.Level)
	assert.Equal(t, &parentID, result.ParentID)
	topics.AssertExpectations(t)
}

func TestCreateTopic_ExceedsMaxLevels(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewCreateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()
	parentID := int64(3)

	// Parent is at level 3, and max is 3 -> child would be level 4
	topics.On("GetTopicByID", ctx, orgID, int64(3)).Return(&entities.Topic{
		ID: 3, OrganizationID: orgID, Level: 3,
	}, nil)
	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID:     orgID,
		Config: datatypes.JSON(`{"topic_max_levels": 3}`),
	}, nil)

	_, err := uc.Execute(ctx, admin.CreateTopicRequest{
		OrgID:    orgID,
		ParentID: &parentID,
		Name:     "Too deep",
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, providers.ErrTopicMaxLevel)
	assert.Contains(t, err.Error(), "exceeds maximum")
	topics.AssertNotCalled(t, "CreateTopic", mock.Anything, mock.Anything)
}

func TestCreateTopic_ParentNotFound(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewCreateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()
	parentID := int64(99)

	topics.On("GetTopicByID", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, admin.CreateTopicRequest{
		OrgID:    orgID,
		ParentID: &parentID,
		Name:     "Orphan",
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestCreateTopic_DefaultMaxLevels(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewCreateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()

	// No topic_max_levels in config -> should default to 3
	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID:     orgID,
		Config: datatypes.JSON(`{}`),
	}, nil)
	topics.On("CreateTopic", ctx, mock.AnythingOfType("*entities.Topic")).Return(int64(1), nil)

	result, err := uc.Execute(ctx, admin.CreateTopicRequest{
		OrgID: orgID,
		Name:  "Root topic",
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, result.Level)
}

func TestCreateTopic_ValidationErrors(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewCreateTopic(orgs, topics)

	tests := []struct {
		name string
		req  admin.CreateTopicRequest
	}{
		{"missing org_id", admin.CreateTopicRequest{Name: "Test"}},
		{"missing name", admin.CreateTopicRequest{OrgID: uuid.New()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}

	topics.AssertNotCalled(t, "CreateTopic", mock.Anything, mock.Anything)
}
