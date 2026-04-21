package admin_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestDeleteTopic_ValidationMissingOrg(t *testing.T) {
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewDeleteTopic(topics)

	err := uc.Execute(context.Background(), admin.DeleteTopicRequest{TopicID: 1})
	assert.ErrorIs(t, err, providers.ErrValidation)
	topics.AssertNotCalled(t, "CountTopicChildren", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteTopic_ValidationMissingID(t *testing.T) {
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewDeleteTopic(topics)

	err := uc.Execute(context.Background(), admin.DeleteTopicRequest{OrgID: uuid.New()})
	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestDeleteTopic_NotFound(t *testing.T) {
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewDeleteTopic(topics)

	orgID := uuid.New()
	topics.On("CountTopicChildren", mock.Anything, orgID, int64(99)).
		Return(int64(0), providers.ErrNotFound)

	err := uc.Execute(context.Background(), admin.DeleteTopicRequest{
		OrgID:   orgID,
		TopicID: 99,
	})
	assert.ErrorIs(t, err, providers.ErrNotFound)
	topics.AssertNotCalled(t, "DeleteTopic", mock.Anything, mock.Anything, mock.Anything)
}

// TestDeleteTopic_ConflictChildren locks in Option A: we refuse a delete while
// the subtree is non-empty. A future switch to Option B would need to update
// this test and the usecase together.
func TestDeleteTopic_ConflictChildren(t *testing.T) {
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewDeleteTopic(topics)

	orgID := uuid.New()
	topics.On("CountTopicChildren", mock.Anything, orgID, int64(1)).
		Return(int64(3), nil)

	err := uc.Execute(context.Background(), admin.DeleteTopicRequest{
		OrgID:   orgID,
		TopicID: 1,
	})
	assert.ErrorIs(t, err, providers.ErrConflict)
	assert.Contains(t, err.Error(), "3 child topics")
	topics.AssertNotCalled(t, "DeleteTopic", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteTopic_CounterError(t *testing.T) {
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewDeleteTopic(topics)

	orgID := uuid.New()
	boom := errors.New("db down")
	topics.On("CountTopicChildren", mock.Anything, orgID, int64(1)).
		Return(int64(0), boom)

	err := uc.Execute(context.Background(), admin.DeleteTopicRequest{
		OrgID:   orgID,
		TopicID: 1,
	})
	assert.ErrorIs(t, err, boom)
}

func TestDeleteTopic_HappyPath(t *testing.T) {
	topics := new(mockproviders.MockTopicProvider)
	uc := admin.NewDeleteTopic(topics)

	orgID := uuid.New()
	topics.On("CountTopicChildren", mock.Anything, orgID, int64(1)).
		Return(int64(0), nil)
	topics.On("DeleteTopic", mock.Anything, orgID, int64(1)).Return(nil)

	err := uc.Execute(context.Background(), admin.DeleteTopicRequest{
		OrgID:   orgID,
		TopicID: 1,
	})
	assert.NoError(t, err)
	topics.AssertExpectations(t)
}
