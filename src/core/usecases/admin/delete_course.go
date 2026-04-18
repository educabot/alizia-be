package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/providers"
)

type DeleteCourseRequest struct {
	OrgID    uuid.UUID
	CourseID int64
}

func (r DeleteCourseRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.CourseID <= 0 {
		return fmt.Errorf("%w: course_id is required", providers.ErrValidation)
	}
	return nil
}

type DeleteCourse interface {
	Execute(ctx context.Context, req DeleteCourseRequest) error
}

type deleteCourseImpl struct {
	courses providers.CourseProvider
}

func NewDeleteCourse(courses providers.CourseProvider) DeleteCourse {
	return &deleteCourseImpl{courses: courses}
}

// Execute deletes a course after verifying it has no blocking dependencies.
// Returns ErrNotFound if the course doesn't belong to the caller's org, or
// ErrConflict listing the blocking counts if course-subjects, students or
// time-slots still reference it. We refuse to cascade because those deletions
// are destructive and not recoverable via the API even though the DB has ON
// DELETE CASCADE at the schema level.
func (uc *deleteCourseImpl) Execute(ctx context.Context, req DeleteCourseRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	deps, err := uc.courses.CountCourseDependencies(ctx, req.OrgID, req.CourseID)
	if err != nil {
		return err
	}
	if !deps.IsEmpty() {
		return fmt.Errorf("%w: course has dependencies (%d course-subjects, %d students, %d time-slots); remove them before deleting",
			providers.ErrConflict, deps.CourseSubjects, deps.Students, deps.TimeSlots)
	}

	return uc.courses.DeleteCourse(ctx, req.OrgID, req.CourseID)
}
