// Package main provides a small CLI helper to generate argon2id hashes for
// seeding users in db/seeds/seed.sql.
//
// Usage:
//
//	go run ./scripts/hash_password [plain-password]
//
// If no argument is provided, defaults to "admin123". The output is a single
// line containing the argon2id hash (PHC-encoded) suitable for pasting into
// SQL. Delegates to team-ai-toolkit/auth.HashPassword so the script stays in
// sync with the verifier used at login time.
package main

import (
	"fmt"
	"os"

	ttauth "github.com/educabot/team-ai-toolkit/auth"
)

func main() {
	plain := "admin123"
	if len(os.Args) > 1 {
		plain = os.Args[1]
	}

	hash, err := ttauth.HashPassword(plain)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(hash)
}
