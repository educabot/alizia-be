package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type GetScheduleRequest struct {
	OrgID    uuid.UUID
	CourseID int64
}

func (r GetScheduleRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.CourseID == 0 {
		return fmt.Errorf("%w: course_id is required", providers.ErrValidation)
	}
	return nil
}

type GetSchedule interface {
	Execute(ctx context.Context, req GetScheduleRequest) ([]entities.TimeSlot, error)
}

type getScheduleImpl struct {
	courses   providers.CourseProvider
	timeSlots providers.TimeSlotProvider
}

func NewGetSchedule(courses providers.CourseProvider, timeSlots providers.TimeSlotProvider) GetSchedule {
	return &getScheduleImpl{courses: courses, timeSlots: timeSlots}
}

func (uc *getScheduleImpl) Execute(ctx context.Context, req GetScheduleRequest) ([]entities.TimeSlot, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Verify course belongs to the org
	if _, err := uc.courses.GetCourse(ctx, req.OrgID, req.CourseID); err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	return uc.timeSlots.ListByCourse(ctx, req.CourseID)
}
