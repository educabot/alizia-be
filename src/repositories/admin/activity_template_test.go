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

var activityColumns = []string{
	"id", "organization_id", "moment", "name", "description", "duration_minutes",
	"created_at", "updated_at",
}

func activityRow(a entities.ActivityTemplate) []driver.Value {
	return []driver.Value{
		a.ID, a.OrganizationID, string(a.Moment), a.Name, a.Description, a.DurationMinutes,
		a.CreatedAt, a.UpdatedAt,
	}
}

func TestActivityTemplateRepo_CreateActivity(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewActivityTemplateRepo(gdb)

	orgID := uuid.New()
	desc := "brainstorm"
	dur := 15
	a := &entities.ActivityTemplate{
		OrganizationID: orgID, Moment: entities.MomentApertura,
		Name: "Lluvia de ideas", Description: &desc, DurationMinutes: &dur,
	}

	sql := `INSERT INTO "activities" ("organization_id","moment","name","description","duration_minutes","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, entities.MomentApertura, "Lluvia de ideas", &desc, &dur, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))
	mock.ExpectCommit()

	id, err := repo.CreateActivity(context.Background(), a)
	require.NoError(t, err)
	assert.Equal(t, int64(1), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestActivityTemplateRepo_ListActivities_AllMoments(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewActivityTemplateRepo(gdb)

	orgID := uuid.New()
	sql := `SELECT * FROM "activities" WHERE organization_id = $1 ORDER BY moment, name LIMIT $2`
	rows := sqlmock.NewRows(activityColumns).
		AddRow(activityRow(entities.ActivityTemplate{ID: 1, OrganizationID: orgID, Moment: entities.MomentApertura, Name: "A"})...)
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, 51). // default 50 + 1
		WillReturnRows(rows)

	items, more, err := repo.ListActivities(context.Background(), orgID, nil, providers.Pagination{})
	require.NoError(t, err)
	assert.Len(t, items, 1)
	assert.False(t, more)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestActivityTemplateRepo_ListActivities_FilterMoment_PagesMore(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewActivityTemplateRepo(gdb)

	orgID := uuid.New()
	moment := entities.MomentDesarrollo

	sql := `SELECT * FROM "activities" WHERE organization_id = $1 AND moment = $2 ORDER BY moment, name LIMIT $3`
	// 3 rows for limit=2 to trigger `more=true`.
	rows := sqlmock.NewRows(activityColumns).
		AddRow(activityRow(entities.ActivityTemplate{ID: 1, OrganizationID: orgID, Moment: moment, Name: "A"})...).
		AddRow(activityRow(entities.ActivityTemplate{ID: 2, OrganizationID: orgID, Moment: moment, Name: "B"})...).
		AddRow(activityRow(entities.ActivityTemplate{ID: 3, OrganizationID: orgID, Moment: moment, Name: "C"})...)
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, moment, 3).
		WillReturnRows(rows)

	items, more, err := repo.ListActivities(
		context.Background(), orgID, &moment,
		providers.Pagination{Limit: 2, Offset: 0},
	)
	require.NoError(t, err)
	assert.Len(t, items, 2)
	assert.True(t, more)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestActivityTemplateRepo_CountByMoment(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewActivityTemplateRepo(gdb)

	orgID := uuid.New()
	sql := `SELECT count(*) FROM "activities" WHERE organization_id = $1 AND moment = $2`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, entities.MomentCierre).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(7)))

	n, err := repo.CountByMoment(context.Background(), orgID, entities.MomentCierre)
	require.NoError(t, err)
	assert.Equal(t, int64(7), n)
	require.NoError(t, mock.ExpectationsWereMet())
}
