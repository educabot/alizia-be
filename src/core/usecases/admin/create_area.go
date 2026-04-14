package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type CreateAreaRequest struct {
	OrgID       uuid.UUID
	Name        string
	Description *string
}

func (r CreateAreaRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.Name == "" {
		return fmt.Errorf("%w: name is required", providers.ErrValidation)
	}
	return nil
}

type CreateArea interface {
	Execute(ctx context.Context, req CreateAreaRequest) (*entities.Area, error)
}

type createAreaImpl struct {
	areas providers.AreaProvider
}

func NewCreateArea(areas providers.AreaProvider) CreateArea {
	return &createAreaImpl{areas: areas}
}

func (uc *createAreaImpl) Execute(ctx context.Context, req CreateAreaRequest) (*entities.Area, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	area := &entities.Area{
		OrganizationID: req.OrgID,
		Name:           req.Name,
		Description:    req.Description,
	}

	id, err := uc.areas.CreateArea(ctx, area)
	if err != nil {
		return nil, err
	}

	area.ID = id
	return area, nil
}
