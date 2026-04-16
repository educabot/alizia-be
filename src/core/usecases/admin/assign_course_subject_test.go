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

func TestAssignCourseSubject_Success(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	users := new(mockproviders.MockUserProvider)
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewAssignCourseSubject(courses, users, cs)

	orgID := uuid.New()
	ctx := context.Background()

	reloaded := &entities.CourseSubject{
		ID: 1, OrganizationID: orgID, CourseID: 1, SubjectID: 2, TeacherID: 3, SchoolYear: 2026,
		Subject: &entities.Subject{ID: 2, Name: "Algebra"},
		Teacher: &entities.User{ID: 3, FirstName: "Ada", LastName: "Lovelace"},
	}

	courses.On("GetCourse", ctx, orgID, int64(1)).Return(&entities.Course{ID: 1, OrganizationID: orgID}, nil)
	users.On("FindByID", ctx, orgID, int64(3)).Return(&entities.User{ID: 3}, nil)
	cs.On("CreateCourseSubject", ctx, mock.AnythingOfType("*entities.CourseSubject")).Return(int64(1), nil)
	cs.On("GetCourseSubject", ctx, orgID, int64(1)).Return(reloaded, nil)

	result, err := uc.Execute(ctx, admin.AssignCourseSubjectRequest{
		OrgID: orgID, CourseID: 1, SubjectID: 2, TeacherID: 3, SchoolYear: 2026,
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, int64(1), result.CourseID)
	assert.Equal(t, int64(2), result.SubjectID)
	assert.NotNil(t, result.Subject, "Subject must be preloaded so the FE can render result.subject.name")
	assert.Equal(t, "Algebra", result.Subject.Name)
	assert.NotNil(t, result.Teacher, "Teacher must be preloaded so the FE can render result.teacher.first_name")
	assert.Equal(t, "Ada", result.Teacher.FirstName)
	courses.AssertExpectations(t)
	users.AssertExpectations(t)
	cs.AssertExpectations(t)
}

func TestAssignCourseSubject_CourseNotFound(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	users := new(mockproviders.MockUserProvider)
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewAssignCourseSubject(courses, users, cs)

	orgID := uuid.New()
	ctx := context.Background()

	courses.On("GetCourse", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, admin.AssignCourseSubjectRequest{
		OrgID: orgID, CourseID: 99, SubjectID: 1, TeacherID: 2, SchoolYear: 2026,
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestAssignCourseSubject_TeacherNotFound(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	users := new(mockproviders.MockUserProvider)
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewAssignCourseSubject(courses, users, cs)

	orgID := uuid.New()
	ctx := context.Background()

	courses.On("GetCourse", ctx, orgID, int64(1)).Return(&entities.Course{ID: 1}, nil)
	users.On("FindByID", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, admin.AssignCourseSubjectRequest{
		OrgID: orgID, CourseID: 1, SubjectID: 1, TeacherID: 99, SchoolYear: 2026,
	})

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestAssignCourseSubject_Duplicate(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	users := new(mockproviders.MockUserProvider)
	csMock := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewAssignCourseSubject(courses, users, csMock)

	orgID := uuid.New()
	ctx := context.Background()

	courses.On("GetCourse", ctx, orgID, int64(1)).Return(&entities.Course{ID: 1}, nil)
	users.On("FindByID", ctx, orgID, int64(3)).Return(&entities.User{ID: 3}, nil)
	csMock.On("CreateCourseSubject", ctx, mock.AnythingOfType("*entities.CourseSubject")).Return(int64(0), providers.ErrConflict)

	_, err := uc.Execute(ctx, admin.AssignCourseSubjectRequest{
		OrgID: orgID, CourseID: 1, SubjectID: 1, TeacherID: 3, SchoolYear: 2026,
	})

	assert.ErrorIs(t, err, providers.ErrConflict)
}

func TestAssignCourseSubject_ValidationErrors(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	users := new(mockproviders.MockUserProvider)
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewAssignCourseSubject(courses, users, cs)

	tests := []struct {
		name string
		req  admin.AssignCourseSubjectRequest
	}{
		{"missing org_id", admin.AssignCourseSubjectRequest{CourseID: 1, SubjectID: 1, TeacherID: 1, SchoolYear: 2026}},
		{"missing course_id", admin.AssignCourseSubjectRequest{OrgID: uuid.New(), SubjectID: 1, TeacherID: 1, SchoolYear: 2026}},
		{"missing subject_id", admin.AssignCourseSubjectRequest{OrgID: uuid.New(), CourseID: 1, TeacherID: 1, SchoolYear: 2026}},
		{"missing teacher_id", admin.AssignCourseSubjectRequest{OrgID: uuid.New(), CourseID: 1, SubjectID: 1, SchoolYear: 2026}},
		{"missing school_year", admin.AssignCourseSubjectRequest{OrgID: uuid.New(), CourseID: 1, SubjectID: 1, TeacherID: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}
	courses.AssertNotCalled(t, "GetCourse", mock.Anything, mock.Anything, mock.Anything)
}
