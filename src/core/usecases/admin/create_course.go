package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type CreateCourseRequest struct {
	OrgID uuid.UUID
	Name  string
	Year  int
}

func (r CreateCourseRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.Name == "" {
		return fmt.Errorf("%w: name is required", providers.ErrValidation)
	}
	if r.Year == 0 {
		return fmt.Errorf("%w: year is required", providers.ErrValidation)
	}
	return nil
}

type CreateCourse interface {
	Execute(ctx context.Context, req CreateCourseRequest) (*entities.Course, error)
}

type createCourseImpl struct {
	courses providers.CourseProvider
}

func NewCreateCourse(courses providers.CourseProvider) CreateCourse {
	return &createCourseImpl{courses: courses}
}

func (uc *createCourseImpl) Execute(ctx context.Context, req CreateCourseRequest) (*entities.Course, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	course := &entities.Course{
		OrganizationID: req.OrgID,
		Name:           req.Name,
		Year:           req.Year,
	}

	id, err := uc.courses.CreateCourse(ctx, course)
	if err != nil {
		return nil, err
	}

	course.ID = id
	return course, nil
}
