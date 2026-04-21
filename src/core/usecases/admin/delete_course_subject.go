package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/providers"
)

type DeleteCourseSubjectRequest struct {
	OrgID           uuid.UUID
	CourseSubjectID int64
}

func (r DeleteCourseSubjectRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.CourseSubjectID <= 0 {
		return fmt.Errorf("%w: course_subject_id is required", providers.ErrValidation)
	}
	return nil
}

type DeleteCourseSubject interface {
	Execute(ctx context.Context, req DeleteCourseSubjectRequest) error
}

type deleteCourseSubjectImpl struct {
	cs providers.CourseSubjectProvider
}

func NewDeleteCourseSubject(cs providers.CourseSubjectProvider) DeleteCourseSubject {
	return &deleteCourseSubjectImpl{cs: cs}
}

// Execute deletes a course-subject after verifying nothing schedulable still
// points at it. Mirrors DeleteArea: we count dependencies and refuse with a
// 409 rather than relying on ON DELETE CASCADE — an admin who deletes a
// course-subject by mistake should not silently lose the course's timetable.
func (uc *deleteCourseSubjectImpl) Execute(ctx context.Context, req DeleteCourseSubjectRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	deps, err := uc.cs.CountCourseSubjectDependencies(ctx, req.OrgID, req.CourseSubjectID)
	if err != nil {
		return err
	}
	if !deps.IsEmpty() {
		return fmt.Errorf("%w: course-subject has dependencies (%d time-slot assignments); remove them before deleting",
			providers.ErrConflict, deps.TimeSlotSubjects)
	}

	return uc.cs.DeleteCourseSubject(ctx, req.OrgID, req.CourseSubjectID)
}
