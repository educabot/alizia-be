// Package queries embeds raw SQL files used by the admin repositories.
//
// Convention: any query longer than ~3 lines or that uses CTEs, window
// functions or subqueries lives in a .sql file here and is loaded with
// embed. Inline raw SQL in Go source is reserved for trivial fragments.
package queries

import _ "embed"

//go:embed shared_class_numbers.sql
var SharedClassNumbers string
