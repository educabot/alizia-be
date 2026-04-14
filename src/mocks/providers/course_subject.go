package providers

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
)

type MockCourseSubjectProvider struct {
	mock.Mock
}

func (m *MockCourseSubjectProvider) CreateCourseSubject(ctx context.Context, cs *entities.CourseSubject) (int64, error) {
	args := m.Called(ctx, cs)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCourseSubjectProvider) ListByCourse(ctx context.Context, courseID int64) ([]entities.CourseSubject, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.CourseSubject), args.Error(1)
}
