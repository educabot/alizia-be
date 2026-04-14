package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type ListAreasRequest struct {
	OrgID uuid.UUID
}

func (r ListAreasRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	return nil
}

type ListAreas interface {
	Execute(ctx context.Context, req ListAreasRequest) ([]entities.Area, error)
}

type listAreasImpl struct {
	areas providers.AreaProvider
}

func NewListAreas(areas providers.AreaProvider) ListAreas {
	return &listAreasImpl{areas: areas}
}

func (uc *listAreasImpl) Execute(ctx context.Context, req ListAreasRequest) ([]entities.Area, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return uc.areas.ListAreas(ctx, req.OrgID)
}
