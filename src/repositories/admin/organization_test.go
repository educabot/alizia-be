package admin_test

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/repositories/admin"
	"github.com/educabot/alizia-be/src/utils/dbtest"
)

var orgColumns = []string{"id", "name", "slug", "config", "created_at", "updated_at"}

func orgRow(o entities.Organization) []driver.Value {
	return []driver.Value{o.ID, o.Name, o.Slug, []byte(o.Config), o.CreatedAt, o.UpdatedAt}
}

func TestOrganizationRepo_FindByID_Found(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewOrganizationRepo(gdb)

	id := uuid.New()
	sql := `SELECT * FROM "organizations" WHERE id = $1 ORDER BY "organizations"."id" LIMIT $2`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows(orgColumns).
			AddRow(orgRow(entities.Organization{ID: id, Name: "Acme", Slug: "acme", Config: []byte(`{}`)})...))

	got, err := repo.FindByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, "Acme", got.Name)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrganizationRepo_FindByID_NotFound(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewOrganizationRepo(gdb)

	id := uuid.New()
	sql := `SELECT * FROM "organizations" WHERE id = $1 ORDER BY "organizations"."id" LIMIT $2`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows(orgColumns))

	_, err := repo.FindByID(context.Background(), id)
	assert.ErrorIs(t, err, providers.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrganizationRepo_FindBySlug(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewOrganizationRepo(gdb)

	id := uuid.New()
	sql := `SELECT * FROM "organizations" WHERE slug = $1 ORDER BY "organizations"."id" LIMIT $2`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs("acme", 1).
		WillReturnRows(sqlmock.NewRows(orgColumns).
			AddRow(orgRow(entities.Organization{ID: id, Name: "Acme", Slug: "acme", Config: []byte(`{}`)})...))

	got, err := repo.FindBySlug(context.Background(), "acme")
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

// TestOrganizationRepo_UpdateConfig verifies the JSONB merge path: raw SQL uses
// the Postgres `||` operator with a jsonb cast. The patch is marshalled to JSON
// and passed as a positional arg — we assert that exact shape is sent to the DB.
func TestOrganizationRepo_UpdateConfig_MergesAndReturnsReloaded(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewOrganizationRepo(gdb)

	id := uuid.New()
	patch := map[string]any{"theme": "dark"}

	updateSQL := `UPDATE "organizations" SET "config"=config || $1::jsonb,"updated_at"=$2 WHERE id = $3`
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(updateSQL)).
		WithArgs(`{"theme":"dark"}`, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// FindByID follow-up to return the reloaded org.
	findSQL := `SELECT * FROM "organizations" WHERE id = $1 ORDER BY "organizations"."id" LIMIT $2`
	mock.ExpectQuery(regexp.QuoteMeta(findSQL)).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows(orgColumns).
			AddRow(orgRow(entities.Organization{ID: id, Name: "Acme", Slug: "acme", Config: []byte(`{"theme":"dark"}`)})...))

	got, err := repo.UpdateConfig(context.Background(), id, patch)
	require.NoError(t, err)
	assert.Contains(t, string(got.Config), "dark")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOrganizationRepo_UpdateConfig_NotFound(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewOrganizationRepo(gdb)

	id := uuid.New()
	updateSQL := `UPDATE "organizations" SET "config"=config || $1::jsonb,"updated_at"=$2 WHERE id = $3`
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(updateSQL)).
		WithArgs(`{"k":"v"}`, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	_, err := repo.UpdateConfig(context.Background(), id, map[string]any{"k": "v"})
	assert.ErrorIs(t, err, providers.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}
