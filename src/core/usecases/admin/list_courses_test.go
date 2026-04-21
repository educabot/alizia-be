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

func TestListCourses_Success(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewListCourses(courses)

	orgID := uuid.New()
	ctx := context.Background()

	expected := []entities.Course{
		{ID: 1, Name: "2do 1era"},
		{ID: 2, Name: "3ro 2da"},
	}
	courses.On("ListCourses", ctx, orgID, mock.AnythingOfType("providers.Pagination")).Return(expected, false, nil)

	result, err := uc.Execute(ctx, admin.ListCoursesRequest{OrgID: orgID})

	assert.NoError(t, err)
	assert.Len(t, result.Items, 2)
	assert.False(t, result.More)
	courses.AssertExpectations(t)
}

func TestListCourses_WithPagination(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewListCourses(courses)

	orgID := uuid.New()
	ctx := context.Background()

	expected := []entities.Course{
		{ID: 1, Name: "2do 1era"},
	}
	courses.On("ListCourses", ctx, orgID, mock.AnythingOfType("providers.Pagination")).Return(expected, true, nil)

	result, err := uc.Execute(ctx, admin.ListCoursesRequest{
		OrgID: orgID,
		Page:  providers.Pagination{Limit: 1, Offset: 0},
	})

	assert.NoError(t, err)
	assert.Len(t, result.Items, 1)
	assert.True(t, result.More)
	courses.AssertExpectations(t)
}

func TestListCourses_ValidationError(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewListCourses(courses)

	_, err := uc.Execute(context.Background(), admin.ListCoursesRequest{})
	assert.ErrorIs(t, err, providers.ErrValidation)
	courses.AssertNotCalled(t, "ListCourses", mock.Anything, mock.Anything, mock.Anything)
}
