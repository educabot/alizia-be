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

type courseSubjectRepo struct {
	db *gorm.DB
}

func NewCourseSubjectRepo(db *gorm.DB) providers.CourseSubjectProvider {
	return &courseSubjectRepo{db: db}
}

func (r *courseSubjectRepo) CreateCourseSubject(ctx context.Context, cs *entities.CourseSubject) (int64, error) {
	if err := r.db.WithContext(ctx).Create(cs).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%w: subject already assigned to course for this year", providers.ErrConflict)
		}
		return 0, err
	}
	return cs.ID, nil
}

func (r *courseSubjectRepo) GetCourseSubject(ctx context.Context, orgID uuid.UUID, id int64) (*entities.CourseSubject, error) {
	var cs entities.CourseSubject
	err := r.db.WithContext(ctx).
		Preload("Subject").
		Preload("Teacher").
		Where("organization_id = ? AND id = ?", orgID, id).
		First(&cs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: course-subject %d", providers.ErrNotFound, id)
		}
		return nil, err
	}
	return &cs, nil
}

func (r *courseSubjectRepo) ListByCourse(ctx context.Context, courseID int64) ([]entities.CourseSubject, error) {
	var results []entities.CourseSubject
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Preload("Subject").
		Preload("Teacher").
		Limit(100).
		Find(&results).Error
	return results, err
}

// ListCourseSubjects returns course-subjects for an org, applying optional
// filters. We filter directly on course_subjects.organization_id (denormalized
// in the schema, indexed via idx_course_subjects_org_year) instead of joining
// with courses — same tenant guarantee, fewer rows scanned.
func (r *courseSubjectRepo) ListCourseSubjects(ctx context.Context, orgID uuid.UUID, filter providers.CourseSubjectFilter) ([]entities.CourseSubject, error) {
	var results []entities.CourseSubject
	q := r.db.WithContext(ctx).Where("organization_id = ?", orgID)

	if filter.CourseID != nil {
		q = q.Where("course_id = ?", *filter.CourseID)
	}
	if filter.SubjectID != nil {
		q = q.Where("subject_id = ?", *filter.SubjectID)
	}
	if filter.TeacherID != nil {
		q = q.Where("teacher_id = ?", *filter.TeacherID)
	}

	err := q.Preload("Subject").
		Preload("Teacher").
		Order("course_id ASC, subject_id ASC").
		Limit(200).
		Find(&results).Error
	return results, err
}
