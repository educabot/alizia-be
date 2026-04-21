package providers

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type MockAreaProvider struct {
	mock.Mock
}

func (m *MockAreaProvider) CreateArea(ctx context.Context, area *entities.Area) (int64, error) {
	args := m.Called(ctx, area)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAreaProvider) GetArea(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Area, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Area), args.Error(1)
}

func (m *MockAreaProvider) ListAreas(ctx context.Context, orgID uuid.UUID, p providers.Pagination) ([]entities.Area, bool, error) {
	args := m.Called(ctx, orgID, p)
	if args.Get(0) == nil {
		return nil, false, args.Error(2)
	}
	return args.Get(0).([]entities.Area), args.Bool(1), args.Error(2)
}

func (m *MockAreaProvider) UpdateArea(ctx context.Context, area *entities.Area) error {
	args := m.Called(ctx, area)
	return args.Error(0)
}

func (m *MockAreaProvider) CountDependencies(ctx context.Context, orgID uuid.UUID, id int64) (providers.AreaDependencies, error) {
	args := m.Called(ctx, orgID, id)
	return args.Get(0).(providers.AreaDependencies), args.Error(1)
}

func (m *MockAreaProvider) DeleteArea(ctx context.Context, orgID uuid.UUID, id int64) error {
	args := m.Called(ctx, orgID, id)
	return args.Error(0)
}
