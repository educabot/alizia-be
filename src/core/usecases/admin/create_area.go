package admin

import (
	"context"
	"fmt"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type CreateAreaRequest struct {
	OrgID int64
	Name  string
}

func (r CreateAreaRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("%w: name is required", providers.ErrValidation)
	}
	return nil
}

type CreateArea interface {
	Execute(ctx context.Context, req CreateAreaRequest) (int64, error)
}

type createAreaImpl struct {
	areas providers.AreaProvider
}

func NewCreateArea(areas providers.AreaProvider) CreateArea {
	return &createAreaImpl{areas: areas}
}

func (uc *createAreaImpl) Execute(ctx context.Context, req CreateAreaRequest) (int64, error) {
	if err := req.Validate(); err != nil {
		return 0, err
	}

	area := &entities.Area{
		OrganizationID: req.OrgID,
		Name:           req.Name,
	}

	return uc.areas.CreateArea(ctx, area)
}
