package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%w: course already exists in this organization", providers.ErrConflict)
		}
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

func (r *courseRepo) ListCourses(ctx context.Context, orgID uuid.UUID, p providers.Pagination) ([]entities.Course, bool, error) {
	p = p.Normalize()
	var courses []entities.Course
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Order("name ASC").
		Offset(p.Offset).
		Limit(p.Limit + 1).
		Find(&courses).Error
	if err != nil {
		return nil, false, err
	}
	more := len(courses) > p.Limit
	if more {
		courses = courses[:p.Limit]
	}
	return courses, more, nil
}

// UpdateCourse writes the course's mutable columns scoped to (org, id). Today
// only `name` is updatable: `year` was dropped in migration 000013 and no other
// columns are mutable. Returns ErrNotFound if the row doesn't belong to the
// org so callers can map to a clean 404 without a pre-read.
func (r *courseRepo) UpdateCourse(ctx context.Context, course *entities.Course) error {
	result := r.db.WithContext(ctx).
		Model(&entities.Course{}).
		Where("organization_id = ? AND id = ?", course.OrganizationID, course.ID).
		Updates(map[string]any{
			"name": course.Name,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: course %d", providers.ErrNotFound, course.ID)
	}
	return nil
}

// CountCourseDependencies counts entities that reference this course and would
// block a destructive delete. course_subjects, students and time_slots all
// cascade at the DB level, but we refuse the delete at the API layer instead
// of silently wiping them — callers must clean up first.
func (r *courseRepo) CountCourseDependencies(ctx context.Context, orgID uuid.UUID, id int64) (providers.CourseDependencies, error) {
	var deps providers.CourseDependencies
	// Tenant check first: counts alone can't distinguish "unknown course" from
	// "course with zero deps", and callers rely on a dedicated NotFound for 404.
	var exists int64
	if err := r.db.WithContext(ctx).
		Model(&entities.Course{}).
		Where("organization_id = ? AND id = ?", orgID, id).
		Count(&exists).Error; err != nil {
		return deps, err
	}
	if exists == 0 {
		return deps, fmt.Errorf("%w: course %d", providers.ErrNotFound, id)
	}
	if err := r.db.WithContext(ctx).
		Model(&entities.CourseSubject{}).
		Where("organization_id = ? AND course_id = ?", orgID, id).
		Count(&deps.CourseSubjects).Error; err != nil {
		return deps, err
	}
	if err := r.db.WithContext(ctx).
		Model(&entities.Student{}).
		Where("course_id = ?", id).
		Count(&deps.Students).Error; err != nil {
		return deps, err
	}
	if err := r.db.WithContext(ctx).
		Model(&entities.TimeSlot{}).
		Where("course_id = ?", id).
		Count(&deps.TimeSlots).Error; err != nil {
		return deps, err
	}
	return deps, nil
}

// DeleteCourse removes a course scoped to (org, id). Callers must verify with
// CountCourseDependencies first.
func (r *courseRepo) DeleteCourse(ctx context.Context, orgID uuid.UUID, id int64) error {
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", orgID, id).
		Delete(&entities.Course{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: course %d", providers.ErrNotFound, id)
	}
	return nil
}
