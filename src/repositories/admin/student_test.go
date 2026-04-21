package admin_test

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/repositories/admin"
	"github.com/educabot/alizia-be/src/utils/dbtest"
)

var studentColumns = []string{"id", "course_id", "name", "created_at", "updated_at"}

func studentRow(s entities.Student) []driver.Value {
	return []driver.Value{s.ID, s.CourseID, s.Name, s.CreatedAt, s.UpdatedAt}
}

func TestStudentRepo_CreateStudent(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewStudentRepo(gdb)

	s := &entities.Student{CourseID: 9, Name: "Juan"}
	sql := `INSERT INTO "students" ("course_id","name","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(9), "Juan", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(77)))
	mock.ExpectCommit()

	id, err := repo.CreateStudent(context.Background(), s)
	require.NoError(t, err)
	assert.Equal(t, int64(77), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStudentRepo_ListByCourse(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewStudentRepo(gdb)

	now := time.Now()
	sql := `SELECT * FROM "students" WHERE course_id = $1 LIMIT $2`
	rows := sqlmock.NewRows(studentColumns).
		AddRow(studentRow(entities.Student{ID: 1, CourseID: 9, Name: "A", TimeTrackedEntity: entities.TimeTrackedEntity{CreatedAt: now, UpdatedAt: now}})...).
		AddRow(studentRow(entities.Student{ID: 2, CourseID: 9, Name: "B", TimeTrackedEntity: entities.TimeTrackedEntity{CreatedAt: now, UpdatedAt: now}})...)
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(9), 500). // boundedListCap
		WillReturnRows(rows)

	items, err := repo.ListByCourse(context.Background(), 9)
	require.NoError(t, err)
	assert.Len(t, items, 2)
	require.NoError(t, mock.ExpectationsWereMet())
}
