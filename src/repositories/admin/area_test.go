package admin_test

import (
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/repositories/admin"
	"github.com/educabot/alizia-be/src/utils/dbtest"
)

var areaColumns = []string{"id", "organization_id", "name", "description", "created_at", "updated_at"}

func areaRow(a entities.Area) []driver.Value {
	return []driver.Value{a.ID, a.OrganizationID, a.Name, a.Description, a.CreatedAt, a.UpdatedAt}
}

func TestAreaRepo_CreateArea(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaRepo(gdb)

	orgID := uuid.New()
	a := &entities.Area{OrganizationID: orgID, Name: "Ciencias"}

	sql := `INSERT INTO "areas" ("organization_id","name","description","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, "Ciencias", nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))
	mock.ExpectCommit()

	id, err := repo.CreateArea(context.Background(), a)
	require.NoError(t, err)
	assert.Equal(t, int64(1), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaRepo_CreateArea_DuplicateUnique(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaRepo(gdb)

	orgID := uuid.New()
	a := &entities.Area{OrganizationID: orgID, Name: "Ciencias"}

	pgErr := &pgconn.PgError{Code: pgerrcode.UniqueViolation}

	sql := `INSERT INTO "areas" ("organization_id","name","description","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, "Ciencias", nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(pgErr)
	mock.ExpectRollback()

	_, err := repo.CreateArea(context.Background(), a)
	assert.ErrorIs(t, err, providers.ErrConflict)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaRepo_GetArea_Found(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaRepo(gdb)
	mock.MatchExpectationsInOrder(false)

	orgID := uuid.New()

	mainSQL := `SELECT * FROM "areas" WHERE organization_id = $1 AND id = $2 ORDER BY "areas"."id" LIMIT $3`
	mock.ExpectQuery(regexp.QuoteMeta(mainSQL)).
		WithArgs(orgID, int64(5), 1).
		WillReturnRows(sqlmock.NewRows(areaColumns).
			AddRow(areaRow(entities.Area{ID: 5, OrganizationID: orgID, Name: "Ciencias"})...))

	// Preload Subjects: WHERE area_id = $1.
	subjectsSQL := `SELECT * FROM "subjects" WHERE "subjects"."area_id" = $1`
	mock.ExpectQuery(regexp.QuoteMeta(subjectsSQL)).
		WithArgs(int64(5)).
		WillReturnRows(sqlmock.NewRows(subjectColumns))

	// Preload Coordinators: WHERE area_id = $1. Empty result skips nested User preload.
	coordSQL := `SELECT * FROM "area_coordinators" WHERE "area_coordinators"."area_id" = $1`
	mock.ExpectQuery(regexp.QuoteMeta(coordSQL)).
		WithArgs(int64(5)).
		WillReturnRows(sqlmock.NewRows(areaCoordColumns))

	got, err := repo.GetArea(context.Background(), orgID, 5)
	require.NoError(t, err)
	assert.Equal(t, "Ciencias", got.Name)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaRepo_GetArea_NotFound(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaRepo(gdb)

	orgID := uuid.New()
	mainSQL := `SELECT * FROM "areas" WHERE organization_id = $1 AND id = $2 ORDER BY "areas"."id" LIMIT $3`
	mock.ExpectQuery(regexp.QuoteMeta(mainSQL)).
		WithArgs(orgID, int64(99), 1).
		WillReturnRows(sqlmock.NewRows(areaColumns))

	_, err := repo.GetArea(context.Background(), orgID, 99)
	assert.ErrorIs(t, err, providers.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaRepo_ListAreas(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaRepo(gdb)

	orgID := uuid.New()
	mainSQL := `SELECT * FROM "areas" WHERE organization_id = $1 ORDER BY name ASC LIMIT $2`
	mock.ExpectQuery(regexp.QuoteMeta(mainSQL)).
		WithArgs(orgID, 500).
		WillReturnRows(sqlmock.NewRows(areaColumns))

	// Empty main result → preloads skipped.
	items, err := repo.ListAreas(context.Background(), orgID)
	require.NoError(t, err)
	assert.Empty(t, items)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaRepo_UpdateArea(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaRepo(gdb)

	orgID := uuid.New()
	desc := "updated"
	area := &entities.Area{ID: 5, OrganizationID: orgID, Name: "Humanidades", Description: &desc}

	// Map keys serialized alphabetically: description, name (+ updated_at appended).
	sql := `UPDATE "areas" SET "description"=$1,"name"=$2,"updated_at"=$3 WHERE organization_id = $4 AND id = $5`
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sql)).
		WithArgs(&desc, "Humanidades", sqlmock.AnyArg(), orgID, int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.UpdateArea(context.Background(), area)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaRepo_UpdateArea_DuplicateUnique(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaRepo(gdb)

	orgID := uuid.New()
	area := &entities.Area{ID: 5, OrganizationID: orgID, Name: "Dup"}

	pgErr := &pgconn.PgError{Code: pgerrcode.UniqueViolation}
	sql := `UPDATE "areas" SET "description"=$1,"name"=$2,"updated_at"=$3 WHERE organization_id = $4 AND id = $5`
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sql)).
		WithArgs(nil, "Dup", sqlmock.AnyArg(), orgID, int64(5)).
		WillReturnError(pgErr)
	mock.ExpectRollback()

	err := repo.UpdateArea(context.Background(), area)
	assert.ErrorIs(t, err, providers.ErrConflict)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaRepo_CountDependencies(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaRepo(gdb)

	orgID := uuid.New()
	sql := `SELECT count(*) FROM "subjects" WHERE organization_id = $1 AND area_id = $2`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(3)))

	deps, err := repo.CountDependencies(context.Background(), orgID, 5)
	require.NoError(t, err)
	assert.Equal(t, int64(3), deps.Subjects)
	assert.False(t, deps.IsEmpty())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaRepo_DeleteArea_Success(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaRepo(gdb)

	orgID := uuid.New()

	deleteCoordsSQL := `DELETE FROM "area_coordinators" WHERE area_id = $1`
	deleteAreaSQL := `DELETE FROM "areas" WHERE organization_id = $1 AND id = $2`

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(deleteCoordsSQL)).
		WithArgs(int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(regexp.QuoteMeta(deleteAreaSQL)).
		WithArgs(orgID, int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.DeleteArea(context.Background(), orgID, 5)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

// TestAreaRepo_DeleteArea_NotFound verifies that when the area DELETE affects
// 0 rows the repo returns ErrNotFound and rolls back the transaction.
func TestAreaRepo_DeleteArea_NotFound(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaRepo(gdb)

	orgID := uuid.New()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "area_coordinators" WHERE area_id = $1`)).
		WithArgs(int64(99)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "areas" WHERE organization_id = $1 AND id = $2`)).
		WithArgs(orgID, int64(99)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := repo.DeleteArea(context.Background(), orgID, 99)
	assert.ErrorIs(t, err, providers.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaRepo_DeleteArea_CoordinatorsDeleteFails(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaRepo(gdb)

	orgID := uuid.New()
	boom := errors.New("FK violation")
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "area_coordinators" WHERE area_id = $1`)).
		WithArgs(int64(5)).
		WillReturnError(boom)
	mock.ExpectRollback()

	err := repo.DeleteArea(context.Background(), orgID, 5)
	assert.ErrorIs(t, err, boom)
	require.NoError(t, mock.ExpectationsWereMet())
}
