package providers

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type MockTopicProvider struct {
	mock.Mock
}

func (m *MockTopicProvider) CreateTopic(ctx context.Context, topic *entities.Topic) (int64, error) {
	args := m.Called(ctx, topic)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTopicProvider) GetTopicByID(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Topic, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Topic), args.Error(1)
}

func (m *MockTopicProvider) GetTopicTree(ctx context.Context, orgID uuid.UUID) ([]entities.Topic, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Topic), args.Error(1)
}

func (m *MockTopicProvider) GetTopicsByLevel(ctx context.Context, orgID uuid.UUID, level int, p providers.Pagination) ([]entities.Topic, bool, error) {
	args := m.Called(ctx, orgID, level, p)
	if args.Get(0) == nil {
		return nil, args.Bool(1), args.Error(2)
	}
	return args.Get(0).([]entities.Topic), args.Bool(1), args.Error(2)
}

func (m *MockTopicProvider) GetTopicsByParent(ctx context.Context, orgID uuid.UUID, parentID *int64, p providers.Pagination) ([]entities.Topic, bool, error) {
	args := m.Called(ctx, orgID, parentID, p)
	if args.Get(0) == nil {
		return nil, args.Bool(1), args.Error(2)
	}
	return args.Get(0).([]entities.Topic), args.Bool(1), args.Error(2)
}

func (m *MockTopicProvider) ListAllTopics(ctx context.Context, orgID uuid.UUID) ([]entities.Topic, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Topic), args.Error(1)
}

func (m *MockTopicProvider) UpdateTopic(ctx context.Context, topic *entities.Topic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockTopicProvider) UpdateTopicLevels(ctx context.Context, orgID uuid.UUID, levels map[int64]int) error {
	args := m.Called(ctx, orgID, levels)
	return args.Error(0)
}
