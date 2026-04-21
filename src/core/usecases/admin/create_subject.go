package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type CreateSubjectRequest struct {
	OrgID       uuid.UUID
	AreaID      int64
	Name        string
	Description *string
}

func (r CreateSubjectRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.AreaID == 0 {
		return fmt.Errorf("%w: area_id is required", providers.ErrValidation)
	}
	if r.Name == "" {
		return fmt.Errorf("%w: name is required", providers.ErrValidation)
	}
	return nil
}

type CreateSubject interface {
	Execute(ctx context.Context, req CreateSubjectRequest) (*entities.Subject, error)
}

type createSubjectImpl struct {
	areas    providers.AreaProvider
	subjects providers.SubjectProvider
}

func NewCreateSubject(areas providers.AreaProvider, subjects providers.SubjectProvider) CreateSubject {
	return &createSubjectImpl{areas: areas, subjects: subjects}
}

func (uc *createSubjectImpl) Execute(ctx context.Context, req CreateSubjectRequest) (*entities.Subject, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Verify area belongs to the org
	if _, err := uc.areas.GetArea(ctx, req.OrgID, req.AreaID); err != nil {
		return nil, fmt.Errorf("area not found: %w", err)
	}

	subject := &entities.Subject{
		OrganizationID: req.OrgID,
		AreaID:         req.AreaID,
		Name:           req.Name,
		Description:    req.Description,
	}

	id, err := uc.subjects.CreateSubject(ctx, subject)
	if err != nil {
		return nil, err
	}

	subject.ID = id
	return subject, nil
}
