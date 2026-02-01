package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yangsijun/bible-tui/internal/db"
)

var testDB *db.DB // only set in tests

func getDB() (*db.DB, error) {
	if testDB != nil {
		return testDB, nil
	}
	return openDB()
}

func openDB() (*db.DB, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("get config dir: %w", err)
	}
	dbPath := filepath.Join(configDir, "bible-tui", "bible.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}
	return db.Open(dbPath)
}
