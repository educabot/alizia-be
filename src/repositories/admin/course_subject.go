package admin

import (
	"context"
	"errors"
	"fmt"

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

func (r *courseSubjectRepo) ListByCourse(ctx context.Context, courseID int64) ([]entities.CourseSubject, error) {
	var results []entities.CourseSubject
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Preload("Subject").
		Preload("Teacher").
		Find(&results).Error
	return results, err
}
