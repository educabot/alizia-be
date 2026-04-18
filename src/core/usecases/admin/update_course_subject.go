package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type UpdateCourseSubjectRequest struct {
	OrgID           uuid.UUID
	CourseSubjectID int64
	// Each patch field is applied only when non-nil. This preserves "unset"
	// semantics for partial updates; a full PUT would lose the ability to
	// express "leave this alone".
	TeacherID  *int64
	StartDate  *time.Time
	EndDate    *time.Time
	SchoolYear *int
}

func (r UpdateCourseSubjectRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.CourseSubjectID <= 0 {
		return fmt.Errorf("%w: course_subject_id is required", providers.ErrValidation)
	}
	if r.TeacherID == nil && r.StartDate == nil && r.EndDate == nil && r.SchoolYear == nil {
		return fmt.Errorf("%w: at least one field must be provided", providers.ErrValidation)
	}
	if r.TeacherID != nil && *r.TeacherID <= 0 {
		return fmt.Errorf("%w: teacher_id must be positive", providers.ErrValidation)
	}
	if r.SchoolYear != nil && *r.SchoolYear <= 0 {
		return fmt.Errorf("%w: school_year must be positive", providers.ErrValidation)
	}
	return nil
}

type UpdateCourseSubject interface {
	Execute(ctx context.Context, req UpdateCourseSubjectRequest) (*entities.CourseSubject, error)
}

type updateCourseSubjectImpl struct {
	courseSubjects providers.CourseSubjectProvider
	users          providers.UserProvider
}

func NewUpdateCourseSubject(
	courseSubjects providers.CourseSubjectProvider,
	users providers.UserProvider,
) UpdateCourseSubject {
	return &updateCourseSubjectImpl{courseSubjects: courseSubjects, users: users}
}

// Execute loads the current course-subject, applies the non-nil patches, and
// persists the result. If TeacherID changes we verify the new teacher exists
// in the same tenant; the repo handles the unique_violation case for the
// (course_id, subject_id, school_year) composite key.
func (uc *updateCourseSubjectImpl) Execute(ctx context.Context, req UpdateCourseSubjectRequest) (*entities.CourseSubject, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	cs, err := uc.courseSubjects.GetCourseSubject(ctx, req.OrgID, req.CourseSubjectID)
	if err != nil {
		return nil, err
	}

	if req.TeacherID != nil && *req.TeacherID != cs.TeacherID {
		if _, err := uc.users.FindByID(ctx, req.OrgID, *req.TeacherID); err != nil {
			return nil, fmt.Errorf("teacher not found: %w", err)
		}
		cs.TeacherID = *req.TeacherID
	}
	if req.StartDate != nil {
		cs.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		cs.EndDate = req.EndDate
	}
	if req.SchoolYear != nil {
		cs.SchoolYear = *req.SchoolYear
	}

	if cs.StartDate != nil && cs.EndDate != nil && cs.StartDate.After(*cs.EndDate) {
		return nil, fmt.Errorf("%w: start_date must be on or before end_date", providers.ErrValidation)
	}

	if err := uc.courseSubjects.UpdateCourseSubject(ctx, cs); err != nil {
		return nil, err
	}

	return uc.courseSubjects.GetCourseSubject(ctx, req.OrgID, cs.ID)
}
