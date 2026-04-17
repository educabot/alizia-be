// Package dbtest provides helpers for writing repo-layer tests without a real
// Postgres. It wires `go-sqlmock` into a GORM `*gorm.DB` configured with the
// Postgres dialect so tests exercise the same SQL generation path as production
// code. Queries must match via `regexp.QuoteMeta(sql)` so we validate the exact
// SQL GORM produces instead of a loose regex — GORM upgrades or query changes
// surface as test failures rather than silent prod breakage.
//
// This helper is a candidate for extraction into team-ai-toolkit once the API
// stabilises across a few repos.
package dbtest

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewMockDB returns a *gorm.DB wired to sqlmock with the Postgres dialect.
// The returned mock uses the default (regex) query matcher, so callers must
// wrap expected SQL with `regexp.QuoteMeta(...)` to assert exact-string match.
func NewMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	gdb, err := gorm.Open(
		postgres.New(postgres.Config{
			Conn:       sqlDB,
			DriverName: "postgres",
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	return gdb, mock
}
