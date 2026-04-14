package providers

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
)

type MockStudentProvider struct {
	mock.Mock
}

func (m *MockStudentProvider) CreateStudent(ctx context.Context, student *entities.Student) (int64, error) {
	args := m.Called(ctx, student)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStudentProvider) ListByCourse(ctx context.Context, courseID int64) ([]entities.Student, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Student), args.Error(1)
}
