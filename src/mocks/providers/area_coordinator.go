package providers

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
)

type MockAreaCoordinatorProvider struct {
	mock.Mock
}

func (m *MockAreaCoordinatorProvider) Assign(ctx context.Context, areaID, userID int64) (*entities.AreaCoordinator, error) {
	args := m.Called(ctx, areaID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AreaCoordinator), args.Error(1)
}

func (m *MockAreaCoordinatorProvider) Remove(ctx context.Context, areaID, userID int64) error {
	args := m.Called(ctx, areaID, userID)
	return args.Error(0)
}

func (m *MockAreaCoordinatorProvider) FindByAreaID(ctx context.Context, areaID int64) ([]entities.AreaCoordinator, error) {
	args := m.Called(ctx, areaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.AreaCoordinator), args.Error(1)
}

func (m *MockAreaCoordinatorProvider) FindByUserID(ctx context.Context, userID int64) ([]entities.AreaCoordinator, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.AreaCoordinator), args.Error(1)
}

func (m *MockAreaCoordinatorProvider) IsCoordinator(ctx context.Context, areaID, userID int64) (bool, error) {
	args := m.Called(ctx, areaID, userID)
	return args.Bool(0), args.Error(1)
}
