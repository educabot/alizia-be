package providers

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
)

type MockActivityTemplateProvider struct {
	mock.Mock
}

func (m *MockActivityTemplateProvider) CreateActivity(ctx context.Context, activity *entities.ActivityTemplate) (int64, error) {
	args := m.Called(ctx, activity)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockActivityTemplateProvider) ListActivities(ctx context.Context, orgID uuid.UUID, moment *entities.ClassMoment) ([]entities.ActivityTemplate, error) {
	args := m.Called(ctx, orgID, moment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.ActivityTemplate), args.Error(1)
}

func (m *MockActivityTemplateProvider) CountByMoment(ctx context.Context, orgID uuid.UUID, moment entities.ClassMoment) (int64, error) {
	args := m.Called(ctx, orgID, moment)
	return args.Get(0).(int64), args.Error(1)
}
