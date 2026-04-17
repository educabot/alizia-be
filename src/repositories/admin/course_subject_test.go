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

var courseSubjectColumns = []string{
	"id", "organization_id", "course_id", "subject_id", "teacher_id",
	"school_year", "start_date", "end_date", "created_at", "updated_at",
}

func courseSubjectRow(cs entities.CourseSubject) []driver.Value {
	return []driver.Value{
		cs.ID, cs.OrganizationID, cs.CourseID, cs.SubjectID, cs.TeacherID,
		cs.SchoolYear, cs.StartDate, cs.EndDate, cs.CreatedAt, cs.UpdatedAt,
	}
}

func TestCourseSubjectRepo_CreateCourseSubject(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewCourseSubjectRepo(gdb)

	orgID := uuid.New()
	cs := &entities.CourseSubject{
		OrganizationID: orgID, CourseID: 1, SubjectID: 2, TeacherID: 3, SchoolYear: 2026,
	}

	sql := `INSERT INTO "course_subjects" ("organization_id","course_id","subject_id","teacher_id","school_year","start_date","end_date","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, int64(1), int64(2), int64(3), 2026, nil, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(99)))
	mock.ExpectCommit()

	id, err := repo.CreateCourseSubject(context.Background(), cs)
	require.NoError(t, err)
	assert.Equal(t, int64(99), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

// TestCourseSubjectRepo_CreateCourseSubject_DuplicateUnique verifies the repo
// translates a Postgres unique_violation into the domain ErrConflict.
func TestCourseSubjectRepo_CreateCourseSubject_DuplicateUnique(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewCourseSubjectRepo(gdb)

	orgID := uuid.New()
	cs := &entities.CourseSubject{OrganizationID: orgID, CourseID: 1, SubjectID: 2, TeacherID: 3, SchoolYear: 2026}

	pgErr := &pgconn.PgError{Code: pgerrcode.UniqueViolation, Message: "duplicate key"}

	sql := `INSERT INTO "course_subjects" ("organization_id","course_id","subject_id","teacher_id","school_year","start_date","end_date","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, int64(1), int64(2), int64(3), 2026, nil, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(pgErr)
	mock.ExpectRollback()

	_, err := repo.CreateCourseSubject(context.Background(), cs)
	assert.ErrorIs(t, err, providers.ErrConflict)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCourseSubjectRepo_CreateCourseSubject_UnrelatedError(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewCourseSubjectRepo(gdb)

	orgID := uuid.New()
	cs := &entities.CourseSubject{OrganizationID: orgID, CourseID: 1, SubjectID: 2, TeacherID: 3, SchoolYear: 2026}

	boom := errors.New("network")
	sql := `INSERT INTO "course_subjects" ("organization_id","course_id","subject_id","teacher_id","school_year","start_date","end_date","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, int64(1), int64(2), int64(3), 2026, nil, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(boom)
	mock.ExpectRollback()

	_, err := repo.CreateCourseSubject(context.Background(), cs)
	assert.ErrorIs(t, err, boom)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCourseSubjectRepo_GetCourseSubject_NotFound(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewCourseSubjectRepo(gdb)

	orgID := uuid.New()
	sql := `SELECT * FROM "course_subjects" WHERE organization_id = $1 AND id = $2 ORDER BY "course_subjects"."id" LIMIT $3`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, int64(99), 1).
		WillReturnRows(sqlmock.NewRows(courseSubjectColumns))

	_, err := repo.GetCourseSubject(context.Background(), orgID, 99)
	assert.ErrorIs(t, err, providers.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCourseSubjectRepo_ListByCourse(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewCourseSubjectRepo(gdb)

	sql := `SELECT * FROM "course_subjects" WHERE course_id = $1 LIMIT $2`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(5), 500).
		WillReturnRows(sqlmock.NewRows(courseSubjectColumns))

	items, err := repo.ListByCourse(context.Background(), 5)
	require.NoError(t, err)
	assert.Empty(t, items)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCourseSubjectRepo_ListCourseSubjects_NoFilter(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewCourseSubjectRepo(gdb)

	orgID := uuid.New()
	sql := `SELECT * FROM "course_subjects" WHERE organization_id = $1 ORDER BY course_id ASC, subject_id ASC LIMIT $2`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, 500).
		WillReturnRows(sqlmock.NewRows(courseSubjectColumns))

	items, err := repo.ListCourseSubjects(context.Background(), orgID, providers.CourseSubjectFilter{})
	require.NoError(t, err)
	assert.Empty(t, items)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCourseSubjectRepo_ListCourseSubjects_AllFilters(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewCourseSubjectRepo(gdb)

	// GORM may reorder preload queries; match by content instead of order.
	mock.MatchExpectationsInOrder(false)

	orgID := uuid.New()
	courseID, subjectID, teacherID := int64(1), int64(2), int64(3)
	filter := providers.CourseSubjectFilter{CourseID: &courseID, SubjectID: &subjectID, TeacherID: &teacherID}

	sql := `SELECT * FROM "course_subjects" WHERE organization_id = $1 AND course_id = $2 AND subject_id = $3 AND teacher_id = $4 ORDER BY course_id ASC, subject_id ASC LIMIT $5`
	rows := sqlmock.NewRows(courseSubjectColumns).
		AddRow(courseSubjectRow(entities.CourseSubject{
			ID: 10, OrganizationID: orgID, CourseID: 1, SubjectID: 2, TeacherID: 3, SchoolYear: 2026,
		})...)
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, courseID, subjectID, teacherID, 500).
		WillReturnRows(rows)

	// Preload Subject fires because 1 row was returned (subject_id=2).
	subjectSQL := `SELECT * FROM "subjects" WHERE "subjects"."id" = $1`
	mock.ExpectQuery(regexp.QuoteMeta(subjectSQL)).
		WithArgs(subjectID).
		WillReturnRows(sqlmock.NewRows(subjectColumns))

	// Preload Teacher fires with teacher_id=3.
	teacherSQL := `SELECT * FROM "users" WHERE "users"."id" = $1`
	mock.ExpectQuery(regexp.QuoteMeta(teacherSQL)).
		WithArgs(teacherID).
		WillReturnRows(sqlmock.NewRows(userColumns))

	items, err := repo.ListCourseSubjects(context.Background(), orgID, filter)
	require.NoError(t, err)
	assert.Len(t, items, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}
