package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/providers"
)

type RemoveCoordinatorRequest struct {
	OrgID  uuid.UUID
	AreaID int64
	UserID int64
}

func (r RemoveCoordinatorRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.AreaID == 0 {
		return fmt.Errorf("%w: area_id is required", providers.ErrValidation)
	}
	if r.UserID == 0 {
		return fmt.Errorf("%w: user_id is required", providers.ErrValidation)
	}
	return nil
}

type RemoveCoordinator interface {
	Execute(ctx context.Context, req RemoveCoordinatorRequest) error
}

type removeCoordinatorImpl struct {
	areas        providers.AreaProvider
	coordinators providers.AreaCoordinatorProvider
}

func NewRemoveCoordinator(
	areas providers.AreaProvider,
	coordinators providers.AreaCoordinatorProvider,
) RemoveCoordinator {
	return &removeCoordinatorImpl{
		areas:        areas,
		coordinators: coordinators,
	}
}

func (uc *removeCoordinatorImpl) Execute(ctx context.Context, req RemoveCoordinatorRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	// Verify area belongs to the org
	if _, err := uc.areas.GetArea(ctx, req.OrgID, req.AreaID); err != nil {
		return fmt.Errorf("area not found: %w", err)
	}

	return uc.coordinators.Remove(ctx, req.AreaID, req.UserID)
}
