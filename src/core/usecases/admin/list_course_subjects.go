package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

// ListCourseSubjectsRequest represents the input for listing course-subjects
// across an organization, optionally filtered by course, subject and/or teacher.
// All filters are optional and combined with AND semantics by the repository.
type ListCourseSubjectsRequest struct {
	OrgID     uuid.UUID
	CourseID  *int64
	SubjectID *int64
	TeacherID *int64
}

func (r ListCourseSubjectsRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	return nil
}

type ListCourseSubjects interface {
	Execute(ctx context.Context, req ListCourseSubjectsRequest) ([]entities.CourseSubject, error)
}

type listCourseSubjectsImpl struct {
	courseSubjects providers.CourseSubjectProvider
}

// NewListCourseSubjects builds the usecase. The repo enforces tenant scoping
// itself (filters by orgID via JOIN with courses), so no pre-validation of
// area/course/subject is needed here.
func NewListCourseSubjects(courseSubjects providers.CourseSubjectProvider) ListCourseSubjects {
	return &listCourseSubjectsImpl{courseSubjects: courseSubjects}
}

func (uc *listCourseSubjectsImpl) Execute(ctx context.Context, req ListCourseSubjectsRequest) ([]entities.CourseSubject, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	filter := providers.CourseSubjectFilter{
		CourseID:  req.CourseID,
		SubjectID: req.SubjectID,
		TeacherID: req.TeacherID,
	}

	return uc.courseSubjects.ListCourseSubjects(ctx, req.OrgID, filter)
}
