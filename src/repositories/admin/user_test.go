package admin_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/repositories/admin"
)

func setupUserTest(t *testing.T) (context.Context, providers.UserProvider, entities.Organization, func()) {
	t.Helper()
	db := testDB(t)
	repo := admin.NewUserRepo(db)
	ctx := context.Background()
	slug := uniqueSlug(t)

	org := entities.Organization{Name: "User Test Org", Slug: slug}
	require.NoError(t, db.Create(&org).Error)

	cleanup := func() { cleanupTestData(t, db, slug) }
	return ctx, repo, org, cleanup
}

func TestUserRepo_Create_And_FindByID(t *testing.T) {
	ctx, repo, org, cleanup := setupUserTest(t)
	defer cleanup()

	user := &entities.User{
		OrganizationID: org.ID,
		Email:          "test@example.com",
		Name:           "Test User",
		Roles:          []entities.UserRole{{Role: entities.RoleTeacher}},
	}

	id, err := repo.Create(ctx, user)
	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	found, err := repo.FindByID(ctx, org.ID, id)
	require.NoError(t, err)
	assert.Equal(t, "Test User", found.Name)
	assert.Equal(t, "test@example.com", found.Email)
	assert.Len(t, found.Roles, 1)
	assert.Equal(t, entities.RoleTeacher, found.Roles[0].Role)
}

func TestUserRepo_FindByID_NotFound(t *testing.T) {
	ctx, repo, org, cleanup := setupUserTest(t)
	defer cleanup()

	_, err := repo.FindByID(ctx, org.ID, 99999)
	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestUserRepo_FindByID_WrongOrg(t *testing.T) {
	ctx, repo, org, cleanup := setupUserTest(t)
	defer cleanup()

	user := &entities.User{
		OrganizationID: org.ID,
		Email:          "isolated@example.com",
		Name:           "Isolated User",
	}
	id, err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Search with different org → not found
	_, err = repo.FindByID(ctx, uuid.New(), id)
	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestUserRepo_FindByEmail(t *testing.T) {
	ctx, repo, org, cleanup := setupUserTest(t)
	defer cleanup()

	user := &entities.User{
		OrganizationID: org.ID,
		Email:          "find-me@example.com",
		Name:           "Findable User",
	}
	_, err := repo.Create(ctx, user)
	require.NoError(t, err)

	found, err := repo.FindByEmail(ctx, org.ID, "find-me@example.com")
	require.NoError(t, err)
	assert.Equal(t, "Findable User", found.Name)
}

func TestUserRepo_FindByEmail_NotFound(t *testing.T) {
	ctx, repo, org, cleanup := setupUserTest(t)
	defer cleanup()

	_, err := repo.FindByEmail(ctx, org.ID, "nope@example.com")
	assert.ErrorIs(t, err, providers.ErrNotFound)
}

func TestUserRepo_FindByOrgID(t *testing.T) {
	ctx, repo, org, cleanup := setupUserTest(t)
	defer cleanup()

	for i, name := range []string{"User A", "User B", "User C"} {
		u := &entities.User{
			OrganizationID: org.ID,
			Email:          fmt.Sprintf("user%d@example.com", i),
			Name:           name,
		}
		_, err := repo.Create(ctx, u)
		require.NoError(t, err)
	}

	users, err := repo.FindByOrgID(ctx, org.ID)
	require.NoError(t, err)
	assert.Len(t, users, 3)
}

func TestUserRepo_AssignRole(t *testing.T) {
	ctx, repo, org, cleanup := setupUserTest(t)
	defer cleanup()

	user := &entities.User{
		OrganizationID: org.ID,
		Email:          "role-test@example.com",
		Name:           "Role User",
	}
	id, err := repo.Create(ctx, user)
	require.NoError(t, err)

	err = repo.AssignRole(ctx, id, entities.RoleCoordinator)
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, org.ID, id)
	require.NoError(t, err)
	assert.True(t, found.HasRole(entities.RoleCoordinator))
}

func TestUserRepo_AssignRole_Duplicate(t *testing.T) {
	ctx, repo, org, cleanup := setupUserTest(t)
	defer cleanup()

	user := &entities.User{
		OrganizationID: org.ID,
		Email:          "dup-role@example.com",
		Name:           "Dup Role User",
		Roles:          []entities.UserRole{{Role: entities.RoleTeacher}},
	}
	id, err := repo.Create(ctx, user)
	require.NoError(t, err)

	err = repo.AssignRole(ctx, id, entities.RoleTeacher)
	assert.ErrorIs(t, err, providers.ErrDuplicate)
}

func TestUserRepo_RemoveRole(t *testing.T) {
	ctx, repo, org, cleanup := setupUserTest(t)
	defer cleanup()

	user := &entities.User{
		OrganizationID: org.ID,
		Email:          "remove-role@example.com",
		Name:           "Remove Role User",
		Roles: []entities.UserRole{
			{Role: entities.RoleTeacher},
			{Role: entities.RoleCoordinator},
		},
	}
	id, err := repo.Create(ctx, user)
	require.NoError(t, err)

	err = repo.RemoveRole(ctx, id, entities.RoleTeacher)
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, org.ID, id)
	require.NoError(t, err)
	assert.Len(t, found.Roles, 1)
	assert.Equal(t, entities.RoleCoordinator, found.Roles[0].Role)
}

func TestUserRepo_RemoveRole_NotFound(t *testing.T) {
	ctx, repo, org, cleanup := setupUserTest(t)
	defer cleanup()

	user := &entities.User{
		OrganizationID: org.ID,
		Email:          "no-role@example.com",
		Name:           "No Role User",
	}
	id, err := repo.Create(ctx, user)
	require.NoError(t, err)

	err = repo.RemoveRole(ctx, id, entities.RoleAdmin)
	assert.ErrorIs(t, err, providers.ErrNotFound)
}
