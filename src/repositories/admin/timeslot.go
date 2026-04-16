package admin

import (
	"context"

	"gorm.io/gorm"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/repositories/admin/queries"
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
		Preload("Subjects.CourseSubject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, course_id, subject_id, teacher_id, school_year")
		}).
		Preload("Subjects.CourseSubject.Subject", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, area_id, name")
		}).
		Preload("Subjects.CourseSubject.Teacher", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, organization_id, email, first_name, last_name, avatar_url")
		}).
		Order("day_of_week, start_time").
		Limit(100).
		Find(&slots).Error
	return slots, err
}

// GetSharedClassNumbers projects the weekly schedule of a course_subject onto a
// total class count and returns the 1-based class numbers that are shared
// (i.e. happen in a time_slot containing more than one course_subject).
func (r *timeSlotRepo) GetSharedClassNumbers(ctx context.Context, courseSubjectID int64, totalClasses int) ([]int, error) {
	type weeklySlot struct {
		WeeklyPosition int  `gorm:"column:weekly_position"`
		IsShared       bool `gorm:"column:is_shared"`
	}
	var slots []weeklySlot
	err := r.db.WithContext(ctx).Raw(queries.SharedClassNumbers, courseSubjectID).Scan(&slots).Error
	if err != nil {
		return nil, err
	}
	if len(slots) == 0 {
		return []int{}, nil
	}

	classesPerWeek := len(slots)
	shared := []int{}
	for classNum := 1; classNum <= totalClasses; classNum++ {
		weekPos := (classNum - 1) % classesPerWeek
		if slots[weekPos].IsShared {
			shared = append(shared, classNum)
		}
	}
	return shared, nil
}
