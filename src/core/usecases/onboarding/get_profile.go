package onboarding

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/providers"
)

type GetProfileRequest struct {
	OrgID  uuid.UUID
	UserID int64
}

func (r GetProfileRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.UserID == 0 {
		return fmt.Errorf("%w: user_id is required", providers.ErrValidation)
	}
	return nil
}

type GetProfile interface {
	Execute(ctx context.Context, req GetProfileRequest) (map[string]any, error)
}

type getProfileImpl struct {
	users providers.UserProvider
}

func NewGetProfile(users providers.UserProvider) GetProfile {
	return &getProfileImpl{users: users}
}

func (uc *getProfileImpl) Execute(ctx context.Context, req GetProfileRequest) (map[string]any, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := uc.users.FindByID(ctx, req.OrgID, req.UserID)
	if err != nil {
		return nil, err
	}

	var data map[string]any
	if len(user.ProfileData) == 0 {
		return map[string]any{}, nil
	}
	if err := json.Unmarshal(user.ProfileData, &data); err != nil {
		return nil, err
	}

	return data, nil
}
