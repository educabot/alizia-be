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

func TestAddStudent_Success(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	students := new(mockproviders.MockStudentProvider)
	uc := admin.NewAddStudent(courses, students)

	orgID := uuid.New()
	ctx := context.Background()

	courses.On("GetCourse", ctx, orgID, int64(1)).Return(&entities.Course{ID: 1, OrganizationID: orgID}, nil)
	students.On("CreateStudent", ctx, mock.AnythingOfType("*entities.Student")).Return(int64(10), nil)

	result, err := uc.Execute(ctx, admin.AddStudentRequest{
		OrgID: orgID, CourseID: 1, Name: "Lucía Martinez",
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(10), result.ID)
	assert.Equal(t, "Lucía Martinez", result.Name)
	courses.AssertExpectations(t)
	students.AssertExpectations(t)
}

func TestAddStudent_CourseNotFound(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	students := new(mockproviders.MockStudentProvider)
	uc := admin.NewAddStudent(courses, students)

	orgID := uuid.New()
	ctx := context.Background()

	courses.On("GetCourse", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, admin.AddStudentRequest{
		OrgID: orgID, CourseID: 99, Name: "Test",
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
	students.AssertNotCalled(t, "CreateStudent", mock.Anything, mock.Anything)
}

func TestAddStudent_ValidationErrors(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	students := new(mockproviders.MockStudentProvider)
	uc := admin.NewAddStudent(courses, students)

	tests := []struct {
		name string
		req  admin.AddStudentRequest
	}{
		{"missing org_id", admin.AddStudentRequest{CourseID: 1, Name: "Test"}},
		{"missing course_id", admin.AddStudentRequest{OrgID: uuid.New(), Name: "Test"}},
		{"missing name", admin.AddStudentRequest{OrgID: uuid.New(), CourseID: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}
	courses.AssertNotCalled(t, "GetCourse", mock.Anything, mock.Anything, mock.Anything)
}
