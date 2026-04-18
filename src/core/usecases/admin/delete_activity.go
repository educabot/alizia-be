package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/providers"
)

type DeleteActivityRequest struct {
	OrgID      uuid.UUID
	ActivityID int64
}

func (r DeleteActivityRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.ActivityID <= 0 {
		return fmt.Errorf("%w: activity_id is required", providers.ErrValidation)
	}
	return nil
}

type DeleteActivity interface {
	Execute(ctx context.Context, req DeleteActivityRequest) error
}

type deleteActivityImpl struct {
	activities providers.ActivityTemplateProvider
}

func NewDeleteActivity(activities providers.ActivityTemplateProvider) DeleteActivity {
	return &deleteActivityImpl{activities: activities}
}

// Execute removes an activity template without a dependency counter: no
// referencing table exists yet (lesson_plans and coordination_documents ship
// in Épica 4/5). When those land, add CountActivityDependencies and gate the
// delete with ErrConflict the same way DeleteArea / DeleteCourse do.
func (uc *deleteActivityImpl) Execute(ctx context.Context, req DeleteActivityRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	return uc.activities.DeleteActivity(ctx, req.OrgID, req.ActivityID)
}
