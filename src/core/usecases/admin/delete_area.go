package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/providers"
)

type DeleteAreaRequest struct {
	OrgID  uuid.UUID
	AreaID int64
}

func (r DeleteAreaRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.AreaID <= 0 {
		return fmt.Errorf("%w: area_id is required", providers.ErrValidation)
	}
	return nil
}

type DeleteArea interface {
	Execute(ctx context.Context, req DeleteAreaRequest) error
}

type deleteAreaImpl struct {
	areas providers.AreaProvider
}

func NewDeleteArea(areas providers.AreaProvider) DeleteArea {
	return &deleteAreaImpl{areas: areas}
}

// Execute deletes an area after verifying it has no blocking dependencies.
// Returns ErrNotFound if the area doesn't belong to the caller's org, or
// ErrConflict with a message listing the blocking entity counts if subjects
// or coordination documents still reference the area. We refuse to cascade
// because those deletions are destructive and not recoverable via the API.
func (uc *deleteAreaImpl) Execute(ctx context.Context, req DeleteAreaRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	// Verify the area exists in this org before counting dependencies —
	// keeps the error surface simple (404 > 409 when both apply).
	if _, err := uc.areas.GetArea(ctx, req.OrgID, req.AreaID); err != nil {
		return err
	}

	deps, err := uc.areas.CountDependencies(ctx, req.OrgID, req.AreaID)
	if err != nil {
		return err
	}
	if !deps.IsEmpty() {
		return fmt.Errorf("%w: area has dependencies (%s); remove them before deleting",
			providers.ErrConflict, formatDependencies(deps))
	}

	return uc.areas.DeleteArea(ctx, req.OrgID, req.AreaID)
}

func formatDependencies(d providers.AreaDependencies) string {
	parts := make([]string, 0, 2)
	if d.Subjects > 0 {
		parts = append(parts, fmt.Sprintf("%d subjects", d.Subjects))
	}
	if d.CoordinationDocuments > 0 {
		parts = append(parts, fmt.Sprintf("%d coordination documents", d.CoordinationDocuments))
	}
	return strings.Join(parts, ", ")
}
