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
	"github.com/educabot/alizia-be/src/repositories/admin"
	"github.com/educabot/alizia-be/src/utils/dbtest"
)

var subjectColumns = []string{
	"id", "organization_id", "area_id", "name", "description", "created_at", "updated_at",
}

func subjectRow(s entities.Subject) []driver.Value {
	return []driver.Value{
		s.ID, s.OrganizationID, s.AreaID, s.Name, s.Description, s.CreatedAt, s.UpdatedAt,
	}
}

func TestSubjectRepo_CreateSubject(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewSubjectRepo(gdb)

	orgID := uuid.New()
	desc := "d"
	s := &entities.Subject{OrganizationID: orgID, AreaID: 5, Name: "Maths", Description: &desc}

	sql := `INSERT INTO "subjects" ("organization_id","area_id","name","description","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, int64(5), "Maths", &desc, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(11)))
	mock.ExpectCommit()

	id, err := repo.CreateSubject(context.Background(), s)
	require.NoError(t, err)
	assert.Equal(t, int64(11), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectRepo_ListSubjectsByArea(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewSubjectRepo(gdb)

	orgID := uuid.New()
	sql := `SELECT * FROM "subjects" WHERE organization_id = $1 AND area_id = $2 ORDER BY name ASC LIMIT $3`
	rows := sqlmock.NewRows(subjectColumns).
		AddRow(subjectRow(entities.Subject{ID: 1, OrganizationID: orgID, AreaID: 5, Name: "A"})...)
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, int64(5), 500).
		WillReturnRows(rows)

	items, err := repo.ListSubjectsByArea(context.Background(), orgID, 5)
	require.NoError(t, err)
	assert.Len(t, items, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectRepo_ListSubjectsByOrg_NoAreaFilter(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewSubjectRepo(gdb)

	orgID := uuid.New()
	sql := `SELECT * FROM "subjects" WHERE organization_id = $1 ORDER BY name ASC LIMIT $2`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, 500).
		WillReturnRows(sqlmock.NewRows(subjectColumns))

	items, err := repo.ListSubjectsByOrg(context.Background(), orgID, nil)
	require.NoError(t, err)
	assert.Empty(t, items)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSubjectRepo_ListSubjectsByOrg_WithAreaFilter(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewSubjectRepo(gdb)

	orgID := uuid.New()
	areaID := int64(3)
	sql := `SELECT * FROM "subjects" WHERE organization_id = $1 AND area_id = $2 ORDER BY name ASC LIMIT $3`
	rows := sqlmock.NewRows(subjectColumns).
		AddRow(subjectRow(entities.Subject{ID: 1, OrganizationID: orgID, AreaID: 3, Name: "Física"})...)
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, areaID, 500).
		WillReturnRows(rows)

	items, err := repo.ListSubjectsByOrg(context.Background(), orgID, &areaID)
	require.NoError(t, err)
	assert.Len(t, items, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}
