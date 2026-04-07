package admin_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/repositories/admin"
)

func TestOrganizationRepo_FindByID(t *testing.T) {
	db := testDB(t)
	repo := admin.NewOrganizationRepo(db)
	ctx := context.Background()
	slug := uniqueSlug(t)
	defer cleanupTestData(t, db, slug)

	org := entities.Organization{Name: "Test Org", Slug: slug}
	require.NoError(t, db.Create(&org).Error)

	found, err := repo.FindByID(ctx, org.ID)
	require.NoError(t, err)
	assert.Equal(t, org.ID, found.ID)
	assert.Equal(t, "Test Org", found.Name)
	assert.Equal(t, slug, found.Slug)
}

func TestOrganizationRepo_FindByID_NotFound(t *testing.T) {
	db := testDB(t)
	repo := admin.NewOrganizationRepo(db)
	ctx := context.Background()

	_, err := repo.FindByID(ctx, uuid.New())
	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestOrganizationRepo_FindBySlug(t *testing.T) {
	db := testDB(t)
	repo := admin.NewOrganizationRepo(db)
	ctx := context.Background()
	slug := uniqueSlug(t)
	defer cleanupTestData(t, db, slug)

	org := entities.Organization{Name: "Slug Org", Slug: slug}
	require.NoError(t, db.Create(&org).Error)

	found, err := repo.FindBySlug(ctx, slug)
	require.NoError(t, err)
	assert.Equal(t, org.ID, found.ID)
	assert.Equal(t, slug, found.Slug)
}

func TestOrganizationRepo_FindBySlug_NotFound(t *testing.T) {
	db := testDB(t)
	repo := admin.NewOrganizationRepo(db)
	ctx := context.Background()

	_, err := repo.FindBySlug(ctx, "nonexistent-slug")
	assert.ErrorIs(t, err, providers.ErrNotFound)
}
