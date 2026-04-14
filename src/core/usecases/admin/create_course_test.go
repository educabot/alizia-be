package admin_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/core/usecases/admin"
	mockproviders "github.com/educabot/alizia-be/src/mocks/providers"
)

func TestCreateCourse_Success(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewCreateCourse(courses)

	orgID := uuid.New()
	ctx := context.Background()

	courses.On("CreateCourse", ctx, mock.AnythingOfType("*entities.Course")).Return(int64(1), nil)

	result, err := uc.Execute(ctx, admin.CreateCourseRequest{
		OrgID: orgID, Name: "2do 1era", Year: 2026,
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "2do 1era", result.Name)
	assert.Equal(t, 2026, result.Year)
	courses.AssertExpectations(t)
}

func TestCreateCourse_ValidationErrors(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewCreateCourse(courses)

	tests := []struct {
		name string
		req  admin.CreateCourseRequest
	}{
		{"missing org_id", admin.CreateCourseRequest{Name: "Test", Year: 2026}},
		{"missing name", admin.CreateCourseRequest{OrgID: uuid.New(), Year: 2026}},
		{"missing year", admin.CreateCourseRequest{OrgID: uuid.New(), Name: "Test"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.req)
			assert.ErrorIs(t, err, providers.ErrValidation)
		})
	}
	courses.AssertNotCalled(t, "CreateCourse", mock.Anything, mock.Anything)
}
