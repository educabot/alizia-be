package admin

import (
	"context"
	"encoding/json"
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

func (r *organizationRepo) UpdateConfig(ctx context.Context, id uuid.UUID, configPatch map[string]any) (*entities.Organization, error) {
	patchJSON, err := json.Marshal(configPatch)
	if err != nil {
		return nil, err
	}

	// PostgreSQL || operator: shallow-merges the patch into existing config.
	// Existing keys not in the patch are preserved; keys in the patch overwrite.
	result := r.db.WithContext(ctx).
		Model(&entities.Organization{}).
		Where("id = ?", id).
		Update("config", gorm.Expr("config || ?::jsonb", string(patchJSON)))
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, providers.ErrNotFound
	}

	return r.FindByID(ctx, id)
}
