//go:build integration

package repository

import (
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const testDSN = "host=localhost user=postgres password=postgres dbname=fiber_api_test port=5432 sslmode=disable"
const testMigrateURL = "postgres://postgres:postgres@localhost:5432/fiber_api_test?sslmode=disable"

// setupTestDB connects to the dedicated test database, applies the existing
// golang-migrate migrations, and truncates tables so each test starts clean.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(postgres.Open(testDSN), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect test DB:", err)
	}

	m, err := migrate.New("file://../../migrations", testMigrateURL)
	if err != nil {
		t.Fatal("Failed to init migrate:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatal("Failed to run migrations:", err)
	}

	truncateAllTables(t, db)

	return db
}

// truncateAllTables empties every table in the public schema so each test
// starts from a clean slate. Tables are discovered dynamically, so new models
// are cleaned automatically without updating this helper. The migrate
// bookkeeping table (schema_migrations) is excluded — truncating it would make
// golang-migrate re-run every migration.
func truncateAllTables(t *testing.T, db *gorm.DB) {
	var tables []string
	if err := db.Raw(`SELECT tablename FROM pg_tables
		WHERE schemaname = 'public' AND tablename <> 'schema_migrations'`).
		Scan(&tables).Error; err != nil {
		t.Fatal("Failed to list tables:", err)
	}

	if len(tables) == 0 {
		return
	}

	if err := db.Exec("TRUNCATE " + strings.Join(tables, ", ") + " RESTART IDENTITY CASCADE").Error; err != nil {
		t.Fatal("Failed to truncate tables:", err)
	}
}
