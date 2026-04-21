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

func TestListCourseSubjects_SuccessNoFilters(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewListCourseSubjects(cs)

	orgID := uuid.New()
	ctx := context.Background()

	expected := []entities.CourseSubject{
		{ID: 1, OrganizationID: orgID, CourseID: 10, SubjectID: 100, TeacherID: 1000},
		{ID: 2, OrganizationID: orgID, CourseID: 11, SubjectID: 101, TeacherID: 1001},
	}
	cs.On("ListCourseSubjects", ctx, orgID, providers.CourseSubjectFilter{}).Return(expected, nil)

	result, err := uc.Execute(ctx, admin.ListCourseSubjectsRequest{OrgID: orgID})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(1), result[0].ID)
	cs.AssertExpectations(t)
}

func TestListCourseSubjects_SuccessFilterByTeacher(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewListCourseSubjects(cs)

	orgID := uuid.New()
	ctx := context.Background()
	tid := int64(42)

	expected := []entities.CourseSubject{
		{ID: 5, OrganizationID: orgID, CourseID: 10, SubjectID: 100, TeacherID: tid},
	}
	cs.On("ListCourseSubjects", ctx, orgID, providers.CourseSubjectFilter{TeacherID: &tid}).Return(expected, nil)

	result, err := uc.Execute(ctx, admin.ListCourseSubjectsRequest{OrgID: orgID, TeacherID: &tid})

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, tid, result[0].TeacherID)
	cs.AssertExpectations(t)
}

func TestListCourseSubjects_SuccessFilterByCourseAndTeacher(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewListCourseSubjects(cs)

	orgID := uuid.New()
	ctx := context.Background()
	cid := int64(10)
	tid := int64(42)

	expected := []entities.CourseSubject{
		{ID: 7, OrganizationID: orgID, CourseID: cid, SubjectID: 100, TeacherID: tid},
	}
	cs.On("ListCourseSubjects", ctx, orgID, providers.CourseSubjectFilter{CourseID: &cid, TeacherID: &tid}).Return(expected, nil)

	result, err := uc.Execute(ctx, admin.ListCourseSubjectsRequest{OrgID: orgID, CourseID: &cid, TeacherID: &tid})

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, cid, result[0].CourseID)
	assert.Equal(t, tid, result[0].TeacherID)
	cs.AssertExpectations(t)
}

func TestListCourseSubjects_ValidationMissingOrgID(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewListCourseSubjects(cs)

	_, err := uc.Execute(context.Background(), admin.ListCourseSubjectsRequest{})

	assert.ErrorIs(t, err, providers.ErrValidation)
	cs.AssertNotCalled(t, "ListCourseSubjects", mock.Anything, mock.Anything, mock.Anything)
}
