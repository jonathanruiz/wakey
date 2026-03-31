package tests

import (
	"os"
	"testing"
	"wakey/internal/config"
)

// setupTestDB creates a temporary SQLite database for testing and returns a
// cleanup function that removes the file and resets the DB connection.
func setupTestDB(t *testing.T) func() {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "wakey_test_*.db")
	if err != nil {
		t.Fatalf("failed to create temp db: %v", err)
	}
	tmpFile.Close()

	config.DBPath = tmpFile.Name()
	config.ResetDB()
	config.CreateConfig()

	return func() {
		config.ResetDB()
		os.Remove(tmpFile.Name())
	}
}
