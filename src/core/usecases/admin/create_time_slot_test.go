package admin_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestCreateTimeSlot_NormalClass(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	courses := new(mockproviders.MockCourseProvider)
	ts := new(mockproviders.MockTimeSlotProvider)
	uc := admin.NewCreateTimeSlot(orgs, courses, ts)

	orgID := uuid.New()
	ctx := context.Background()

	courses.On("GetCourse", ctx, orgID, int64(1)).Return(&entities.Course{ID: 1}, nil)
	ts.On("CreateTimeSlot", ctx, mock.AnythingOfType("*entities.TimeSlot")).Return(int64(1), nil)

	result, err := uc.Execute(ctx, admin.CreateTimeSlotRequest{
		OrgID:            orgID,
		CourseID:         1,
		DayOfWeek:        1,
		StartTime:        "08:00",
		EndTime:          "09:30",
		CourseSubjectIDs: []int64{1},
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Len(t, result.Subjects, 1)
	// orgs.FindByID should NOT be called for single subject
	orgs.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestCreateTimeSlot_SharedClass(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	courses := new(mockproviders.MockCourseProvider)
	ts := new(mockproviders.MockTimeSlotProvider)
	uc := admin.NewCreateTimeSlot(orgs, courses, ts)

	orgID := uuid.New()
	ctx := context.Background()

	courses.On("GetCourse", ctx, orgID, int64(1)).Return(&entities.Course{ID: 1}, nil)
	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID: orgID, Config: datatypes.JSON(`{"shared_classes_enabled": true}`),
	}, nil)
	ts.On("CreateTimeSlot", ctx, mock.AnythingOfType("*entities.TimeSlot")).Return(int64(2), nil)

	result, err := uc.Execute(ctx, admin.CreateTimeSlotRequest{
		OrgID:            orgID,
		CourseID:         1,
		DayOfWeek:        3,
		StartTime:        "08:00",
		EndTime:          "09:30",
		CourseSubjectIDs: []int64{1, 2},
	})

	assert.NoError(t, err)
	assert.Len(t, result.Subjects, 2)
	orgs.AssertExpectations(t)
}

func TestCreateTimeSlot_SharedDisabled(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	courses := new(mockproviders.MockCourseProvider)
	ts := new(mockproviders.MockTimeSlotProvider)
	uc := admin.NewCreateTimeSlot(orgs, courses, ts)

	orgID := uuid.New()
	ctx := context.Background()

	courses.On("GetCourse", ctx, orgID, int64(1)).Return(&entities.Course{ID: 1}, nil)
	orgs.On("FindByID", ctx, orgID).Return(&entities.Organization{
		ID: orgID, Config: datatypes.JSON(`{"shared_classes_enabled": false}`),
	}, nil)

	_, err := uc.Execute(ctx, admin.CreateTimeSlotRequest{
		OrgID:            orgID,
		CourseID:         1,
		DayOfWeek:        3,
		StartTime:        "08:00",
		EndTime:          "09:30",
		CourseSubjectIDs: []int64{1, 2},
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, providers.ErrValidation)
	assert.Contains(t, err.Error(), "shared classes are not enabled")
	ts.AssertNotCalled(t, "CreateTimeSlot", mock.Anything, mock.Anything)
}

func TestCreateTimeSlot_ValidationErrors(t *testing.T) {
	orgs := new(mockproviders.MockOrganizationProvider)
	courses := new(mockproviders.MockCourseProvider)
	ts := new(mockproviders.MockTimeSlotProvider)
	uc := admin.NewCreateTimeSlot(orgs, courses, ts)

	tests := []struct {
		name string
		req  admin.CreateTimeSlotRequest
	}{
		{"missing org_id", admin.CreateTimeSlotRequest{CourseID: 1, StartTime: "08:00", EndTime: "09:00", CourseSubjectIDs: []int64{1}}},
		{"missing course_id", admin.CreateTimeSlotRequest{OrgID: uuid.New(), StartTime: "08:00", EndTime: "09:00", CourseSubjectIDs: []int64{1}}},
		{"missing start_time", admin.CreateTimeSlotRequest{OrgID: uuid.New(), CourseID: 1, EndTime: "09:00", CourseSubjectIDs: []int64{1}}},
		{"missing end_time", admin.CreateTimeSlotRequest{OrgID: uuid.New(), CourseID: 1, StartTime: "08:00", CourseSubjectIDs: []int64{1}}},
		{"day_of_week negative", admin.CreateTimeSlotRequest{OrgID: uuid.New(), CourseID: 1, DayOfWeek: -1, StartTime: "08:00", EndTime: "09:00", CourseSubjectIDs: []int64{1}}},
		{"day_of_week too high", admin.CreateTimeSlotRequest{OrgID: uuid.New(), CourseID: 1, DayOfWeek: 7, StartTime: "08:00", EndTime: "09:00", CourseSubjectIDs: []int64{1}}},
		{"invalid start_time format", admin.CreateTimeSlotRequest{OrgID: uuid.New(), CourseID: 1, StartTime: "8am", EndTime: "09:00", CourseSubjectIDs: []int64{1}}},
		{"invalid end_time format", admin.CreateTimeSlotRequest{OrgID: uuid.New(), CourseID: 1, StartTime: "08:00", EndTime: "nine", CourseSubjectIDs: []int64{1}}},
		{"start_time after end_time", admin.CreateTimeSlotRequest{OrgID: uuid.New(), CourseID: 1, StartTime: "10:00", EndTime: "08:00", CourseSubjectIDs: []int64{1}}},
		{"start_time equals end_time", admin.CreateTimeSlotRequest{OrgID: uuid.New(), CourseID: 1, StartTime: "08:00", EndTime: "08:00", CourseSubjectIDs: []int64{1}}},
		{"no subjects", admin.CreateTimeSlotRequest{OrgID: uuid.New(), CourseID: 1, StartTime: "08:00", EndTime: "09:00"}},
		{"too many subjects", admin.CreateTimeSlotRequest{OrgID: uuid.New(), CourseID: 1, StartTime: "08:00", EndTime: "09:00", CourseSubjectIDs: []int64{1, 2, 3}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}
	courses.AssertNotCalled(t, "GetCourse", mock.Anything, mock.Anything, mock.Anything)
}
