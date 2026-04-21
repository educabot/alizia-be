package admin_test

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/repositories/admin"
	"github.com/educabot/alizia-be/src/utils/dbtest"
)

var areaCoordColumns = []string{"id", "area_id", "user_id", "created_at"}

func areaCoordRow(ac entities.AreaCoordinator) []driver.Value {
	return []driver.Value{ac.ID, ac.AreaID, ac.UserID, ac.CreatedAt}
}

func TestAreaCoordinatorRepo_Assign(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaCoordinatorRepo(gdb)

	sql := `INSERT INTO "area_coordinators" ("area_id","user_id","created_at") VALUES ($1,$2,$3) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(5), int64(9), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))
	mock.ExpectCommit()

	ac, err := repo.Assign(context.Background(), 5, 9)
	require.NoError(t, err)
	assert.Equal(t, int64(1), ac.ID)
	assert.Equal(t, int64(5), ac.AreaID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaCoordinatorRepo_Assign_DuplicateUnique(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaCoordinatorRepo(gdb)

	pgErr := &pgconn.PgError{Code: pgerrcode.UniqueViolation, Message: "duplicate"}

	sql := `INSERT INTO "area_coordinators" ("area_id","user_id","created_at") VALUES ($1,$2,$3) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(5), int64(9), sqlmock.AnyArg()).
		WillReturnError(pgErr)
	mock.ExpectRollback()

	_, err := repo.Assign(context.Background(), 5, 9)
	assert.ErrorIs(t, err, providers.ErrConflict)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaCoordinatorRepo_Remove_Success(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaCoordinatorRepo(gdb)

	sql := `DELETE FROM "area_coordinators" WHERE area_id = $1 AND user_id = $2`
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sql)).
		WithArgs(int64(5), int64(9)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Remove(context.Background(), 5, 9)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaCoordinatorRepo_Remove_NotFound(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaCoordinatorRepo(gdb)

	sql := `DELETE FROM "area_coordinators" WHERE area_id = $1 AND user_id = $2`
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sql)).
		WithArgs(int64(5), int64(9)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Remove(context.Background(), 5, 9)
	assert.ErrorIs(t, err, providers.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaCoordinatorRepo_FindByAreaID(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaCoordinatorRepo(gdb)

	now := time.Now()
	sql := `SELECT * FROM "area_coordinators" WHERE area_id = $1`
	rows := sqlmock.NewRows(areaCoordColumns).
		AddRow(areaCoordRow(entities.AreaCoordinator{ID: 1, AreaID: 5, UserID: 10, CreatedAt: now})...).
		AddRow(areaCoordRow(entities.AreaCoordinator{ID: 2, AreaID: 5, UserID: 11, CreatedAt: now})...)
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(5)).
		WillReturnRows(rows)

	items, err := repo.FindByAreaID(context.Background(), 5)
	require.NoError(t, err)
	assert.Len(t, items, 2)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaCoordinatorRepo_FindByUserID(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaCoordinatorRepo(gdb)

	sql := `SELECT * FROM "area_coordinators" WHERE user_id = $1`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(9)).
		WillReturnRows(sqlmock.NewRows(areaCoordColumns).
			AddRow(areaCoordRow(entities.AreaCoordinator{ID: 1, AreaID: 5, UserID: 9})...))

	items, err := repo.FindByUserID(context.Background(), 9)
	require.NoError(t, err)
	assert.Len(t, items, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaCoordinatorRepo_IsCoordinator_True(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaCoordinatorRepo(gdb)

	sql := `SELECT count(*) FROM "area_coordinators" WHERE area_id = $1 AND user_id = $2`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(5), int64(9)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(1)))

	yes, err := repo.IsCoordinator(context.Background(), 5, 9)
	require.NoError(t, err)
	assert.True(t, yes)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAreaCoordinatorRepo_IsCoordinator_False(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewAreaCoordinatorRepo(gdb)

	sql := `SELECT count(*) FROM "area_coordinators" WHERE area_id = $1 AND user_id = $2`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(5), int64(9)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))

	yes, err := repo.IsCoordinator(context.Background(), 5, 9)
	require.NoError(t, err)
	assert.False(t, yes)
	require.NoError(t, mock.ExpectationsWereMet())
}
