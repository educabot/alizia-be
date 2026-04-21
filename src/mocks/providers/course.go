package providers

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
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

func (m *MockCourseProvider) ListCourses(ctx context.Context, orgID uuid.UUID, p providers.Pagination) ([]entities.Course, bool, error) {
	args := m.Called(ctx, orgID, p)
	if args.Get(0) == nil {
		return nil, false, args.Error(2)
	}
	return args.Get(0).([]entities.Course), args.Bool(1), args.Error(2)
}

func (m *MockCourseProvider) UpdateCourse(ctx context.Context, course *entities.Course) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockCourseProvider) CountCourseDependencies(ctx context.Context, orgID uuid.UUID, id int64) (providers.CourseDependencies, error) {
	args := m.Called(ctx, orgID, id)
	return args.Get(0).(providers.CourseDependencies), args.Error(1)
}

func (m *MockCourseProvider) DeleteCourse(ctx context.Context, orgID uuid.UUID, id int64) error {
	args := m.Called(ctx, orgID, id)
	return args.Error(0)
}
