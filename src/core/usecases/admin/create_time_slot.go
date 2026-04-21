package admin

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type CreateTimeSlotRequest struct {
	OrgID            uuid.UUID
	CourseID         int64
	DayOfWeek        int
	StartTime        string
	EndTime          string
	CourseSubjectIDs []int64
}

var timeFormatRe = regexp.MustCompile(`^\d{2}:\d{2}$`)

func (r CreateTimeSlotRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.CourseID == 0 {
		return fmt.Errorf("%w: course_id is required", providers.ErrValidation)
	}
	if r.DayOfWeek < 0 || r.DayOfWeek > 6 {
		return fmt.Errorf("%w: day_of_week must be between 0 and 6", providers.ErrValidation)
	}
	if r.StartTime == "" {
		return fmt.Errorf("%w: start_time is required", providers.ErrValidation)
	}
	if r.EndTime == "" {
		return fmt.Errorf("%w: end_time is required", providers.ErrValidation)
	}
	if !timeFormatRe.MatchString(r.StartTime) {
		return fmt.Errorf("%w: start_time must be in HH:MM format", providers.ErrValidation)
	}
	if !timeFormatRe.MatchString(r.EndTime) {
		return fmt.Errorf("%w: end_time must be in HH:MM format", providers.ErrValidation)
	}
	if r.StartTime >= r.EndTime {
		return fmt.Errorf("%w: start_time must be before end_time", providers.ErrValidation)
	}
	if len(r.CourseSubjectIDs) == 0 {
		return fmt.Errorf("%w: at least one course_subject_id is required", providers.ErrValidation)
	}
	if len(r.CourseSubjectIDs) > 2 {
		return fmt.Errorf("%w: maximum 2 course_subject_ids per slot", providers.ErrValidation)
	}
	return nil
}

type CreateTimeSlot interface {
	Execute(ctx context.Context, req CreateTimeSlotRequest) (*entities.TimeSlot, error)
}

type createTimeSlotImpl struct {
	orgs      providers.OrganizationProvider
	courses   providers.CourseProvider
	timeSlots providers.TimeSlotProvider
}

func NewCreateTimeSlot(
	orgs providers.OrganizationProvider,
	courses providers.CourseProvider,
	timeSlots providers.TimeSlotProvider,
) CreateTimeSlot {
	return &createTimeSlotImpl{orgs: orgs, courses: courses, timeSlots: timeSlots}
}

func (uc *createTimeSlotImpl) Execute(ctx context.Context, req CreateTimeSlotRequest) (*entities.TimeSlot, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if _, err := uc.courses.GetCourse(ctx, req.OrgID, req.CourseID); err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	if len(req.CourseSubjectIDs) > 1 {
		org, err := uc.orgs.FindByID(ctx, req.OrgID)
		if err != nil {
			return nil, err
		}
		if !entities.ParseOrgConfig(org.Config).SharedClassesEnabled {
			return nil, fmt.Errorf("%w: shared classes are not enabled for this organization", providers.ErrValidation)
		}
	}

	slot := &entities.TimeSlot{
		CourseID:  req.CourseID,
		DayOfWeek: req.DayOfWeek,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	for _, csID := range req.CourseSubjectIDs {
		slot.Subjects = append(slot.Subjects, entities.TimeSlotSubject{
			CourseSubjectID: csID,
		})
	}

	id, err := uc.timeSlots.CreateTimeSlot(ctx, slot)
	if err != nil {
		return nil, err
	}

	slot.ID = id
	return slot, nil
}
