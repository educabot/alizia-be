package admin

import (
	"context"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type ListAreas interface {
	Execute(ctx context.Context, orgID int64) ([]entities.Area, error)
}

type listAreasImpl struct {
	areas providers.AreaProvider
}

func NewListAreas(areas providers.AreaProvider) ListAreas {
	return &listAreasImpl{areas: areas}
}

func (uc *listAreasImpl) Execute(ctx context.Context, orgID int64) ([]entities.Area, error) {
	return uc.areas.ListAreas(ctx, orgID)
}
