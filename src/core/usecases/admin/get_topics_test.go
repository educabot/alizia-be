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

func TestGetTopics_TreeSuccess(t *testing.T) {
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewGetTopics(topics)

	orgID := uuid.New()
	ctx := context.Background()

	expected := []entities.Topic{
		{ID: 1, Name: "Root", Level: 1, Children: []entities.Topic{
			{ID: 2, Name: "Child", Level: 2, Children: []entities.Topic{}},
		}},
	}
	topics.On("GetTopicTree", ctx, orgID).Return(expected, nil)

	result, err := uc.Execute(ctx, admin.GetTopicsRequest{OrgID: orgID})

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Root", result[0].Name)
	assert.Len(t, result[0].Children, 1)
	topics.AssertExpectations(t)
}

func TestGetTopics_ByLevel(t *testing.T) {
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewGetTopics(topics)

	orgID := uuid.New()
	ctx := context.Background()
	level := 2

	expected := []entities.Topic{
		{ID: 2, Name: "Level 2 Topic A", Level: 2},
		{ID: 3, Name: "Level 2 Topic B", Level: 2},
	}
	topics.On("GetTopicsByLevel", ctx, orgID, 2).Return(expected, nil)

	result, err := uc.Execute(ctx, admin.GetTopicsRequest{OrgID: orgID, Level: &level})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	topics.AssertExpectations(t)
	topics.AssertNotCalled(t, "GetTopicTree", mock.Anything, mock.Anything)
}

func TestGetTopics_ValidationError(t *testing.T) {
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewGetTopics(topics)

	_, err := uc.Execute(context.Background(), admin.GetTopicsRequest{})
	assert.ErrorIs(t, err, providers.ErrValidation)
	topics.AssertNotCalled(t, "GetTopicTree", mock.Anything, mock.Anything)
}

func TestGetTopics_InvalidLevel(t *testing.T) {
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewGetTopics(topics)

	level := 0
	_, err := uc.Execute(context.Background(), admin.GetTopicsRequest{
		OrgID: uuid.New(),
		Level: &level,
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestGetTopics_ByParent(t *testing.T) {
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewGetTopics(topics)

	orgID := uuid.New()
	ctx := context.Background()
	pid := int64(7)

	expected := []entities.Topic{
		{ID: 10, Name: "Child A", Level: 2, ParentID: &pid},
		{ID: 11, Name: "Child B", Level: 2, ParentID: &pid},
	}
	topics.On("GetTopicsByParent", ctx, orgID, &pid).Return(expected, nil)

	result, err := uc.Execute(ctx, admin.GetTopicsRequest{
		OrgID:     orgID,
		ParentID:  &pid,
		SetParent: true,
	})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Child A", result[0].Name)
	topics.AssertExpectations(t)
	topics.AssertNotCalled(t, "GetTopicTree", mock.Anything, mock.Anything)
	topics.AssertNotCalled(t, "GetTopicsByLevel", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetTopics_ByRootParent(t *testing.T) {
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewGetTopics(topics)

	orgID := uuid.New()
	ctx := context.Background()

	expected := []entities.Topic{
		{ID: 1, Name: "Root A", Level: 1},
		{ID: 2, Name: "Root B", Level: 1},
	}
	topics.On("GetTopicsByParent", ctx, orgID, (*int64)(nil)).Return(expected, nil)

	result, err := uc.Execute(ctx, admin.GetTopicsRequest{
		OrgID:     orgID,
		SetParent: true,
	})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	topics.AssertExpectations(t)
	topics.AssertNotCalled(t, "GetTopicTree", mock.Anything, mock.Anything)
}
