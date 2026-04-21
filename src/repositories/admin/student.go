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

type studentRepo struct {
	db *gorm.DB
}

func NewStudentRepo(db *gorm.DB) providers.StudentProvider {
	return &studentRepo{db: db}
}

func (r *studentRepo) CreateStudent(ctx context.Context, student *entities.Student) (int64, error) {
	if err := r.db.WithContext(ctx).Create(student).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%w: student already exists in this course", providers.ErrConflict)
		}
		return 0, err
	}
	return student.ID, nil
}

func (r *studentRepo) ListByCourse(ctx context.Context, courseID int64) ([]entities.Student, error) {
	var students []entities.Student
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Limit(boundedListCap).
		Find(&students).Error
	return students, err
}
