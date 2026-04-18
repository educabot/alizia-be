package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

// UpdateSubjectRequest patches a subject. Fields follow the same "nil means
// leave alone" convention used by UpdateCourseSubject; SetDescription gates the
// null-clearing case that can't be expressed with a pointer alone.
type UpdateSubjectRequest struct {
	OrgID          uuid.UUID
	SubjectID      int64
	Name           *string
	AreaID         *int64
	Description    *string
	SetDescription bool
}

func (r UpdateSubjectRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.SubjectID <= 0 {
		return fmt.Errorf("%w: subject_id is required", providers.ErrValidation)
	}
	if r.Name == nil && r.AreaID == nil && !r.SetDescription {
		return fmt.Errorf("%w: at least one field must be provided", providers.ErrValidation)
	}
	if r.Name != nil && strings.TrimSpace(*r.Name) == "" {
		return fmt.Errorf("%w: name must not be blank", providers.ErrValidation)
	}
	if r.AreaID != nil && *r.AreaID <= 0 {
		return fmt.Errorf("%w: area_id must be positive", providers.ErrValidation)
	}
	return nil
}

type UpdateSubject interface {
	Execute(ctx context.Context, req UpdateSubjectRequest) (*entities.Subject, error)
}

type updateSubjectImpl struct {
	areas    providers.AreaProvider
	subjects providers.SubjectProvider
}

func NewUpdateSubject(areas providers.AreaProvider, subjects providers.SubjectProvider) UpdateSubject {
	return &updateSubjectImpl{areas: areas, subjects: subjects}
}

// Execute loads the subject, applies the non-nil patches, and persists. If
// AreaID changes we verify the destination area belongs to the same tenant so
// admins can't move a subject across orgs by crafting a payload.
func (uc *updateSubjectImpl) Execute(ctx context.Context, req UpdateSubjectRequest) (*entities.Subject, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	current, err := uc.subjects.GetSubject(ctx, req.OrgID, req.SubjectID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		current.Name = strings.TrimSpace(*req.Name)
	}
	if req.AreaID != nil && *req.AreaID != current.AreaID {
		if _, err := uc.areas.GetArea(ctx, req.OrgID, *req.AreaID); err != nil {
			return nil, fmt.Errorf("area not found: %w", err)
		}
		current.AreaID = *req.AreaID
	}
	if req.SetDescription {
		current.Description = req.Description
	}

	if err := uc.subjects.UpdateSubject(ctx, current); err != nil {
		return nil, err
	}
	return current, nil
}
