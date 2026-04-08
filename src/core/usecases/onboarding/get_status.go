package onboarding

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/providers"
)

type GetStatusRequest struct {
	OrgID  uuid.UUID
	UserID int64
}

func (r GetStatusRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.UserID == 0 {
		return fmt.Errorf("%w: user_id is required", providers.ErrValidation)
	}
	return nil
}

type GetStatusResponse struct {
	Completed   bool    `json:"completed"`
	CompletedAt *string `json:"completed_at"`
}

type GetStatus interface {
	Execute(ctx context.Context, req GetStatusRequest) (*GetStatusResponse, error)
}

type getStatusImpl struct {
	users providers.UserProvider
}

func NewGetStatus(users providers.UserProvider) GetStatus {
	return &getStatusImpl{users: users}
}

func (uc *getStatusImpl) Execute(ctx context.Context, req GetStatusRequest) (*GetStatusResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := uc.users.FindByID(ctx, req.OrgID, req.UserID)
	if err != nil {
		return nil, err
	}

	resp := &GetStatusResponse{
		Completed: user.OnboardingCompletedAt != nil,
	}
	if user.OnboardingCompletedAt != nil {
		t := user.OnboardingCompletedAt.Format("2006-01-02T15:04:05Z07:00")
		resp.CompletedAt = &t
	}

	return resp, nil
}
