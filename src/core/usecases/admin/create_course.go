package admin

import (
	"context"
	"fmt"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type CreateCourseRequest struct {
	OrgID int64
	Name  string
	Year  int
}

func (r CreateCourseRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("%w: name is required", providers.ErrValidation)
	}
	return nil
}

type CreateCourse interface {
	Execute(ctx context.Context, req CreateCourseRequest) (int64, error)
}

type createCourseImpl struct {
	courses providers.CourseProvider
}

func NewCreateCourse(courses providers.CourseProvider) CreateCourse {
	return &createCourseImpl{courses: courses}
}

func (uc *createCourseImpl) Execute(ctx context.Context, req CreateCourseRequest) (int64, error) {
	if err := req.Validate(); err != nil {
		return 0, err
	}

	course := &entities.Course{
		OrganizationID: req.OrgID,
		Name:           req.Name,
		Year:           req.Year,
	}

	return uc.courses.CreateCourse(ctx, course)
}
