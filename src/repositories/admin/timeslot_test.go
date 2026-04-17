package admin_test

import (
	"context"
	"database/sql/driver"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/repositories/admin"
	"github.com/educabot/alizia-be/src/repositories/admin/queries"
	"github.com/educabot/alizia-be/src/utils/dbtest"
)

// sharedClassNumbersSQL mirrors what sqlmock observes after pgx rebinds the
// two `?` placeholders in the embedded query to positional `$1`, `$2` params.
// We rebuild the expected string from the same source so the test fails fast
// if the .sql file drifts.
var sharedClassNumbersSQL = func() string {
	s := strings.Replace(queries.SharedClassNumbers, "?", "$1", 1)
	s = strings.Replace(s, "?", "$2", 1)
	return s
}()

var timeSlotColumns = []string{
	"id", "course_id", "day_of_week", "start_time", "end_time", "created_at", "updated_at",
}

func timeSlotRow(s entities.TimeSlot) []driver.Value {
	return []driver.Value{s.ID, s.CourseID, s.DayOfWeek, s.StartTime, s.EndTime, s.CreatedAt, s.UpdatedAt}
}

func TestTimeSlotRepo_CreateTimeSlot(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTimeSlotRepo(gdb)

	slot := &entities.TimeSlot{CourseID: 1, DayOfWeek: 1, StartTime: "08:00", EndTime: "09:00"}

	sql := `INSERT INTO "time_slots" ("course_id","day_of_week","start_time","end_time","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(1), 1, "08:00", "09:00", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(10)))
	mock.ExpectCommit()

	id, err := repo.CreateTimeSlot(context.Background(), slot)
	require.NoError(t, err)
	assert.Equal(t, int64(10), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

// TestTimeSlotRepo_ListByCourse_Empty exercises the happy path with no rows so
// GORM skips the chain of nested Preloads (Subjects → CourseSubject → Subject
// → Teacher). That keeps the assertion focused on the main query without
// duplicating the Preload fan-out already covered elsewhere.
func TestTimeSlotRepo_ListByCourse_Empty(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTimeSlotRepo(gdb)

	sql := `SELECT * FROM "time_slots" WHERE course_id = $1 ORDER BY day_of_week, start_time LIMIT $2`
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(1), 500).
		WillReturnRows(sqlmock.NewRows(timeSlotColumns))

	items, err := repo.ListByCourse(context.Background(), 1)
	require.NoError(t, err)
	assert.Empty(t, items)
	require.NoError(t, mock.ExpectationsWereMet())
}

// TestTimeSlotRepo_GetSharedClassNumbers_NoSlots verifies the early return when
// the CTE returns zero rows: no iteration over classes, empty result.
func TestTimeSlotRepo_GetSharedClassNumbers_NoSlots(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTimeSlotRepo(gdb)

	orgID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(sharedClassNumbersSQL)).
		WithArgs(int64(42), orgID).
		WillReturnRows(sqlmock.NewRows([]string{"weekly_position", "is_shared"}))

	shared, err := repo.GetSharedClassNumbers(context.Background(), orgID, 42, 10)
	require.NoError(t, err)
	assert.Empty(t, shared)
	require.NoError(t, mock.ExpectationsWereMet())
}

// TestTimeSlotRepo_GetSharedClassNumbers_ProjectsSharedFlag runs the projection
// loop over a 3-class-per-week schedule where only position 1 is shared, then
// walks 7 total classes. Expected shared class numbers: 2 and 5 (0-indexed
// positions 1, 1 over classesPerWeek=3).
func TestTimeSlotRepo_GetSharedClassNumbers_ProjectsSharedFlag(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewTimeSlotRepo(gdb)

	orgID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(sharedClassNumbersSQL)).
		WithArgs(int64(42), orgID).
		WillReturnRows(sqlmock.NewRows([]string{"weekly_position", "is_shared"}).
			AddRow(1, false).
			AddRow(2, true).
			AddRow(3, false))

	shared, err := repo.GetSharedClassNumbers(context.Background(), orgID, 42, 7)
	require.NoError(t, err)
	// classNum 1..7 → weekPos (classNum-1) % 3.
	// wp=0 false; wp=1 true (class 2); wp=2 false; wp=0 false (class 4);
	// wp=1 true (class 5); wp=2 false; wp=0 false.
	assert.Equal(t, []int{2, 5}, shared)
	require.NoError(t, mock.ExpectationsWereMet())
}
