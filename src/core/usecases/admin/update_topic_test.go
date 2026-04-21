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

func intPtr(v int64) *int64   { return &v }
func strPtr(s string) *string { return &s }

// Tree used across move-related tests:
//
//	1 (L1)
//	├── 2 (L2)
//	│   ├── 4 (L3)
//	│   └── 5 (L3)
//	└── 3 (L2)
//
// We rebuild this in each test so a mutation in one does not leak.
func sampleTree(orgID uuid.UUID) []entities.Topic {
	return []entities.Topic{
		{ID: 1, OrganizationID: orgID, Level: 1},
		{ID: 2, OrganizationID: orgID, Level: 2, ParentID: intPtr(1)},
		{ID: 3, OrganizationID: orgID, Level: 2, ParentID: intPtr(1)},
		{ID: 4, OrganizationID: orgID, Level: 3, ParentID: intPtr(2)},
		{ID: 5, OrganizationID: orgID, Level: 3, ParentID: intPtr(2)},
	}
}

func TestUpdateTopic_FieldOnlyUpdate_SkipsTreeWork(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewUpdateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()

	topics.On("GetTopicByID", ctx, orgID, int64(2)).Return(&entities.Topic{
		ID: 2, OrganizationID: orgID, Level: 2, ParentID: intPtr(1),
	}, nil)
	topics.On("UpdateTopic", ctx, mock.AnythingOfType("*entities.Topic")).Return(nil)

	newName := "Álgebra"
	result, err := uc.Execute(ctx, admin.UpdateTopicRequest{
		OrgID: orgID, TopicID: 2, Name: &newName,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Álgebra", result.Name)
	assert.Equal(t, 2, result.Level)
	topics.AssertNotCalled(t, "ListAllTopics", mock.Anything, mock.Anything)
	topics.AssertNotCalled(t, "UpdateTopicLevels", mock.Anything, mock.Anything, mock.Anything)
	orgs.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestUpdateTopic_RejectsSelfParent(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewUpdateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()
	selfID := int64(2)

	topics.On("GetTopicByID", ctx, orgID, int64(2)).Return(&entities.Topic{
		ID: 2, OrganizationID: orgID, Level: 2, ParentID: intPtr(1),
	}, nil)

	_, err := uc.Execute(ctx, admin.UpdateTopicRequest{
		OrgID: orgID, TopicID: 2, SetParent: true, ParentID: &selfID,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	assert.Contains(t, err.Error(), "own parent")
}

func TestUpdateTopic_RejectsCycleToDescendant(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewUpdateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()
	descendantID := int64(4) // 4 is a descendant of 1

	topics.On("GetTopicByID", ctx, orgID, int64(1)).Return(&entities.Topic{
		ID: 1, OrganizationID: orgID, Level: 1,
	}, nil)
	// The move targets topic 4 as parent; the UC still resolves it as parent first.
	topics.On("GetTopicByID", ctx, orgID, int64(4)).Return(&entities.Topic{
		ID: 4, OrganizationID: orgID, Level: 3, ParentID: intPtr(2),
	}, nil)
	topics.On("ListAllTopics", ctx, orgID).Return(sampleTree(orgID), nil)

	_, err := uc.Execute(ctx, admin.UpdateTopicRequest{
		OrgID: orgID, TopicID: 1, SetParent: true, ParentID: &descendantID,
	})

	assert.ErrorIs(t, err, providers.ErrValidation)
	assert.Contains(t, err.Error(), "descendants")
	topics.AssertNotCalled(t, "UpdateTopic", mock.Anything, mock.Anything)
	topics.AssertNotCalled(t, "UpdateTopicLevels", mock.Anything, mock.Anything, mock.Anything)
}

func TestUpdateTopic_MovePromotesToRoot(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewUpdateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()

	topics.On("GetTopicByID", ctx, orgID, int64(2)).Return(&entities.Topic{
		ID: 2, OrganizationID: orgID, Level: 2, ParentID: intPtr(1),
	}, nil)
	topics.On("ListAllTopics", ctx, orgID).Return(sampleTree(orgID), nil)
	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID: orgID, Config: datatypes.JSON(`{"topic_max_levels":3}`),
	}, nil)
	topics.On("UpdateTopic", ctx, mock.MatchedBy(func(t *entities.Topic) bool {
		return t.ID == 2 && t.ParentID == nil && t.Level == 1
	})).Return(nil)
	// Moving 2 from L2 to L1 shifts descendants 4 and 5 from L3 to L2.
	topics.On("UpdateTopicLevels", ctx, orgID, mock.MatchedBy(func(m map[int64]int) bool {
		return m[2] == 1 && m[4] == 2 && m[5] == 2 && len(m) == 3
	})).Return(nil)

	result, err := uc.Execute(ctx, admin.UpdateTopicRequest{
		OrgID: orgID, TopicID: 2, SetParent: true, ParentID: nil,
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, result.Level)
	assert.Nil(t, result.ParentID)
	topics.AssertExpectations(t)
}

func TestUpdateTopic_MoveExceedsMaxLevelsInSubtree(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewUpdateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()
	newParent := int64(3)

	// Moving 2 (L2) under 3 (L2) -> 2 becomes L3, descendants 4/5 become L4.
	// With max=3 the subtree check must fail before anything persists.
	topics.On("GetTopicByID", ctx, orgID, int64(2)).Return(&entities.Topic{
		ID: 2, OrganizationID: orgID, Level: 2, ParentID: intPtr(1),
	}, nil)
	topics.On("GetTopicByID", ctx, orgID, int64(3)).Return(&entities.Topic{
		ID: 3, OrganizationID: orgID, Level: 2, ParentID: intPtr(1),
	}, nil)
	topics.On("ListAllTopics", ctx, orgID).Return(sampleTree(orgID), nil)
	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID: orgID, Config: datatypes.JSON(`{"topic_max_levels":3}`),
	}, nil)

	_, err := uc.Execute(ctx, admin.UpdateTopicRequest{
		OrgID: orgID, TopicID: 2, SetParent: true, ParentID: &newParent,
	})

	assert.ErrorIs(t, err, providers.ErrTopicMaxLevel)
	assert.Contains(t, err.Error(), "subtree depth")
	topics.AssertNotCalled(t, "UpdateTopic", mock.Anything, mock.Anything)
	topics.AssertNotCalled(t, "UpdateTopicLevels", mock.Anything, mock.Anything, mock.Anything)
}

func TestUpdateTopic_MoveExceedsMaxLevelsAtNode(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewUpdateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()
	newParent := int64(4) // L3, and the moved node is a childless leaf so no subtree depth issue

	// Leaf topic 6 has no descendants so only the node check should fire.
	topics.On("GetTopicByID", ctx, orgID, int64(6)).Return(&entities.Topic{
		ID: 6, OrganizationID: orgID, Level: 2, ParentID: intPtr(1),
	}, nil)
	topics.On("GetTopicByID", ctx, orgID, int64(4)).Return(&entities.Topic{
		ID: 4, OrganizationID: orgID, Level: 3, ParentID: intPtr(2),
	}, nil)
	tree := append(sampleTree(orgID), entities.Topic{
		ID: 6, OrganizationID: orgID, Level: 2, ParentID: intPtr(1),
	})
	topics.On("ListAllTopics", ctx, orgID).Return(tree, nil)
	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID: orgID, Config: datatypes.JSON(`{"topic_max_levels":3}`),
	}, nil)

	_, err := uc.Execute(ctx, admin.UpdateTopicRequest{
		OrgID: orgID, TopicID: 6, SetParent: true, ParentID: &newParent,
	})

	assert.ErrorIs(t, err, providers.ErrTopicMaxLevel)
	assert.Contains(t, err.Error(), "level 4")
}

func TestUpdateTopic_MoveSucceedsWhenWithinLimit(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewUpdateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()
	newParent := int64(3)

	// With max=4, moving 2 under 3 makes 2→L3, 4/5→L4. All fit.
	topics.On("GetTopicByID", ctx, orgID, int64(2)).Return(&entities.Topic{
		ID: 2, OrganizationID: orgID, Level: 2, ParentID: intPtr(1),
	}, nil)
	topics.On("GetTopicByID", ctx, orgID, int64(3)).Return(&entities.Topic{
		ID: 3, OrganizationID: orgID, Level: 2, ParentID: intPtr(1),
	}, nil)
	topics.On("ListAllTopics", ctx, orgID).Return(sampleTree(orgID), nil)
	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID: orgID, Config: datatypes.JSON(`{"topic_max_levels":4}`),
	}, nil)
	topics.On("UpdateTopic", ctx, mock.MatchedBy(func(t *entities.Topic) bool {
		return t.ID == 2 && t.ParentID != nil && *t.ParentID == 3 && t.Level == 3
	})).Return(nil)
	topics.On("UpdateTopicLevels", ctx, orgID, mock.MatchedBy(func(m map[int64]int) bool {
		return m[2] == 3 && m[4] == 4 && m[5] == 4 && len(m) == 3
	})).Return(nil)

	result, err := uc.Execute(ctx, admin.UpdateTopicRequest{
		OrgID: orgID, TopicID: 2, SetParent: true, ParentID: &newParent,
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, result.Level)
	assert.Equal(t, int64(3), *result.ParentID)
	topics.AssertExpectations(t)
}

func TestUpdateTopic_UsesDefaultMaxLevelsWhenConfigMissing(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewUpdateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()
	newParent := int64(4)

	// Leaf promoted under L3 -> L4. Default max is 3, so this must fail.
	topics.On("GetTopicByID", ctx, orgID, int64(6)).Return(&entities.Topic{
		ID: 6, OrganizationID: orgID, Level: 1, ParentID: nil,
	}, nil)
	topics.On("GetTopicByID", ctx, orgID, int64(4)).Return(&entities.Topic{
		ID: 4, OrganizationID: orgID, Level: 3, ParentID: intPtr(2),
	}, nil)
	tree := append(sampleTree(orgID), entities.Topic{ID: 6, OrganizationID: orgID, Level: 1})
	topics.On("ListAllTopics", ctx, orgID).Return(tree, nil)
	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID: orgID, Config: datatypes.JSON(`{}`),
	}, nil)

	_, err := uc.Execute(ctx, admin.UpdateTopicRequest{
		OrgID: orgID, TopicID: 6, SetParent: true, ParentID: &newParent,
	})

	assert.ErrorIs(t, err, providers.ErrTopicMaxLevel)
}

func TestUpdateTopic_ParentNotFound(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewUpdateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()
	missingParent := int64(99)

	topics.On("GetTopicByID", ctx, orgID, int64(2)).Return(&entities.Topic{
		ID: 2, OrganizationID: orgID, Level: 2, ParentID: intPtr(1),
	}, nil)
	topics.On("GetTopicByID", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, admin.UpdateTopicRequest{
		OrgID: orgID, TopicID: 2, SetParent: true, ParentID: &missingParent,
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestUpdateTopic_TopicNotFound(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewUpdateTopic(orgs, topics)

	orgID := uuid.New()
	ctx := context.Background()

	topics.On("GetTopicByID", ctx, orgID, int64(42)).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, admin.UpdateTopicRequest{
		OrgID: orgID, TopicID: 42, Name: strPtr("x"),
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestUpdateTopic_ValidationErrors(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewUpdateTopic(orgs, topics)

	tests := []struct {
		name string
		req  admin.UpdateTopicRequest
	}{
		{"missing org_id", admin.UpdateTopicRequest{TopicID: 1}},
		{"missing topic_id", admin.UpdateTopicRequest{OrgID: uuid.New()}},
		{"empty name", admin.UpdateTopicRequest{OrgID: uuid.New(), TopicID: 1, Name: strPtr("")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}

	topics.AssertNotCalled(t, "UpdateTopic", mock.Anything, mock.Anything)
}
