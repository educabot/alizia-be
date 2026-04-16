package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type GetCourseSubjectRequest struct {
	OrgID           uuid.UUID
	CourseSubjectID int64
}

func (r GetCourseSubjectRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.CourseSubjectID <= 0 {
		return fmt.Errorf("%w: course_subject_id is required", providers.ErrValidation)
	}
	return nil
}

type GetCourseSubject interface {
	Execute(ctx context.Context, req GetCourseSubjectRequest) (*entities.CourseSubject, error)
}

type getCourseSubjectImpl struct {
	courseSubjects providers.CourseSubjectProvider
}

func NewGetCourseSubject(courseSubjects providers.CourseSubjectProvider) GetCourseSubject {
	return &getCourseSubjectImpl{courseSubjects: courseSubjects}
}

func (uc *getCourseSubjectImpl) Execute(ctx context.Context, req GetCourseSubjectRequest) (*entities.CourseSubject, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return uc.courseSubjects.GetCourseSubject(ctx, req.OrgID, req.CourseSubjectID)
}
