package providers

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
)

type MockCourseProvider struct {
	mock.Mock
}

func (m *MockCourseProvider) CreateCourse(ctx context.Context, course *entities.Course) (int64, error) {
	args := m.Called(ctx, course)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCourseProvider) GetCourse(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Course, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Course), args.Error(1)
}

func (m *MockCourseProvider) ListCourses(ctx context.Context, orgID uuid.UUID) ([]entities.Course, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Course), args.Error(1)
}
