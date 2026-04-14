package admin_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestGetSchedule_Success(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	ts := new(mockproviders.MockTimeSlotProvider)
	uc := admin.NewGetSchedule(courses, ts)

	orgID := uuid.New()
	ctx := context.Background()

	courses.On("GetCourse", ctx, orgID, int64(1)).Return(&entities.Course{ID: 1}, nil)
	expected := []entities.TimeSlot{
		{ID: 1, CourseID: 1, DayOfWeek: 1, StartTime: "08:00", EndTime: "09:30"},
		{ID: 2, CourseID: 1, DayOfWeek: 1, StartTime: "09:45", EndTime: "11:15"},
	}
	ts.On("ListByCourse", ctx, int64(1)).Return(expected, nil)

	result, err := uc.Execute(ctx, admin.GetScheduleRequest{OrgID: orgID, CourseID: 1})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	courses.AssertExpectations(t)
	ts.AssertExpectations(t)
}

func TestGetSchedule_CourseNotFound(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	ts := new(mockproviders.MockTimeSlotProvider)
	uc := admin.NewGetSchedule(courses, ts)

	orgID := uuid.New()
	ctx := context.Background()

	courses.On("GetCourse", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, admin.GetScheduleRequest{OrgID: orgID, CourseID: 99})
	assert.ErrorIs(t, err, providers.ErrNotFound)
	ts.AssertNotCalled(t, "ListByCourse", mock.Anything, mock.Anything)
}

func TestGetSchedule_ValidationErrors(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	ts := new(mockproviders.MockTimeSlotProvider)
	uc := admin.NewGetSchedule(courses, ts)

	tests := []struct {
		name string
		req  admin.GetScheduleRequest
	}{
		{"missing org_id", admin.GetScheduleRequest{CourseID: 1}},
		{"missing course_id", admin.GetScheduleRequest{OrgID: uuid.New()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}
	courses.AssertNotCalled(t, "GetCourse", mock.Anything, mock.Anything, mock.Anything)
}
