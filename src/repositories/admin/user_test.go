package admin_test

import (
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/educabot/alizia-be/src/core/entities"
	"github.com/educabot/alizia-be/src/core/providers"
	"github.com/educabot/alizia-be/src/repositories/admin"
	"github.com/educabot/alizia-be/src/utils/dbtest"
)

// userColumns lists every column GORM expects to scan into entities.User,
// including the JSONB columns stored as []byte (profile_data) and the hidden
// password_hash. Shared with course_subject preload tests.
var userColumns = []string{
	"id", "organization_id", "email", "first_name", "last_name",
	"password_hash", "avatar_url", "onboarding_completed_at", "profile_data",
	"created_at", "updated_at",
}

func userRow(u entities.User) []driver.Value {
	return []driver.Value{
		u.ID, u.OrganizationID, u.Email, u.FirstName, u.LastName,
		u.PasswordHash, u.AvatarURL, u.OnboardingCompletedAt, []byte(u.ProfileData),
		u.CreatedAt, u.UpdatedAt,
	}
}

var userRoleColumns = []string{"id", "user_id", "role"}

func TestUserRepo_FindByID_Found(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewUserRepo(gdb)

	orgID := uuid.New()

	mainSQL := `SELECT * FROM "users" WHERE organization_id = $1 AND id = $2 ORDER BY "users"."id" LIMIT $3`
	mock.ExpectQuery(regexp.QuoteMeta(mainSQL)).
		WithArgs(orgID, int64(7), 1).
		WillReturnRows(sqlmock.NewRows(userColumns).
			AddRow(userRow(entities.User{
				ID: 7, OrganizationID: orgID, Email: "a@b.c",
				FirstName: "A", LastName: "B", ProfileData: []byte(`{}`),
			})...))

	// Preload Roles: WHERE user_id = $1 for the single returned user.
	rolesSQL := `SELECT * FROM "user_roles" WHERE "user_roles"."user_id" = $1`
	mock.ExpectQuery(regexp.QuoteMeta(rolesSQL)).
		WithArgs(int64(7)).
		WillReturnRows(sqlmock.NewRows(userRoleColumns).
			AddRow(int64(1), int64(7), string(entities.RoleCoordinator)))

	got, err := repo.FindByID(context.Background(), orgID, 7)
	require.NoError(t, err)
	assert.Equal(t, "a@b.c", got.Email)
	require.Len(t, got.Roles, 1)
	assert.Equal(t, entities.RoleCoordinator, got.Roles[0].Role)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_FindByID_NotFound(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewUserRepo(gdb)

	orgID := uuid.New()
	mainSQL := `SELECT * FROM "users" WHERE organization_id = $1 AND id = $2 ORDER BY "users"."id" LIMIT $3`
	mock.ExpectQuery(regexp.QuoteMeta(mainSQL)).
		WithArgs(orgID, int64(99), 1).
		WillReturnRows(sqlmock.NewRows(userColumns))

	_, err := repo.FindByID(context.Background(), orgID, 99)
	assert.ErrorIs(t, err, providers.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_FindByEmail(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewUserRepo(gdb)

	orgID := uuid.New()

	mainSQL := `SELECT * FROM "users" WHERE organization_id = $1 AND email = $2 ORDER BY "users"."id" LIMIT $3`
	mock.ExpectQuery(regexp.QuoteMeta(mainSQL)).
		WithArgs(orgID, "a@b.c", 1).
		WillReturnRows(sqlmock.NewRows(userColumns).
			AddRow(userRow(entities.User{
				ID: 7, OrganizationID: orgID, Email: "a@b.c", ProfileData: []byte(`{}`),
			})...))

	rolesSQL := `SELECT * FROM "user_roles" WHERE "user_roles"."user_id" = $1`
	mock.ExpectQuery(regexp.QuoteMeta(rolesSQL)).
		WithArgs(int64(7)).
		WillReturnRows(sqlmock.NewRows(userRoleColumns))

	got, err := repo.FindByEmail(context.Background(), orgID, "a@b.c")
	require.NoError(t, err)
	assert.Equal(t, int64(7), got.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_FindByOrgID(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewUserRepo(gdb)

	orgID := uuid.New()

	mainSQL := `SELECT * FROM "users" WHERE organization_id = $1`
	mock.ExpectQuery(regexp.QuoteMeta(mainSQL)).
		WithArgs(orgID).
		WillReturnRows(sqlmock.NewRows(userColumns).
			AddRow(userRow(entities.User{ID: 1, OrganizationID: orgID, Email: "a", ProfileData: []byte(`{}`)})...).
			AddRow(userRow(entities.User{ID: 2, OrganizationID: orgID, Email: "b", ProfileData: []byte(`{}`)})...))

	// Preload Roles: GORM batches the IN clause: IN (1,2).
	rolesSQL := `SELECT * FROM "user_roles" WHERE "user_roles"."user_id" IN ($1,$2)`
	mock.ExpectQuery(regexp.QuoteMeta(rolesSQL)).
		WithArgs(int64(1), int64(2)).
		WillReturnRows(sqlmock.NewRows(userRoleColumns))

	items, err := repo.FindByOrgID(context.Background(), orgID)
	require.NoError(t, err)
	assert.Len(t, items, 2)
	require.NoError(t, mock.ExpectationsWereMet())
}

// userInsertSQL mirrors the exact SQL GORM generates for a User INSERT.
// Column order is NOT the struct field order: `profile_data` is emitted last
// because its `default:'{}'` tag makes GORM wait on it and RETURN the DB value.
// The RETURNING clause also includes "profile_data" for the same reason.
const userInsertSQL = `INSERT INTO "users" ("organization_id","email","first_name","last_name","password_hash","avatar_url","onboarding_completed_at","created_at","updated_at","profile_data") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "profile_data","id"`

func TestUserRepo_Create(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewUserRepo(gdb)

	orgID := uuid.New()
	u := &entities.User{
		OrganizationID: orgID, Email: "new@x.com",
		FirstName: "N", LastName: "U", ProfileData: []byte(`{}`),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(userInsertSQL)).
		WithArgs(
			orgID, "new@x.com", "N", "U",
			nil, nil, nil,
			sqlmock.AnyArg(), sqlmock.AnyArg(), // created_at, updated_at
			sqlmock.AnyArg(), // profile_data (last)
		).
		WillReturnRows(sqlmock.NewRows([]string{"profile_data", "id"}).AddRow([]byte(`{}`), int64(5)))
	mock.ExpectCommit()

	id, err := repo.Create(context.Background(), u)
	require.NoError(t, err)
	assert.Equal(t, int64(5), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepo_Create_DuplicatedKey verifies the repo translates GORM's
// translated duplicate error into providers.ErrDuplicate.
func TestUserRepo_Create_DuplicatedKey(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewUserRepo(gdb)

	u := &entities.User{OrganizationID: uuid.New(), Email: "dup@x.com", ProfileData: []byte(`{}`)}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(userInsertSQL)).
		WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()

	_, err := repo.Create(context.Background(), u)
	assert.ErrorIs(t, err, providers.ErrDuplicate)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_Create_OtherError(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewUserRepo(gdb)

	u := &entities.User{OrganizationID: uuid.New(), Email: "x", ProfileData: []byte(`{}`)}
	boom := errors.New("network down")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(userInsertSQL)).WillReturnError(boom)
	mock.ExpectRollback()

	_, err := repo.Create(context.Background(), u)
	assert.ErrorIs(t, err, boom)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_AssignRole(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewUserRepo(gdb)

	sql := `INSERT INTO "user_roles" ("user_id","role") VALUES ($1,$2) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(7), entities.RoleTeacher).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))
	mock.ExpectCommit()

	err := repo.AssignRole(context.Background(), 7, entities.RoleTeacher)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_AssignRole_Duplicate(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewUserRepo(gdb)

	sql := `INSERT INTO "user_roles" ("user_id","role") VALUES ($1,$2) RETURNING "id"`
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(7), entities.RoleTeacher).
		WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()

	err := repo.AssignRole(context.Background(), 7, entities.RoleTeacher)
	assert.ErrorIs(t, err, providers.ErrDuplicate)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_RemoveRole_Success(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewUserRepo(gdb)

	sql := `DELETE FROM "user_roles" WHERE user_id = $1 AND role = $2`
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sql)).
		WithArgs(int64(7), entities.RoleAdmin).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.RemoveRole(context.Background(), 7, entities.RoleAdmin)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_RemoveRole_NotFound(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewUserRepo(gdb)

	sql := `DELETE FROM "user_roles" WHERE user_id = $1 AND role = $2`
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sql)).
		WithArgs(int64(7), entities.RoleAdmin).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.RemoveRole(context.Background(), 7, entities.RoleAdmin)
	assert.ErrorIs(t, err, providers.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepo_CompleteOnboarding(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewUserRepo(gdb)

	orgID := uuid.New()
	sql := `UPDATE "users" SET "onboarding_completed_at"=$1,"updated_at"=$2 WHERE id = $3 AND organization_id = $4`
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sql)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), int64(7), orgID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.CompleteOnboarding(context.Background(), orgID, 7)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepo_UpdateProfileData verifies the raw `?::jsonb` cast is applied
// with the marshalled JSON as the positional arg.
func TestUserRepo_UpdateProfileData(t *testing.T) {
	gdb, mock := dbtest.NewMockDB(t)
	repo := admin.NewUserRepo(gdb)

	orgID := uuid.New()
	data := map[string]any{"key": "value"}

	sql := `UPDATE "users" SET "profile_data"=$1::jsonb,"updated_at"=$2 WHERE id = $3 AND organization_id = $4`
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(sql)).
		WithArgs(`{"key":"value"}`, sqlmock.AnyArg(), int64(7), orgID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.UpdateProfileData(context.Background(), orgID, 7, data)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
