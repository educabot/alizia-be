package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type ListUsersRequest struct {
	OrgID      uuid.UUID
	Role       *string
	AreaID     *int64
	Search     *string
	Pagination providers.Pagination
}

func (r ListUsersRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.Role != nil && !isValidRole(*r.Role) {
		return fmt.Errorf("%w: role must be admin, coordinator or teacher", providers.ErrValidation)
	}
	if r.AreaID != nil && *r.AreaID <= 0 {
		return fmt.Errorf("%w: area_id must be positive", providers.ErrValidation)
	}
	return nil
}

func isValidRole(role string) bool {
	switch entities.Role(role) {
	case entities.RoleAdmin, entities.RoleCoordinator, entities.RoleTeacher:
		return true
	}
	return false
}

type ListUsersResponse struct {
	Items []entities.User
	More  bool
}

type ListUsers interface {
	Execute(ctx context.Context, req ListUsersRequest) (ListUsersResponse, error)
}

type listUsersImpl struct {
	users providers.UserProvider
}

func NewListUsers(users providers.UserProvider) ListUsers {
	return &listUsersImpl{users: users}
}

func (uc *listUsersImpl) Execute(ctx context.Context, req ListUsersRequest) (ListUsersResponse, error) {
	if err := req.Validate(); err != nil {
		return ListUsersResponse{}, err
	}

	filter := providers.UserFilter{
		AreaID: req.AreaID,
		Search: req.Search,
	}
	if req.Role != nil {
		r := entities.Role(*req.Role)
		filter.Role = &r
	}

	items, more, err := uc.users.ListUsers(ctx, req.OrgID, filter, req.Pagination)
	if err != nil {
		return ListUsersResponse{}, err
	}
	return ListUsersResponse{Items: items, More: more}, nil
}
