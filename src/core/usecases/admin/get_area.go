package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type GetAreaRequest struct {
	OrgID  uuid.UUID
	AreaID int64
}

func (r GetAreaRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.AreaID <= 0 {
		return fmt.Errorf("%w: area_id is required", providers.ErrValidation)
	}
	return nil
}

type GetArea interface {
	Execute(ctx context.Context, req GetAreaRequest) (*entities.Area, error)
}

type getAreaImpl struct {
	areas providers.AreaProvider
}

func NewGetArea(areas providers.AreaProvider) GetArea {
	return &getAreaImpl{areas: areas}
}

func (uc *getAreaImpl) Execute(ctx context.Context, req GetAreaRequest) (*entities.Area, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return uc.areas.GetArea(ctx, req.OrgID, req.AreaID)
}
