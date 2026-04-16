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

func TestGetCourse_Success(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewGetCourse(courses)

	orgID := uuid.New()
	ctx := context.Background()

	courses.On("GetCourse", ctx, orgID, int64(1)).Return(&entities.Course{
		ID: 1, OrganizationID: orgID, Name: "2do 1era",
		Students:       []entities.Student{{ID: 1, Name: "Lucía"}},
		CourseSubjects: []entities.CourseSubject{{ID: 1, SubjectID: 1, TeacherID: 2}},
	}, nil)

	result, err := uc.Execute(ctx, admin.GetCourseRequest{OrgID: orgID, CourseID: 1})

	assert.NoError(t, err)
	assert.Equal(t, "2do 1era", result.Name)
	assert.Len(t, result.Students, 1)
	assert.Len(t, result.CourseSubjects, 1)
	courses.AssertExpectations(t)
}

func TestGetCourse_NotFound(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewGetCourse(courses)

	orgID := uuid.New()
	ctx := context.Background()

	courses.On("GetCourse", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, admin.GetCourseRequest{OrgID: orgID, CourseID: 99})
	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestGetCourse_ValidationErrors(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewGetCourse(courses)

	tests := []struct {
		name string
		req  admin.GetCourseRequest
	}{
		{"missing org_id", admin.GetCourseRequest{CourseID: 1}},
		{"missing course_id", admin.GetCourseRequest{OrgID: uuid.New()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}
	courses.AssertNotCalled(t, "GetCourse", mock.Anything, mock.Anything, mock.Anything)
}
