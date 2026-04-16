package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/educabot/alizia-be/src/core/providers"
)

type GetSharedClassNumbersRequest struct {
	OrgID           uuid.UUID
	CourseSubjectID int64
	TotalClasses    int
}

func (r GetSharedClassNumbersRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return fmt.Errorf("%w: org_id is required", providers.ErrValidation)
	}
	if r.CourseSubjectID == 0 {
		return fmt.Errorf("%w: course_subject_id is required", providers.ErrValidation)
	}
	if r.TotalClasses <= 0 {
		return fmt.Errorf("%w: total_classes must be greater than 0", providers.ErrValidation)
	}
	return nil
}

type GetSharedClassNumbersResponse struct {
	CourseSubjectID    int64 `json:"course_subject_id"`
	TotalClasses       int   `json:"total_classes"`
	SharedClassNumbers []int `json:"shared_class_numbers"`
}

type GetSharedClassNumbers interface {
	Execute(ctx context.Context, req GetSharedClassNumbersRequest) (*GetSharedClassNumbersResponse, error)
}

type getSharedClassNumbersImpl struct {
	timeSlots providers.TimeSlotProvider
}

func NewGetSharedClassNumbers(timeSlots providers.TimeSlotProvider) GetSharedClassNumbers {
	return &getSharedClassNumbersImpl{timeSlots: timeSlots}
}

func (uc *getSharedClassNumbersImpl) Execute(ctx context.Context, req GetSharedClassNumbersRequest) (*GetSharedClassNumbersResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	numbers, err := uc.timeSlots.GetSharedClassNumbers(ctx, req.CourseSubjectID, req.TotalClasses)
	if err != nil {
		return nil, err
	}

	return &GetSharedClassNumbersResponse{
		CourseSubjectID:    req.CourseSubjectID,
		TotalClasses:       req.TotalClasses,
		SharedClassNumbers: numbers,
	}, nil
}
