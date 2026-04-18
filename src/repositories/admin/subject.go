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

type subjectRepo struct {
	db *gorm.DB
}

func NewSubjectRepo(db *gorm.DB) providers.SubjectProvider {
	return &subjectRepo{db: db}
}

func (r *subjectRepo) CreateSubject(ctx context.Context, subject *entities.Subject) (int64, error) {
	if err := r.db.WithContext(ctx).Create(subject).Error; err != nil {
		return 0, err
	}
	return subject.ID, nil
}

func (r *subjectRepo) ListSubjectsByArea(ctx context.Context, orgID uuid.UUID, areaID int64) ([]entities.Subject, error) {
	var subjects []entities.Subject
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND area_id = ?", orgID, areaID).
		Order("name ASC").Limit(boundedListCap).
		Find(&subjects).Error
	return subjects, err
}

// ListSubjectsByOrg returns all subjects for an org, optionally filtered by area.
func (r *subjectRepo) ListSubjectsByOrg(ctx context.Context, orgID uuid.UUID, areaID *int64) ([]entities.Subject, error) {
	var subjects []entities.Subject
	q := r.db.WithContext(ctx).Where("organization_id = ?", orgID)
	if areaID != nil {
		q = q.Where("area_id = ?", *areaID)
	}
	err := q.Order("name ASC").Limit(boundedListCap).Find(&subjects).Error
	return subjects, err
}

func (r *subjectRepo) GetSubject(ctx context.Context, orgID uuid.UUID, id int64) (*entities.Subject, error) {
	var subject entities.Subject
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", orgID, id).
		First(&subject).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: subject %d", providers.ErrNotFound, id)
		}
		return nil, err
	}
	return &subject, nil
}

// UpdateSubject writes the subject's mutable columns scoped to (org, id).
// Returns ErrNotFound if the row doesn't belong to the org so callers get a
// clean 404 without a pre-read.
func (r *subjectRepo) UpdateSubject(ctx context.Context, subject *entities.Subject) error {
	result := r.db.WithContext(ctx).
		Model(&entities.Subject{}).
		Where("organization_id = ? AND id = ?", subject.OrganizationID, subject.ID).
		Updates(map[string]any{
			"name":        subject.Name,
			"description": subject.Description,
			"area_id":     subject.AreaID,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: subject %d", providers.ErrNotFound, subject.ID)
	}
	return nil
}

// CountSubjectDependencies counts rows in course_subjects that reference this
// subject. The FK has no ON DELETE action, so without this check DeleteSubject
// would surface a raw FK violation as a 500. Tenant check first so callers
// can map missing rows to 404 without a pre-read.
func (r *subjectRepo) CountSubjectDependencies(ctx context.Context, orgID uuid.UUID, id int64) (providers.SubjectDependencies, error) {
	var deps providers.SubjectDependencies
	var exists int64
	if err := r.db.WithContext(ctx).
		Model(&entities.Subject{}).
		Where("organization_id = ? AND id = ?", orgID, id).
		Count(&exists).Error; err != nil {
		return deps, err
	}
	if exists == 0 {
		return deps, fmt.Errorf("%w: subject %d", providers.ErrNotFound, id)
	}
	if err := r.db.WithContext(ctx).
		Model(&entities.CourseSubject{}).
		Where("organization_id = ? AND subject_id = ?", orgID, id).
		Count(&deps.CourseSubjects).Error; err != nil {
		return deps, err
	}
	return deps, nil
}

// DeleteSubject removes a subject scoped to (org, id). Caller must verify with
// CountSubjectDependencies first.
func (r *subjectRepo) DeleteSubject(ctx context.Context, orgID uuid.UUID, id int64) error {
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", orgID, id).
		Delete(&entities.Subject{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: subject %d", providers.ErrNotFound, id)
	}
	return nil
}
