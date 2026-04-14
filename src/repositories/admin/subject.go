package admin

import (
	"context"

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
		Find(&subjects).Error
	return subjects, err
}
