package admin_test

import (
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/repositories/admin"
	"github.com/educabot/alizia-be/src/utils/dbtest"
)

// topicColumns lists the columns sqlmock rows must expose so GORM can scan into
// entities.Topic. Keep aligned with the struct in src/core/entities/topic.go.
var topicColumns = []string{
	"id", "organization_id", "parent_id", "name", "description", "level",
	"created_at", "updated_at",
}

func topicRow(t entities.Topic) []driver.Value {
	return []driver.Value{
		t.ID, t.OrganizationID, t.ParentID, t.Name, t.Description, t.Level,
		t.CreatedAt, t.UpdatedAt,
	}
}

func TestTopicRepo_GetTopicByID_Found(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTopicRepo(gdb)

	orgID := uuid.New()
	now := time.Now()
	want := entities.Topic{
		ID: 42, OrganizationID: orgID, Name: "Álgebra", Level: 2,
		TimeTrackedEntity: entities.TimeTrackedEntity{CreatedAt: now, UpdatedAt: now},
	}

	sql := `SELECT * FROM "topics" WHERE organization_id = $1 AND id = $2 ORDER BY "topics"."id" LIMIT $3`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, int64(42), 1).
		WillReturnRows(sqlmock.NewRows(topicColumns).AddRow(topicRow(want)...))

	got, err := repo.GetTopicByID(context.Background(), orgID, 42)
	require.NoError(t, err)
	assert.Equal(t, "Álgebra", got.Name)
	assert.Equal(t, int64(42), got.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTopicRepo_GetTopicByID_NotFound(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTopicRepo(gdb)

	orgID := uuid.New()
	sql := `SELECT * FROM "topics" WHERE organization_id = $1 AND id = $2 ORDER BY "topics"."id" LIMIT $3`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, int64(99), 1).
		WillReturnRows(sqlmock.NewRows(topicColumns))

	_, err := repo.GetTopicByID(context.Background(), orgID, 99)
	assert.ErrorIs(t, err, providers.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTopicRepo_GetTopicTree_BuildsHierarchy(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTopicRepo(gdb)

	orgID := uuid.New()
	pid1 := int64(1)
	pid2 := int64(2)

	// Flat rows sorted by (level ASC, name ASC) — buildTree assembles from this.
	flat := []entities.Topic{
		{ID: 1, OrganizationID: orgID, Name: "Root", Level: 1},
		{ID: 2, OrganizationID: orgID, ParentID: &pid1, Name: "Child", Level: 2},
		{ID: 3, OrganizationID: orgID, ParentID: &pid2, Name: "Grandchild", Level: 3},
	}

	sql := `SELECT * FROM "topics" WHERE organization_id = $1 ORDER BY level ASC, name ASC LIMIT $2`
	rows := sqlmock.NewRows(topicColumns)
	for _, tp := range flat {
		rows.AddRow(topicRow(tp)...)
	}
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, 1000). // maxTopicTreeSize
		WillReturnRows(rows)

	tree, err := repo.GetTopicTree(context.Background(), orgID)
	require.NoError(t, err)
	require.Len(t, tree, 1, "one root")
	assert.Equal(t, "Root", tree[0].Name)
	require.Len(t, tree[0].Children, 1, "root has one child")
	assert.Equal(t, "Child", tree[0].Children[0].Name)
	require.Len(t, tree[0].Children[0].Children, 1, "child has grandchild — multi-depth build")
	assert.Equal(t, "Grandchild", tree[0].Children[0].Children[0].Name)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTopicRepo_GetTopicsByLevel_PaginationMore(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTopicRepo(gdb)

	orgID := uuid.New()
	page := providers.Pagination{Limit: 2, Offset: 0}

	// Repo fetches limit+1 to detect "more"; return 3 rows for limit=2.
	sql := `SELECT * FROM "topics" WHERE organization_id = $1 AND level = $2 ORDER BY name ASC LIMIT $3`
	rows := sqlmock.NewRows(topicColumns).
		AddRow(topicRow(entities.Topic{ID: 1, OrganizationID: orgID, Name: "A", Level: 2})...).
		AddRow(topicRow(entities.Topic{ID: 2, OrganizationID: orgID, Name: "B", Level: 2})...).
		AddRow(topicRow(entities.Topic{ID: 3, OrganizationID: orgID, Name: "C", Level: 2})...)
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, 2, 3). // limit+1 = 3
		WillReturnRows(rows)

	items, more, err := repo.GetTopicsByLevel(context.Background(), orgID, 2, page)
	require.NoError(t, err)
	assert.Len(t, items, 2, "trim to requested limit")
	assert.True(t, more, "3 rows for limit=2 means another page exists")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTopicRepo_GetTopicsByLevel_NoMore(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTopicRepo(gdb)

	orgID := uuid.New()
	page := providers.Pagination{Limit: 10, Offset: 0}

	sql := `SELECT * FROM "topics" WHERE organization_id = $1 AND level = $2 ORDER BY name ASC LIMIT $3`
	rows := sqlmock.NewRows(topicColumns).
		AddRow(topicRow(entities.Topic{ID: 1, OrganizationID: orgID, Name: "Only", Level: 2})...)
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, 2, 11).
		WillReturnRows(rows)

	items, more, err := repo.GetTopicsByLevel(context.Background(), orgID, 2, page)
	require.NoError(t, err)
	assert.Len(t, items, 1)
	assert.False(t, more)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTopicRepo_GetTopicsByParent_RootParent(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTopicRepo(gdb)

	orgID := uuid.New()

	// nil parent → WHERE parent_id IS NULL (not `= NULL`).
	sql := `SELECT * FROM "topics" WHERE organization_id = $1 AND parent_id IS NULL ORDER BY name ASC LIMIT $2`
	rows := sqlmock.NewRows(topicColumns).
		AddRow(topicRow(entities.Topic{ID: 1, OrganizationID: orgID, Name: "R1", Level: 1})...)
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, 51). // default limit 50 + 1
		WillReturnRows(rows)

	items, more, err := repo.GetTopicsByParent(context.Background(), orgID, nil, providers.Pagination{})
	require.NoError(t, err)
	assert.Len(t, items, 1)
	assert.False(t, more)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTopicRepo_GetTopicsByParent_SpecificParent(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTopicRepo(gdb)

	orgID := uuid.New()
	parent := int64(7)

	sql := `SELECT * FROM "topics" WHERE organization_id = $1 AND parent_id = $2 ORDER BY name ASC LIMIT $3`
	rows := sqlmock.NewRows(topicColumns).
		AddRow(topicRow(entities.Topic{ID: 10, OrganizationID: orgID, ParentID: &parent, Name: "Kid", Level: 2})...)
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, parent, 51).
		WillReturnRows(rows)

	items, _, err := repo.GetTopicsByParent(context.Background(), orgID, &parent, providers.Pagination{})
	require.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, &parent, items[0].ParentID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTopicRepo_ListAllTopics_NoCap(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTopicRepo(gdb)

	orgID := uuid.New()

	// Crucially: no LIMIT clause — cycle detection needs the full graph.
	sql := `SELECT * FROM "topics" WHERE organization_id = $1 ORDER BY level ASC, name ASC`
	rows := sqlmock.NewRows(topicColumns).
		AddRow(topicRow(entities.Topic{ID: 1, OrganizationID: orgID, Name: "A", Level: 1})...).
		AddRow(topicRow(entities.Topic{ID: 2, OrganizationID: orgID, Name: "B", Level: 1})...)
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID).
		WillReturnRows(rows)

	items, err := repo.ListAllTopics(context.Background(), orgID)
	require.NoError(t, err)
	assert.Len(t, items, 2)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTopicRepo_UpdateTopic(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTopicRepo(gdb)

	orgID := uuid.New()
	desc := "updated desc"
	pid := int64(5)
	topic := &entities.Topic{
		ID: 10, OrganizationID: orgID, ParentID: &pid,
		Name: "New Name", Description: &desc, Level: 3,
	}

	// GORM sorts map[string]any keys alphabetically:
	// description, level, name, parent_id, (updated_at appended).
	sql := `UPDATE "topics" SET "description"=$1,"level"=$2,"name"=$3,"parent_id"=$4,"updated_at"=$5 WHERE organization_id = $6 AND id = $7`
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sql)).
		WithArgs(&desc, 3, "New Name", &pid, sqlmock.AnyArg(), orgID, int64(10)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.UpdateTopic(context.Background(), topic)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTopicRepo_UpdateTopicLevels_Empty_NoQuery(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTopicRepo(gdb)

	// Empty map short-circuits: no BEGIN, no UPDATE, no COMMIT.
	err := repo.UpdateTopicLevels(context.Background(), uuid.New(), map[int64]int{})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTopicRepo_UpdateTopicLevels_TransactionalUpdates(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTopicRepo(gdb)

	orgID := uuid.New()
	// Single-entry map keeps the test deterministic (map iteration order is random).
	levels := map[int64]int{7: 3}

	sql := `UPDATE "topics" SET "level"=$1,"updated_at"=$2 WHERE organization_id = $3 AND id = $4`
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sql)).
		WithArgs(3, sqlmock.AnyArg(), orgID, int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.UpdateTopicLevels(context.Background(), orgID, levels)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTopicRepo_CreateTopic_ReturnsID(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTopicRepo(gdb)

	orgID := uuid.New()
	desc := "desc"
	topic := &entities.Topic{
		OrganizationID: orgID, Name: "New", Description: &desc, Level: 1,
	}

	sql := `INSERT INTO "topics" ("organization_id","parent_id","name","description","level","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, nil, "New", &desc, 1, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(123)))
	mock.ExpectCommit()

	id, err := repo.CreateTopic(context.Background(), topic)
	require.NoError(t, err)
	assert.Equal(t, int64(123), id)
	assert.Equal(t, int64(123), topic.ID, "repo should backfill ID on the passed pointer")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTopicRepo_GetTopicByID_DBError(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTopicRepo(gdb)

	orgID := uuid.New()
	boom := errors.New("connection closed")
	sql := `SELECT * FROM "topics" WHERE organization_id = $1 AND id = $2 ORDER BY "topics"."id" LIMIT $3`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(orgID, int64(1), 1).
		WillReturnError(boom)

	_, err := repo.GetTopicByID(context.Background(), orgID, 1)
	assert.ErrorIs(t, err, boom, "unrelated DB errors must propagate as-is")
	require.NoError(t, mock.ExpectationsWereMet())
}
