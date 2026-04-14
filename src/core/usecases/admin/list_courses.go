package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type ListCoursesRequest struct {
	OrgID uuid.UUID
}

func (r ListCoursesRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	return nil
}

type ListCourses interface {
	Execute(ctx context.Context, req ListCoursesRequest) ([]entities.Course, error)
}

type listCoursesImpl struct {
	courses providers.CourseProvider
}

func NewListCourses(courses providers.CourseProvider) ListCourses {
	return &listCoursesImpl{courses: courses}
}

func (uc *listCoursesImpl) Execute(ctx context.Context, req ListCoursesRequest) ([]entities.Course, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return uc.courses.ListCourses(ctx, req.OrgID)
}
