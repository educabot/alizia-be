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

func TestUpdateCourse_ValidationMissingID(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewUpdateCourse(courses)

	name := "3ro A"
	_, err := uc.Execute(context.Background(), admin.UpdateCourseRequest{
		OrgID: uuid.New(),
		Name:  &name,
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
	courses.AssertNotCalled(t, "GetCourse", mock.Anything, mock.Anything, mock.Anything)
}

// TestUpdateCourse_ValidationNoFields asserts we reject empty PATCH bodies at
// validation time so the handler never issues a no-op UPDATE.
func TestUpdateCourse_ValidationNoFields(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewUpdateCourse(courses)

	_, err := uc.Execute(context.Background(), admin.UpdateCourseRequest{
		OrgID:    uuid.New(),
		CourseID: 1,
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestUpdateCourse_ValidationBlankName(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewUpdateCourse(courses)

	blank := "   "
	_, err := uc.Execute(context.Background(), admin.UpdateCourseRequest{
		OrgID:    uuid.New(),
		CourseID: 1,
		Name:     &blank,
	})
	assert.ErrorIs(t, err, providers.ErrValidation)
}

func TestUpdateCourse_NotFound(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewUpdateCourse(courses)

	orgID := uuid.New()
	name := "3ro A"
	courses.On("GetCourse", mock.Anything, orgID, int64(99)).Return(nil, providers.ErrNotFound)

	_, err := uc.Execute(context.Background(), admin.UpdateCourseRequest{
		OrgID:    orgID,
		CourseID: 99,
		Name:     &name,
	})
	assert.ErrorIs(t, err, providers.ErrNotFound)
	courses.AssertNotCalled(t, "UpdateCourse", mock.Anything, mock.Anything)
}

// TestUpdateCourse_HappyPath trims whitespace on the incoming name and returns
// the reloaded row (not the pre-update one) so the FE sees server-side state.
func TestUpdateCourse_HappyPath(t *testing.T) {
	courses := new(mockproviders.MockCourseProvider)
	uc := admin.NewUpdateCourse(courses)

	orgID := uuid.New()
	current := &entities.Course{ID: 1, OrganizationID: orgID, Name: "old"}
	reloaded := &entities.Course{ID: 1, OrganizationID: orgID, Name: "3ro A"}

	courses.On("GetCourse", mock.Anything, orgID, int64(1)).Return(current, nil).Once()
	courses.On("UpdateCourse", mock.Anything, mock.MatchedBy(func(c *entities.Course) bool {
		return c.ID == 1 && c.Name == "3ro A"
	})).Return(nil)
	courses.On("GetCourse", mock.Anything, orgID, int64(1)).Return(reloaded, nil).Once()

	name := "  3ro A  "
	got, err := uc.Execute(context.Background(), admin.UpdateCourseRequest{
		OrgID:    orgID,
		CourseID: 1,
		Name:     &name,
	})
	assert.NoError(t, err)
	assert.Equal(t, "3ro A", got.Name)
	courses.AssertExpectations(t)
}
