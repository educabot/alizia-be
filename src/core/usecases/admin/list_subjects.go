package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type ListSubjectsRequest struct {
	OrgID  uuid.UUID
	AreaID int64
}

func (r ListSubjectsRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.AreaID == 0 {
		return fmt.Errorf("%w: area_id is required", providers.ErrValidation)
	}
	return nil
}

type ListSubjects interface {
	Execute(ctx context.Context, req ListSubjectsRequest) ([]entities.Subject, error)
}

type listSubjectsImpl struct {
	areas    providers.AreaProvider
	subjects providers.SubjectProvider
}

func NewListSubjects(areas providers.AreaProvider, subjects providers.SubjectProvider) ListSubjects {
	return &listSubjectsImpl{areas: areas, subjects: subjects}
}

func (uc *listSubjectsImpl) Execute(ctx context.Context, req ListSubjectsRequest) ([]entities.Subject, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Verify area belongs to the org
	if _, err := uc.areas.GetArea(ctx, req.OrgID, req.AreaID); err != nil {
		return nil, fmt.Errorf("area not found: %w", err)
	}

	return uc.subjects.ListSubjectsByArea(ctx, req.OrgID, req.AreaID)
}
