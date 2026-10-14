//go:build integration

package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func countDBConnections(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var n int64
	if err := db.Raw(`SELECT count(*) FROM pg_stat_activity WHERE datname = current_database() AND backend_type = 'client backend'`).
		Scan(&n).Error; err != nil {
		t.Fatal("Failed to count connections:", err)
	}
	return n
}

func TestSetupTestDB_ReleasesConnections(t *testing.T) {
	probe := setupTestDB(t)
	baseline := countDBConnections(t, probe)

	for i := 0; i < 5; i++ {
		t.Run(fmt.Sprintf("setup_%d", i), func(t *testing.T) {
			setupTestDB(t)
		})
	}

	// Postgres drops a backend shortly after the client closes it, so poll.
	after := countDBConnections(t, probe)
	deadline := time.Now().Add(5 * time.Second)
	for after > baseline && time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
		after = countDBConnections(t, probe)
	}

	assert.LessOrEqual(t, after, baseline, "setupTestDB leaked connections")
}
