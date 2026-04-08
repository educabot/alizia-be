package admin

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
)

type organizationRepo struct {
	db *gorm.DB
}

func NewOrganizationRepo(db *gorm.DB) providers.OrganizationProvider {
	return &organizationRepo{db: db}
}

func (r *organizationRepo) FindByID(ctx context.Context, id uuid.UUID) (*entities.Organization, error) {
	var org entities.Organization
	err := r.db.WithContext(ctx).First(&org, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, providers.ErrNotFound
	}
	return &org, err
}

func (r *organizationRepo) FindBySlug(ctx context.Context, slug string) (*entities.Organization, error) {
	var org entities.Organization
	err := r.db.WithContext(ctx).First(&org, "slug = ?", slug).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, providers.ErrNotFound
	}
	return &org, err
}
