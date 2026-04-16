package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type UpdateAreaRequest struct {
	OrgID  uuid.UUID
	AreaID int64
	// Name is updated when non-nil. An empty string is rejected.
	Name *string
	// Description is updated when SetDescription is true. A nil Description with
	// SetDescription=true clears the field; a non-nil value sets it.
	Description    *string
	SetDescription bool
}

func (r UpdateAreaRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.AreaID <= 0 {
		return fmt.Errorf("%w: area_id is required", providers.ErrValidation)
	}
	if r.Name != nil && *r.Name == "" {
		return fmt.Errorf("%w: name cannot be empty", providers.ErrValidation)
	}
	return nil
}

type UpdateArea interface {
	Execute(ctx context.Context, req UpdateAreaRequest) (*entities.Area, error)
}

type updateAreaImpl struct {
	areas providers.AreaProvider
}

func NewUpdateArea(areas providers.AreaProvider) UpdateArea {
	return &updateAreaImpl{areas: areas}
}

func (uc *updateAreaImpl) Execute(ctx context.Context, req UpdateAreaRequest) (*entities.Area, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	area, err := uc.areas.GetArea(ctx, req.OrgID, req.AreaID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		area.Name = *req.Name
	}
	if req.SetDescription {
		area.Description = req.Description
	}

	if err := uc.areas.UpdateArea(ctx, area); err != nil {
		return nil, err
	}
	return area, nil
}
