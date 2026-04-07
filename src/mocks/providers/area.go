package providers

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
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

func (m *MockAreaProvider) ListAreas(ctx context.Context, orgID uuid.UUID) ([]entities.Area, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Area), args.Error(1)
}
