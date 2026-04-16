package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

// ListAllSubjectsRequest lists subjects across the org. AreaID is optional;
// when supplied, results are filtered to that area (which must belong to the org).
type ListAllSubjectsRequest struct {
	OrgID  uuid.UUID
	AreaID *int64
}

func (r ListAllSubjectsRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.AreaID != nil && *r.AreaID <= 0 {
		return fmt.Errorf("%w: area_id must be > 0", providers.ErrValidation)
	}
	return nil
}

type ListAllSubjects interface {
	Execute(ctx context.Context, req ListAllSubjectsRequest) ([]entities.Subject, error)
}

type listAllSubjectsImpl struct {
	areas    providers.AreaProvider
	subjects providers.SubjectProvider
}

func NewListAllSubjects(areas providers.AreaProvider, subjects providers.SubjectProvider) ListAllSubjects {
	return &listAllSubjectsImpl{areas: areas, subjects: subjects}
}

func (uc *listAllSubjectsImpl) Execute(ctx context.Context, req ListAllSubjectsRequest) ([]entities.Subject, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// When an area filter is supplied, verify the area belongs to the org.
	if req.AreaID != nil && *req.AreaID > 0 {
		if _, err := uc.areas.GetArea(ctx, req.OrgID, *req.AreaID); err != nil {
			return nil, fmt.Errorf("area not found: %w", err)
		}
	}

	return uc.subjects.ListSubjectsByOrg(ctx, req.OrgID, req.AreaID)
}
