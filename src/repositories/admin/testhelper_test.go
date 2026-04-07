package admin_test

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/educabot/alizia-be/src/core/entities"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://postgres:postgres@localhost:5480/alizia?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("skipping repo test: cannot connect to database: %v", err)
	}

	// Ensure enum type exists
	db.Exec("DO $$ BEGIN CREATE TYPE member_role AS ENUM ('teacher', 'coordinator', 'admin'); EXCEPTION WHEN duplicate_object THEN null; END $$;")

	// Auto-migrate tables for test
	err = db.AutoMigrate(&entities.Organization{}, &entities.User{}, &entities.UserRole{})
	if err != nil {
		t.Fatalf("failed to auto-migrate: %v", err)
	}

	return db
}

func cleanupTestData(t *testing.T, db *gorm.DB, orgSlug string) {
	t.Helper()

	var org entities.Organization
	if err := db.Where("slug = ?", orgSlug).First(&org).Error; err != nil {
		return // nothing to clean
	}

	db.Where("user_id IN (SELECT id FROM users WHERE organization_id = ?)", org.ID).Delete(&entities.UserRole{})
	db.Where("organization_id = ?", org.ID).Delete(&entities.User{})
	db.Delete(&org)
}

func uniqueSlug(t *testing.T) string {
	return fmt.Sprintf("test-%s", t.Name())
}
