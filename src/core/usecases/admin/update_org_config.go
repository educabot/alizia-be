package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type UpdateOrgConfigRequest struct {
	OrgID       uuid.UUID
	ConfigPatch map[string]any
}

func (r UpdateOrgConfigRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if len(r.ConfigPatch) == 0 {
		return fmt.Errorf("%w: config patch cannot be empty", providers.ErrValidation)
	}
	return nil
}

type UpdateOrgConfig interface {
	Execute(ctx context.Context, req UpdateOrgConfigRequest) (*entities.Organization, error)
}

type updateOrgConfigImpl struct {
	orgs providers.OrganizationProvider
}

func NewUpdateOrgConfig(orgs providers.OrganizationProvider) UpdateOrgConfig {
	return &updateOrgConfigImpl{orgs: orgs}
}

func (uc *updateOrgConfigImpl) Execute(ctx context.Context, req UpdateOrgConfigRequest) (*entities.Organization, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Verify org exists before attempting update
	if _, err := uc.orgs.FindByID(ctx, req.OrgID); err != nil {
		return nil, err
	}

	return uc.orgs.UpdateConfig(ctx, req.OrgID, req.ConfigPatch)
}
