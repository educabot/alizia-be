package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type GetOrganizationRequest struct {
	OrgID uuid.UUID
}

func (r GetOrganizationRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	return nil
}

type GetOrganization interface {
	Execute(ctx context.Context, req GetOrganizationRequest) (*entities.Organization, error)
}

type getOrganizationImpl struct {
	orgs providers.OrganizationProvider
}

func NewGetOrganization(orgs providers.OrganizationProvider) GetOrganization {
	return &getOrganizationImpl{orgs: orgs}
}

func (uc *getOrganizationImpl) Execute(ctx context.Context, req GetOrganizationRequest) (*entities.Organization, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return uc.orgs.FindByID(ctx, req.OrgID)
}
