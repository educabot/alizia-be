package admin_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func mustDate(t *testing.T, s string) time.Time {
	t.Helper()
	out, err := time.Parse("2006-01-02", s)
	require.NoError(t, err)
	return out
}

func TestUpdateCourseSubject_ValidationMissingID(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	users := new(mockproviders.MockUserProvider)
	uc := admin.NewUpdateCourseSubject(cs, users)

	teacher := int64(9)
	_, err := uc.Execute(context.Background(), admin.UpdateCourseSubjectRequest{
		OrgID:     uuid.New(),
		TeacherID: &teacher,
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
	cs.AssertNotCalled(t, "GetCourseSubject", mock.Anything, mock.Anything, mock.Anything)
}

// TestUpdateCourseSubject_ValidationNoFields asserts we reject empty PATCH
// bodies at validation time — the handler should never issue a no-op UPDATE.
func TestUpdateCourseSubject_ValidationNoFields(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	users := new(mockproviders.MockUserProvider)
	uc := admin.NewUpdateCourseSubject(cs, users)

	_, err := uc.Execute(context.Background(), admin.UpdateCourseSubjectRequest{
		OrgID:           uuid.New(),
		CourseSubjectID: 1,
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
	cs.AssertNotCalled(t, "GetCourseSubject", mock.Anything, mock.Anything, mock.Anything)
}

func TestUpdateCourseSubject_NotFound(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	users := new(mockproviders.MockUserProvider)
	uc := admin.NewUpdateCourseSubject(cs, users)

	orgID := uuid.New()
	teacher := int64(9)
	cs.On("GetCourseSubject", mock.Anything, orgID, int64(99)).
		Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(context.Background(), admin.UpdateCourseSubjectRequest{
		OrgID:           orgID,
		CourseSubjectID: 99,
		TeacherID:       &teacher,
	})
	assert.ErrorIs(t, err, providers.ErrNotFound)
	users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything, mock.Anything)
	cs.AssertNotCalled(t, "UpdateCourseSubject", mock.Anything, mock.Anything)
}

// TestUpdateCourseSubject_TeacherMissing verifies tenant verification: if the
// teacher doesn't belong to the org, we must not write the UPDATE.
func TestUpdateCourseSubject_TeacherMissing(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	users := new(mockproviders.MockUserProvider)
	uc := admin.NewUpdateCourseSubject(cs, users)

	orgID := uuid.New()
	current := &entities.CourseSubject{ID: 1, OrganizationID: orgID, TeacherID: 10}
	cs.On("GetCourseSubject", mock.Anything, orgID, int64(1)).Return(current, nil)

	newTeacher := int64(42)
	users.On("FindByID", mock.Anything, orgID, newTeacher).
		Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(context.Background(), admin.UpdateCourseSubjectRequest{
		OrgID:           orgID,
		CourseSubjectID: 1,
		TeacherID:       &newTeacher,
	})
	assert.ErrorIs(t, err, providers.ErrNotFound)
	cs.AssertNotCalled(t, "UpdateCourseSubject", mock.Anything, mock.Anything)
}

// TestUpdateCourseSubject_SameTeacherSkipsLookup checks that when the teacher
// isn't actually changing we don't waste a FindByID roundtrip.
func TestUpdateCourseSubject_SameTeacherSkipsLookup(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	users := new(mockproviders.MockUserProvider)
	uc := admin.NewUpdateCourseSubject(cs, users)

	orgID := uuid.New()
	current := &entities.CourseSubject{ID: 1, OrganizationID: orgID, TeacherID: 10, SchoolYear: 2026}
	reloaded := *current
	reloaded.SchoolYear = 2027

	same := int64(10)
	year := 2027
	cs.On("GetCourseSubject", mock.Anything, orgID, int64(1)).Return(current, nil).Once()
	cs.On("UpdateCourseSubject", mock.Anything, mock.MatchedBy(func(e *entities.CourseSubject) bool {
		return e.ID == 1 && e.TeacherID == 10 && e.SchoolYear == 2027
	})).Return(nil)
	cs.On("GetCourseSubject", mock.Anything, orgID, int64(1)).Return(&reloaded, nil).Once()

	got, err := uc.Execute(context.Background(), admin.UpdateCourseSubjectRequest{
		OrgID:           orgID,
		CourseSubjectID: 1,
		TeacherID:       &same,
		SchoolYear:      &year,
	})
	assert.NoError(t, err)
	assert.Equal(t, 2027, got.SchoolYear)
	users.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything, mock.Anything)
	cs.AssertExpectations(t)
}

// TestUpdateCourseSubject_DateOrder protects the invariant that start_date
// precedes end_date so we never write an upside-down range to the DB.
func TestUpdateCourseSubject_DateOrder(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	users := new(mockproviders.MockUserProvider)
	uc := admin.NewUpdateCourseSubject(cs, users)

	orgID := uuid.New()
	current := &entities.CourseSubject{ID: 1, OrganizationID: orgID, TeacherID: 10}
	cs.On("GetCourseSubject", mock.Anything, orgID, int64(1)).Return(current, nil)

	start := mustDate(t, "2026-12-31")
	end := mustDate(t, "2026-01-01")

	_, err := uc.Execute(context.Background(), admin.UpdateCourseSubjectRequest{
		OrgID:           orgID,
		CourseSubjectID: 1,
		StartDate:       &start,
		EndDate:         &end,
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
	cs.AssertNotCalled(t, "UpdateCourseSubject", mock.Anything, mock.Anything)
}

// TestUpdateCourseSubject_HappyPath covers the main reassign-teacher flow the
// FE hits most often: verify teacher exists → persist → return reloaded row.
func TestUpdateCourseSubject_HappyPath(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	users := new(mockproviders.MockUserProvider)
	uc := admin.NewUpdateCourseSubject(cs, users)

	orgID := uuid.New()
	current := &entities.CourseSubject{ID: 1, OrganizationID: orgID, TeacherID: 10}
	reloaded := *current
	reloaded.TeacherID = 42
	reloaded.Teacher = &entities.User{ID: 42, FirstName: "Ana"}

	newTeacher := int64(42)
	cs.On("GetCourseSubject", mock.Anything, orgID, int64(1)).Return(current, nil).Once()
	users.On("FindByID", mock.Anything, orgID, newTeacher).
		Return(&entities.User{ID: 42}, nil)
	cs.On("UpdateCourseSubject", mock.Anything, mock.MatchedBy(func(e *entities.CourseSubject) bool {
		return e.ID == 1 && e.TeacherID == 42
	})).Return(nil)
	cs.On("GetCourseSubject", mock.Anything, orgID, int64(1)).Return(&reloaded, nil).Once()

	got, err := uc.Execute(context.Background(), admin.UpdateCourseSubjectRequest{
		OrgID:           orgID,
		CourseSubjectID: 1,
		TeacherID:       &newTeacher,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(42), got.TeacherID)
	assert.NotNil(t, got.Teacher)
	cs.AssertExpectations(t)
	users.AssertExpectations(t)
}

// TestUpdateCourseSubject_RepoConflict asserts we surface the conflict error
// verbatim so the HTTP layer can translate it to 409 without double-wrapping.
func TestUpdateCourseSubject_RepoConflict(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	users := new(mockproviders.MockUserProvider)
	uc := admin.NewUpdateCourseSubject(cs, users)

	orgID := uuid.New()
	current := &entities.CourseSubject{ID: 1, OrganizationID: orgID, TeacherID: 10, SchoolYear: 2026}
	cs.On("GetCourseSubject", mock.Anything, orgID, int64(1)).Return(current, nil).Once()
	year := 2027
	cs.On("UpdateCourseSubject", mock.Anything, mock.Anything).
		Return(providers.ErrConflict)

	_, err := uc.Execute(context.Background(), admin.UpdateCourseSubjectRequest{
		OrgID:           orgID,
		CourseSubjectID: 1,
		SchoolYear:      &year,
	})
	assert.ErrorIs(t, err, providers.ErrConflict)
}
