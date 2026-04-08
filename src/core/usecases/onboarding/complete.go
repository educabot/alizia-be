package onboarding

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/providers"
)

type CompleteRequest struct {
	OrgID  uuid.UUID
	UserID int64
}

func (r CompleteRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.UserID == 0 {
		return fmt.Errorf("%w: user_id is required", providers.ErrValidation)
	}
	return nil
}

type Complete interface {
	Execute(ctx context.Context, req CompleteRequest) error
}

type completeImpl struct {
	users providers.UserProvider
}

func NewComplete(users providers.UserProvider) Complete {
	return &completeImpl{users: users}
}

func (uc *completeImpl) Execute(ctx context.Context, req CompleteRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	// Verify user exists in the org
	if _, err := uc.users.FindByID(ctx, req.OrgID, req.UserID); err != nil {
		return err
	}

	return uc.users.CompleteOnboarding(ctx, req.UserID)
}
