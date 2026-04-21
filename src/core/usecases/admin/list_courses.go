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
	Page  providers.Pagination
}

func (r ListCoursesRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	return nil
}

type ListCoursesResponse struct {
	Items []entities.Course
	More  bool
}

type ListCourses interface {
	Execute(ctx context.Context, req ListCoursesRequest) (*ListCoursesResponse, error)
}

type listCoursesImpl struct {
	courses providers.CourseProvider
}

func NewListCourses(courses providers.CourseProvider) ListCourses {
	return &listCoursesImpl{courses: courses}
}

func (uc *listCoursesImpl) Execute(ctx context.Context, req ListCoursesRequest) (*ListCoursesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	items, more, err := uc.courses.ListCourses(ctx, req.OrgID, req.Page)
	if err != nil {
		return nil, err
	}
	return &ListCoursesResponse{Items: items, More: more}, nil
}
