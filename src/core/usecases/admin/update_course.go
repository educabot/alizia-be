package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

// UpdateCourseRequest is intentionally narrow: the `courses` table only exposes
// `name` as a mutable column today (`year` was dropped in migration 000013 and
// moved to course_subjects). When more fields become editable, extend this
// struct alongside the repo — do not let the handler accept fields the repo
// silently ignores.
type UpdateCourseRequest struct {
	OrgID    uuid.UUID
	CourseID int64
	Name     *string
}

func (r UpdateCourseRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.CourseID <= 0 {
		return fmt.Errorf("%w: course_id is required", providers.ErrValidation)
	}
	if r.Name == nil {
		return fmt.Errorf("%w: at least one field must be provided", providers.ErrValidation)
	}
	if strings.TrimSpace(*r.Name) == "" {
		return fmt.Errorf("%w: name must not be blank", providers.ErrValidation)
	}
	return nil
}

type UpdateCourse interface {
	Execute(ctx context.Context, req UpdateCourseRequest) (*entities.Course, error)
}

type updateCourseImpl struct {
	courses providers.CourseProvider
}

func NewUpdateCourse(courses providers.CourseProvider) UpdateCourse {
	return &updateCourseImpl{courses: courses}
}

// Execute patches a course's mutable fields and returns the reloaded row.
// GetCourse preloads students and course-subjects so the FE can keep its local
// cache coherent after the PATCH without a follow-up fetch.
func (uc *updateCourseImpl) Execute(ctx context.Context, req UpdateCourseRequest) (*entities.Course, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	current, err := uc.courses.GetCourse(ctx, req.OrgID, req.CourseID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		current.Name = strings.TrimSpace(*req.Name)
	}

	if err := uc.courses.UpdateCourse(ctx, current); err != nil {
		return nil, err
	}

	return uc.courses.GetCourse(ctx, req.OrgID, req.CourseID)
}
