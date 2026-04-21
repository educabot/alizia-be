package admin_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestDeleteCourse_ValidationMissingOrg(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewDeleteCourse(courses)

	err := uc.Execute(context.Background(), admin.DeleteCourseRequest{CourseID: 1})
	assert.ErrorIs(t, err, providers.ErrValidation)
	courses.AssertNotCalled(t, "CountCourseDependencies", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteCourse_ValidationMissingID(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewDeleteCourse(courses)

	err := uc.Execute(context.Background(), admin.DeleteCourseRequest{OrgID: uuid.New()})
	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestDeleteCourse_NotFound(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewDeleteCourse(courses)

	orgID := uuid.New()
	courses.On("CountCourseDependencies", mock.Anything, orgID, int64(99)).
		Return(providers.CourseDependencies{}, providers.ErrNotFound)

	err := uc.Execute(context.Background(), admin.DeleteCourseRequest{
		OrgID:    orgID,
		CourseID: 99,
	})
	assert.ErrorIs(t, err, providers.ErrNotFound)
	courses.AssertNotCalled(t, "DeleteCourse", mock.Anything, mock.Anything, mock.Anything)
}

// TestDeleteCourse_ConflictCourseSubjects covers the most common blocker:
// an admin cannot wipe a course while it still has subject assignments.
func TestDeleteCourse_ConflictCourseSubjects(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewDeleteCourse(courses)

	orgID := uuid.New()
	courses.On("CountCourseDependencies", mock.Anything, orgID, int64(1)).
		Return(providers.CourseDependencies{CourseSubjects: 2}, nil)

	err := uc.Execute(context.Background(), admin.DeleteCourseRequest{
		OrgID:    orgID,
		CourseID: 1,
	})
	assert.ErrorIs(t, err, providers.ErrConflict)
	courses.AssertNotCalled(t, "DeleteCourse", mock.Anything, mock.Anything, mock.Anything)
}

// TestDeleteCourse_ConflictStudents guards against silently deleting student
// records via ON DELETE CASCADE when the admin expected an error.
func TestDeleteCourse_ConflictStudents(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewDeleteCourse(courses)

	orgID := uuid.New()
	courses.On("CountCourseDependencies", mock.Anything, orgID, int64(1)).
		Return(providers.CourseDependencies{Students: 25}, nil)

	err := uc.Execute(context.Background(), admin.DeleteCourseRequest{
		OrgID:    orgID,
		CourseID: 1,
	})
	assert.ErrorIs(t, err, providers.ErrConflict)
	courses.AssertNotCalled(t, "DeleteCourse", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteCourse_CounterError(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewDeleteCourse(courses)

	orgID := uuid.New()
	boom := errors.New("db down")
	courses.On("CountCourseDependencies", mock.Anything, orgID, int64(1)).
		Return(providers.CourseDependencies{}, boom)

	err := uc.Execute(context.Background(), admin.DeleteCourseRequest{
		OrgID:    orgID,
		CourseID: 1,
	})
	assert.ErrorIs(t, err, boom)
}

func TestDeleteCourse_HappyPath(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewDeleteCourse(courses)

	orgID := uuid.New()
	courses.On("CountCourseDependencies", mock.Anything, orgID, int64(1)).
		Return(providers.CourseDependencies{}, nil)
	courses.On("DeleteCourse", mock.Anything, orgID, int64(1)).Return(nil)

	err := uc.Execute(context.Background(), admin.DeleteCourseRequest{
		OrgID:    orgID,
		CourseID: 1,
	})
	assert.NoError(t, err)
	courses.AssertExpectations(t)
}
