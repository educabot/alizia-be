package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type AssignCourseSubjectRequest struct {
	OrgID      uuid.UUID
	CourseID   int64
	SubjectID  int64
	TeacherID  int64
	SchoolYear int
	StartDate  *time.Time
	EndDate    *time.Time
}

func (r AssignCourseSubjectRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.CourseID == 0 {
		return fmt.Errorf("%w: course_id is required", providers.ErrValidation)
	}
	if r.SubjectID == 0 {
		return fmt.Errorf("%w: subject_id is required", providers.ErrValidation)
	}
	if r.TeacherID == 0 {
		return fmt.Errorf("%w: teacher_id is required", providers.ErrValidation)
	}
	if r.SchoolYear == 0 {
		return fmt.Errorf("%w: school_year is required", providers.ErrValidation)
	}
	return nil
}

type AssignCourseSubject interface {
	Execute(ctx context.Context, req AssignCourseSubjectRequest) (*entities.CourseSubject, error)
}

type assignCourseSubjectImpl struct {
	courses        providers.CourseProvider
	users          providers.UserProvider
	courseSubjects providers.CourseSubjectProvider
}

func NewAssignCourseSubject(
	courses providers.CourseProvider,
	users providers.UserProvider,
	courseSubjects providers.CourseSubjectProvider,
) AssignCourseSubject {
	return &assignCourseSubjectImpl{
		courses:        courses,
		users:          users,
		courseSubjects: courseSubjects,
	}
}

func (uc *assignCourseSubjectImpl) Execute(ctx context.Context, req AssignCourseSubjectRequest) (*entities.CourseSubject, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Verify course belongs to the org
	if _, err := uc.courses.GetCourse(ctx, req.OrgID, req.CourseID); err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	// Verify teacher belongs to the org
	if _, err := uc.users.FindByID(ctx, req.OrgID, req.TeacherID); err != nil {
		return nil, fmt.Errorf("teacher not found: %w", err)
	}

	cs := &entities.CourseSubject{
		OrganizationID: req.OrgID,
		CourseID:       req.CourseID,
		SubjectID:      req.SubjectID,
		TeacherID:      req.TeacherID,
		SchoolYear:     req.SchoolYear,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
	}

	id, err := uc.courseSubjects.CreateCourseSubject(ctx, cs)
	if err != nil {
		return nil, err
	}

	// Reload through the repo so the response includes Subject and Teacher.
	// The FE contract requires both fields populated (see
	// docs/frontend-breaking-changes-dtos.md §7); returning the in-memory
	// struct here would silently omit them via `omitempty`.
	return uc.courseSubjects.GetCourseSubject(ctx, req.OrgID, id)
}
