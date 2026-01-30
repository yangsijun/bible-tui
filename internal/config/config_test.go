package config

import (
	"os"
	"testing"

	"github.com/sijun-dong/bible-tui/internal/db"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Setup: Create in-memory database and migrate
	database, err := db.OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory failed: %v", err)
	}
	defer database.Close()

	if err := database.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	// Act: Load config without any settings (should use defaults)
	cfg, err := LoadConfig(database)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Assert: Check defaults
	if cfg.ThemeName != "dark" {
		t.Errorf("ThemeName: got %q, want %q", cfg.ThemeName, "dark")
	}
	if cfg.FontSize != 2 {
		t.Errorf("FontSize: got %d, want %d", cfg.FontSize, 2)
	}
	if cfg.VersionCode != "GAE" {
		t.Errorf("VersionCode: got %q, want %q", cfg.VersionCode, "GAE")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Setup: Create in-memory database and migrate
	database, err := db.OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory failed: %v", err)
	}
	defer database.Close()

	if err := database.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	// Act: Create and save custom config
	originalCfg := &Config{
		ThemeName:   "solarized",
		FontSize:    3,
		VersionCode: "KJV",
	}

	if err := SaveConfig(database, originalCfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Act: Load the config back
	loadedCfg, err := LoadConfig(database)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Assert: Check that loaded config matches saved config
	if loadedCfg.ThemeName != originalCfg.ThemeName {
		t.Errorf("ThemeName: got %q, want %q", loadedCfg.ThemeName, originalCfg.ThemeName)
	}
	if loadedCfg.FontSize != originalCfg.FontSize {
		t.Errorf("FontSize: got %d, want %d", loadedCfg.FontSize, originalCfg.FontSize)
	}
	if loadedCfg.VersionCode != originalCfg.VersionCode {
		t.Errorf("VersionCode: got %q, want %q", loadedCfg.VersionCode, originalCfg.VersionCode)
	}
}

func TestConfigPersistence(t *testing.T) {
	// Setup: Create temporary file for database
	tmpFile, err := os.CreateTemp("", "test-config-*.db")
	if err != nil {
		t.Fatalf("CreateTemp failed: %v", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// First connection: Create and save config
	database1, err := db.Open(tmpPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	if err := database1.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	originalCfg := &Config{
		ThemeName:   "nord",
		FontSize:    1,
		VersionCode: "NIV",
	}

	if err := SaveConfig(database1, originalCfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	database1.Close()

	// Second connection: Reopen database and load config
	database2, err := db.Open(tmpPath)
	if err != nil {
		t.Fatalf("Open (second) failed: %v", err)
	}
	defer database2.Close()

	loadedCfg, err := LoadConfig(database2)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Assert: Check that config persisted across connections
	if loadedCfg.ThemeName != originalCfg.ThemeName {
		t.Errorf("ThemeName: got %q, want %q", loadedCfg.ThemeName, originalCfg.ThemeName)
	}
	if loadedCfg.FontSize != originalCfg.FontSize {
		t.Errorf("FontSize: got %d, want %d", loadedCfg.FontSize, originalCfg.FontSize)
	}
	if loadedCfg.VersionCode != originalCfg.VersionCode {
		t.Errorf("VersionCode: got %q, want %q", loadedCfg.VersionCode, originalCfg.VersionCode)
	}
}
