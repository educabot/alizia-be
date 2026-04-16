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

func TestGetCourseSubject_Success(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewGetCourseSubject(cs)

	orgID := uuid.New()
	ctx := context.Background()

	expected := &entities.CourseSubject{
		ID: 7, OrganizationID: orgID, CourseID: 1, SubjectID: 2, TeacherID: 3,
		Subject: &entities.Subject{ID: 2, Name: "Algebra"},
		Teacher: &entities.User{ID: 3, FirstName: "Ada"},
	}
	cs.On("GetCourseSubject", ctx, orgID, int64(7)).Return(expected, nil)

	result, err := uc.Execute(ctx, admin.GetCourseSubjectRequest{OrgID: orgID, CourseSubjectID: 7})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	cs.AssertExpectations(t)
}

func TestGetCourseSubject_NotFound(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewGetCourseSubject(cs)

	orgID := uuid.New()
	ctx := context.Background()

	cs.On("GetCourseSubject", ctx, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(ctx, admin.GetCourseSubjectRequest{OrgID: orgID, CourseSubjectID: 99})

	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestGetCourseSubject_ValidationErrors(t *testing.T) {
	cs := new(mockproviders.MockCourseSubjectProvider)
	uc := admin.NewGetCourseSubject(cs)

	tests := []struct {
		name string
		req  admin.GetCourseSubjectRequest
	}{
		{"missing org_id", admin.GetCourseSubjectRequest{CourseSubjectID: 1}},
		{"missing id", admin.GetCourseSubjectRequest{OrgID: uuid.New()}},
		{"negative id", admin.GetCourseSubjectRequest{OrgID: uuid.New(), CourseSubjectID: -1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}

	cs.AssertNotCalled(t, "GetCourseSubject", mock.Anything, mock.Anything, mock.Anything)
}
