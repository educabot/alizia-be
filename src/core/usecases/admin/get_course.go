package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type GetCourseRequest struct {
	OrgID    uuid.UUID
	CourseID int64
}

func (r GetCourseRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.CourseID == 0 {
		return fmt.Errorf("%w: course_id is required", providers.ErrValidation)
	}
	return nil
}

type GetCourse interface {
	Execute(ctx context.Context, req GetCourseRequest) (*entities.Course, error)
}

type getCourseImpl struct {
	courses providers.CourseProvider
}

func NewGetCourse(courses providers.CourseProvider) GetCourse {
	return &getCourseImpl{courses: courses}
}

func (uc *getCourseImpl) Execute(ctx context.Context, req GetCourseRequest) (*entities.Course, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return uc.courses.GetCourse(ctx, req.OrgID, req.CourseID)
}
