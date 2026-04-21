package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/providers"
)

type DeleteSubjectRequest struct {
	OrgID     uuid.UUID
	SubjectID int64
}

func (r DeleteSubjectRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.SubjectID <= 0 {
		return fmt.Errorf("%w: subject_id is required", providers.ErrValidation)
	}
	return nil
}

type DeleteSubject interface {
	Execute(ctx context.Context, req DeleteSubjectRequest) error
}

type deleteSubjectImpl struct {
	subjects providers.SubjectProvider
}

func NewDeleteSubject(subjects providers.SubjectProvider) DeleteSubject {
	return &deleteSubjectImpl{subjects: subjects}
}

// Execute refuses the delete if any course_subject references the subject.
// The `subjects.id` FK on course_subjects has no ON DELETE action, so without
// this check the delete would surface as a raw 500 from PG.
func (uc *deleteSubjectImpl) Execute(ctx context.Context, req DeleteSubjectRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	deps, err := uc.subjects.CountSubjectDependencies(ctx, req.OrgID, req.SubjectID)
	if err != nil {
		return err
	}
	if !deps.IsEmpty() {
		return fmt.Errorf("%w: subject has dependencies (%d course-subjects); remove them before deleting",
			providers.ErrConflict, deps.CourseSubjects)
	}

	return uc.subjects.DeleteSubject(ctx, req.OrgID, req.SubjectID)
}
