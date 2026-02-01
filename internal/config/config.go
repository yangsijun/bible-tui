package config

import (
	"fmt"
	"strconv"

	"github.com/yangsijun/bible-tui/internal/db"
)

type Config struct {
	ThemeName   string
	FontSize    int
	VersionCode string
}

func LoadConfig(database *db.DB) (*Config, error) {
	cfg := &Config{
		ThemeName:   "dark",
		FontSize:    2,
		VersionCode: "GAE",
	}

	themeName, err := database.GetSetting("theme_name")
	if err != nil {
		return nil, fmt.Errorf("get theme_name setting: %w", err)
	}
	if themeName != "" {
		cfg.ThemeName = themeName
	}

	fontSize, err := database.GetSetting("font_size")
	if err != nil {
		return nil, fmt.Errorf("get font_size setting: %w", err)
	}
	if fontSize != "" {
		size, err := strconv.Atoi(fontSize)
		if err != nil {
			return nil, fmt.Errorf("parse font_size: %w", err)
		}
		cfg.FontSize = size
	}

	versionCode, err := database.GetSetting("default_version")
	if err != nil {
		return nil, fmt.Errorf("get default_version setting: %w", err)
	}
	if versionCode != "" {
		cfg.VersionCode = versionCode
	}

	return cfg, nil
}

func SaveConfig(database *db.DB, cfg *Config) error {
	if err := database.SetSetting("theme_name", cfg.ThemeName); err != nil {
		return fmt.Errorf("set theme_name: %w", err)
	}

	if err := database.SetSetting("font_size", strconv.Itoa(cfg.FontSize)); err != nil {
		return fmt.Errorf("set font_size: %w", err)
	}

	if err := database.SetSetting("default_version", cfg.VersionCode); err != nil {
		return fmt.Errorf("set default_version: %w", err)
	}

	return nil
}
