// Package main provides a small CLI helper to generate bcrypt hashes for
// seeding users in db/seeds/seed.sql.
//
// Usage:
//
//	go run ./scripts/hash_password [plain-password]
//
// If no argument is provided, defaults to "admin123". The output is a single
// line containing the bcrypt hash (cost 12) suitable for pasting into SQL.
package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	plain := "admin123"
	if len(os.Args) > 1 {
		plain = os.Args[1]
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plain), 12)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(hash))
}
