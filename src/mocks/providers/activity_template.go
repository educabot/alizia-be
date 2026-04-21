package providers

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type MockActivityTemplateProvider struct {
	mock.Mock
}

func (m *MockActivityTemplateProvider) CreateActivity(ctx context.Context, activity *entities.ActivityTemplate) (int64, error) {
	args := m.Called(ctx, activity)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockActivityTemplateProvider) ListActivities(ctx context.Context, orgID uuid.UUID, moment *entities.ClassMoment, p providers.Pagination) ([]entities.ActivityTemplate, bool, error) {
	args := m.Called(ctx, orgID, moment, p)
	if args.Get(0) == nil {
		return nil, args.Bool(1), args.Error(2)
	}
	return args.Get(0).([]entities.ActivityTemplate), args.Bool(1), args.Error(2)
}

func (m *MockActivityTemplateProvider) CountByMoment(ctx context.Context, orgID uuid.UUID, moment entities.ClassMoment) (int64, error) {
	args := m.Called(ctx, orgID, moment)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockActivityTemplateProvider) GetActivity(ctx context.Context, orgID uuid.UUID, id int64) (*entities.ActivityTemplate, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ActivityTemplate), args.Error(1)
}

func (m *MockActivityTemplateProvider) UpdateActivity(ctx context.Context, activity *entities.ActivityTemplate) error {
	args := m.Called(ctx, activity)
	return args.Error(0)
}

func (m *MockActivityTemplateProvider) DeleteActivity(ctx context.Context, orgID uuid.UUID, id int64) error {
	args := m.Called(ctx, orgID, id)
	return args.Error(0)
}
