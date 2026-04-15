package admin

import (
	"context"

	"gorm.io/gorm"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type timeSlotRepo struct {
	db *gorm.DB
}

func NewTimeSlotRepo(db *gorm.DB) providers.TimeSlotProvider {
	return &timeSlotRepo{db: db}
}

func (r *timeSlotRepo) CreateTimeSlot(ctx context.Context, slot *entities.TimeSlot) (int64, error) {
	if err := r.db.WithContext(ctx).Create(slot).Error; err != nil {
		return 0, err
	}
	return slot.ID, nil
}

func (r *timeSlotRepo) ListByCourse(ctx context.Context, courseID int64) ([]entities.TimeSlot, error) {
	var slots []entities.TimeSlot
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Preload("Subjects").
		Preload("Subjects.CourseSubject").
		Preload("Subjects.CourseSubject.Subject").
		Preload("Subjects.CourseSubject.Teacher").
		Order("day_of_week, start_time").
		Limit(100).
		Find(&slots).Error
	return slots, err
}
