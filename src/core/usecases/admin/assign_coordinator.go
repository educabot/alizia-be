package admin

import (
	"context"

	"github.com/educabot/alizia-be/src/core/providers"
)

type AssignCoordinator interface {
	Execute(ctx context.Context, orgID, areaID, userID int64) error
}

type assignCoordinatorImpl struct {
	areas providers.AreaProvider
}

func NewAssignCoordinator(areas providers.AreaProvider) AssignCoordinator {
	return &assignCoordinatorImpl{areas: areas}
}

func (uc *assignCoordinatorImpl) Execute(ctx context.Context, orgID, areaID, userID int64) error {
	area, err := uc.areas.GetArea(ctx, orgID, areaID)
	if err != nil {
		return err
	}

	area.CoordinatorID = &userID
	_, err = uc.areas.CreateArea(ctx, area) // TODO: update method
	return err
}
