package providers

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
)

type MockTimeSlotProvider struct {
	mock.Mock
}

func (m *MockTimeSlotProvider) CreateTimeSlot(ctx context.Context, slot *entities.TimeSlot) (int64, error) {
	args := m.Called(ctx, slot)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTimeSlotProvider) ListByCourse(ctx context.Context, courseID int64) ([]entities.TimeSlot, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.TimeSlot), args.Error(1)
}

func (m *MockTimeSlotProvider) GetSharedClassNumbers(ctx context.Context, courseSubjectID int64, totalClasses int) ([]int, error) {
	args := m.Called(ctx, courseSubjectID, totalClasses)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}
