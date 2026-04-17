package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type courseRepo struct {
	db *gorm.DB
}

func NewCourseRepo(db *gorm.DB) providers.CourseProvider {
	return &courseRepo{db: db}
}

func (r *courseRepo) CreateCourse(ctx context.Context, course *entities.Course) (int64, error) {
	if err := r.db.WithContext(ctx).Create(course).Error; err != nil {
		return 0, err
	}
	return course.ID, nil
}

func (r *courseRepo) GetCourse(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Course, error) {
	var course entities.Course
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", orgID, id).
		Preload("Students").
		Preload("CourseSubjects").
		Preload("CourseSubjects.Subject").
		Preload("CourseSubjects.Teacher").
		First(&course).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: course %d", providers.ErrNotFound, id)
		}
		return nil, err
	}
	return &course, nil
}

func (r *courseRepo) ListCourses(ctx context.Context, orgID uuid.UUID) ([]entities.Course, error) {
	var courses []entities.Course
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Order("name ASC").Limit(boundedListCap).
		Find(&courses).Error
	return courses, err
}
