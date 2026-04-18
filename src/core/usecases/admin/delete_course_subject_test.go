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

func TestDeleteCourseSubject_ValidationMissingOrg(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewDeleteCourseSubject(cs)

	err := uc.Execute(context.Background(), admin.DeleteCourseSubjectRequest{CourseSubjectID: 1})
	assert.ErrorIs(t, err, providers.ErrValidation)
	cs.AssertNotCalled(t, "CountCourseSubjectDependencies", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteCourseSubject_ValidationMissingID(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewDeleteCourseSubject(cs)

	err := uc.Execute(context.Background(), admin.DeleteCourseSubjectRequest{OrgID: uuid.New()})
	assert.ErrorIs(t, err, providers.ErrValidation)
}

// TestDeleteCourseSubject_NotFound surfaces ErrNotFound from the counter so
// the HTTP layer returns 404 — we rely on the repo's tenant check here rather
// than a separate Get call.
func TestDeleteCourseSubject_NotFound(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewDeleteCourseSubject(cs)

	orgID := uuid.New()
	cs.On("CountCourseSubjectDependencies", mock.Anything, orgID, int64(99)).
		Return(providers.CourseSubjectDependencies{}, providers.ErrNotFound)

	err := uc.Execute(context.Background(), admin.DeleteCourseSubjectRequest{
		OrgID:           orgID,
		CourseSubjectID: 99,
	})
	assert.ErrorIs(t, err, providers.ErrNotFound)
	cs.AssertNotCalled(t, "DeleteCourseSubject", mock.Anything, mock.Anything, mock.Anything)
}

// TestDeleteCourseSubject_Conflict asserts we refuse the delete when the
// course-subject is still scheduled — we want a 409 instead of the DB silently
// cascading through time_slot_subjects.
func TestDeleteCourseSubject_Conflict(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewDeleteCourseSubject(cs)

	orgID := uuid.New()
	cs.On("CountCourseSubjectDependencies", mock.Anything, orgID, int64(1)).
		Return(providers.CourseSubjectDependencies{TimeSlotSubjects: 3}, nil)

	err := uc.Execute(context.Background(), admin.DeleteCourseSubjectRequest{
		OrgID:           orgID,
		CourseSubjectID: 1,
	})
	assert.ErrorIs(t, err, providers.ErrConflict)
	cs.AssertNotCalled(t, "DeleteCourseSubject", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteCourseSubject_CounterError(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewDeleteCourseSubject(cs)

	orgID := uuid.New()
	boom := errors.New("db down")
	cs.On("CountCourseSubjectDependencies", mock.Anything, orgID, int64(1)).
		Return(providers.CourseSubjectDependencies{}, boom)

	err := uc.Execute(context.Background(), admin.DeleteCourseSubjectRequest{
		OrgID:           orgID,
		CourseSubjectID: 1,
	})
	assert.ErrorIs(t, err, boom)
}

// TestDeleteCourseSubject_HappyPath covers the clean path: no schedule refs,
// delete returns nil.
func TestDeleteCourseSubject_HappyPath(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewDeleteCourseSubject(cs)

	orgID := uuid.New()
	cs.On("CountCourseSubjectDependencies", mock.Anything, orgID, int64(1)).
		Return(providers.CourseSubjectDependencies{}, nil)
	cs.On("DeleteCourseSubject", mock.Anything, orgID, int64(1)).Return(nil)

	err := uc.Execute(context.Background(), admin.DeleteCourseSubjectRequest{
		OrgID:           orgID,
		CourseSubjectID: 1,
	})
	assert.NoError(t, err)
	cs.AssertExpectations(t)
}
