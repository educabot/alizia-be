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

func TestDeleteSubject_ValidationMissingOrg(t *testing.T) {
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewDeleteSubject(subjects)

	err := uc.Execute(context.Background(), admin.DeleteSubjectRequest{SubjectID: 1})
	assert.ErrorIs(t, err, providers.ErrValidation)
	subjects.AssertNotCalled(t, "CountSubjectDependencies", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteSubject_ValidationMissingID(t *testing.T) {
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewDeleteSubject(subjects)

	err := uc.Execute(context.Background(), admin.DeleteSubjectRequest{OrgID: uuid.New()})
	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestDeleteSubject_NotFound(t *testing.T) {
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewDeleteSubject(subjects)

	orgID := uuid.New()
	subjects.On("CountSubjectDependencies", mock.Anything, orgID, int64(99)).
		Return(providers.SubjectDependencies{}, providers.ErrNotFound)

	err := uc.Execute(context.Background(), admin.DeleteSubjectRequest{
		OrgID:     orgID,
		SubjectID: 99,
	})
	assert.ErrorIs(t, err, providers.ErrNotFound)
	subjects.AssertNotCalled(t, "DeleteSubject", mock.Anything, mock.Anything, mock.Anything)
}

// TestDeleteSubject_ConflictCourseSubjects is the main blocker case: the subject
// is still assigned to at least one course.
func TestDeleteSubject_ConflictCourseSubjects(t *testing.T) {
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewDeleteSubject(subjects)

	orgID := uuid.New()
	subjects.On("CountSubjectDependencies", mock.Anything, orgID, int64(1)).
		Return(providers.SubjectDependencies{CourseSubjects: 4}, nil)

	err := uc.Execute(context.Background(), admin.DeleteSubjectRequest{
		OrgID:     orgID,
		SubjectID: 1,
	})
	assert.ErrorIs(t, err, providers.ErrConflict)
	assert.Contains(t, err.Error(), "4 course-subjects")
	subjects.AssertNotCalled(t, "DeleteSubject", mock.Anything, mock.Anything, mock.Anything)
}

func TestDeleteSubject_CounterError(t *testing.T) {
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewDeleteSubject(subjects)

	orgID := uuid.New()
	boom := errors.New("db down")
	subjects.On("CountSubjectDependencies", mock.Anything, orgID, int64(1)).
		Return(providers.SubjectDependencies{}, boom)

	err := uc.Execute(context.Background(), admin.DeleteSubjectRequest{
		OrgID:     orgID,
		SubjectID: 1,
	})
	assert.ErrorIs(t, err, boom)
}

func TestDeleteSubject_HappyPath(t *testing.T) {
	subjects := new(mockproviders.MockSubjectProvider)
	uc := admin.NewDeleteSubject(subjects)

	orgID := uuid.New()
	subjects.On("CountSubjectDependencies", mock.Anything, orgID, int64(1)).
		Return(providers.SubjectDependencies{}, nil)
	subjects.On("DeleteSubject", mock.Anything, orgID, int64(1)).Return(nil)

	err := uc.Execute(context.Background(), admin.DeleteSubjectRequest{
		OrgID:     orgID,
		SubjectID: 1,
	})
	assert.NoError(t, err)
	subjects.AssertExpectations(t)
}
