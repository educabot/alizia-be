package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type AssignCoordinatorRequest struct {
	OrgID  uuid.UUID
	AreaID int64
	UserID int64
}

func (r AssignCoordinatorRequest) Validate() error {
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

type AssignCoordinator interface {
	Execute(ctx context.Context, req AssignCoordinatorRequest) (*entities.AreaCoordinator, error)
}

type assignCoordinatorImpl struct {
	areas        providers.AreaProvider
	users        providers.UserProvider
	coordinators providers.AreaCoordinatorProvider
}

func NewAssignCoordinator(
	areas providers.AreaProvider,
	users providers.UserProvider,
	coordinators providers.AreaCoordinatorProvider,
) AssignCoordinator {
	return &assignCoordinatorImpl{
		areas:        areas,
		users:        users,
		coordinators: coordinators,
	}
}

func (uc *assignCoordinatorImpl) Execute(ctx context.Context, req AssignCoordinatorRequest) (*entities.AreaCoordinator, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Verify area belongs to the org
	if _, err := uc.areas.GetArea(ctx, req.OrgID, req.AreaID); err != nil {
		return nil, fmt.Errorf("area not found: %w", err)
	}

	// Verify user belongs to the org and has coordinator role
	user, err := uc.users.FindByID(ctx, req.OrgID, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if !user.HasRole(entities.RoleCoordinator) && !user.HasRole(entities.RoleAdmin) {
		return nil, fmt.Errorf("%w: user must have coordinator or admin role", providers.ErrValidation)
	}

	ac, err := uc.coordinators.Assign(ctx, req.AreaID, req.UserID)
	if err != nil {
		return nil, err
	}

	return ac, nil
}
