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

var courseColumns = []string{"id", "organization_id", "name", "created_at", "updated_at"}

func courseRow(c entities.Course) []driver.Value {
	return []driver.Value{c.ID, c.OrganizationID, c.Name, c.CreatedAt, c.UpdatedAt}
}

func TestCourseRepo_CreateCourse(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewCourseRepo(gdb)

	orgID := uuid.New()
	c := &entities.Course{OrganizationID: orgID, Name: "6° A"}

	sql := `INSERT INTO "courses" ("organization_id","name","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, "6° A", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(3)))
	mock.ExpectCommit()

	id, err := repo.CreateCourse(context.Background(), c)
	require.NoError(t, err)
	assert.Equal(t, int64(3), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCourseRepo_ListCourses(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewCourseRepo(gdb)

	orgID := uuid.New()
	sql := `SELECT * FROM "courses" WHERE organization_id = $1 ORDER BY name ASC LIMIT $2`
	rows := sqlmock.NewRows(courseColumns).
		AddRow(courseRow(entities.Course{ID: 1, OrganizationID: orgID, Name: "A"})...).
		AddRow(courseRow(entities.Course{ID: 2, OrganizationID: orgID, Name: "B"})...)
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, 51).
		WillReturnRows(rows)

	items, more, err := repo.ListCourses(context.Background(), orgID, providers.Pagination{})
	require.NoError(t, err)
	assert.Len(t, items, 2)
	assert.False(t, more)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCourseRepo_GetCourse_NotFound(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewCourseRepo(gdb)

	orgID := uuid.New()
	sql := `SELECT * FROM "courses" WHERE organization_id = $1 AND id = $2 ORDER BY "courses"."id" LIMIT $3`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, int64(99), 1).
		WillReturnRows(sqlmock.NewRows(courseColumns))

	_, err := repo.GetCourse(context.Background(), orgID, 99)
	assert.ErrorIs(t, err, providers.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

// TestCourseRepo_GetCourse_Found verifies the happy path with empty preload
// children. We return 1 course row from the main query and empty result sets
// for the Students and CourseSubjects preloads — GORM still fires those
// queries to honour the Preload chain; downstream subject/teacher preloads are
// skipped because the CourseSubjects parent is empty.
func TestCourseRepo_GetCourse_Found(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewCourseRepo(gdb)

	orgID := uuid.New()

	mainSQL := `SELECT * FROM "courses" WHERE organization_id = $1 AND id = $2 ORDER BY "courses"."id" LIMIT $3`
	mock.ExpectQuery(regexp.QuoteMeta(mainSQL)).
		WithArgs(orgID, int64(42), 1).
		WillReturnRows(sqlmock.NewRows(courseColumns).
			AddRow(courseRow(entities.Course{ID: 42, OrganizationID: orgID, Name: "Found"})...))

	// Preload order is decided by GORM, not by the `.Preload()` call order.
	// Use unordered matching so the assertion survives GORM-version reshuffles.
	mock.MatchExpectationsInOrder(false)

	studentsSQL := `SELECT * FROM "students" WHERE "students"."course_id" = $1`
	mock.ExpectQuery(regexp.QuoteMeta(studentsSQL)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "course_id", "name", "created_at", "updated_at"}))

	csSQL := `SELECT * FROM "course_subjects" WHERE "course_subjects"."course_id" = $1`
	mock.ExpectQuery(regexp.QuoteMeta(csSQL)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "organization_id", "course_id", "subject_id", "teacher_id",
			"school_year", "start_date", "end_date", "created_at", "updated_at",
		}))

	// Subject/Teacher preloads skipped — empty parent slice.

	got, err := repo.GetCourse(context.Background(), orgID, 42)
	require.NoError(t, err)
	assert.Equal(t, "Found", got.Name)
	require.NoError(t, mock.ExpectationsWereMet())
}
