package styles

import "testing"

func TestGetTheme_Dark(t *testing.T) {
	theme := GetTheme("dark")
	if theme == nil {
		t.Fatal("GetTheme(\"dark\") returned nil")
	}
	if theme.Name != "dark" {
		t.Errorf("expected Name \"dark\", got %q", theme.Name)
	}
}

func TestGetTheme_AllPresets(t *testing.T) {
	for _, name := range AllThemeNames() {
		theme := GetTheme(name)
		if theme == nil {
			t.Errorf("GetTheme(%q) returned nil", name)
			continue
		}
		if theme.Name != name {
			t.Errorf("GetTheme(%q).Name = %q, want %q", name, theme.Name, name)
		}
	}
}

func TestGetTheme_Unknown(t *testing.T) {
	theme := GetTheme("unknown")
	if theme == nil {
		t.Fatal("GetTheme(\"unknown\") returned nil")
	}
	if theme.Name != "dark" {
		t.Errorf("expected fallback to \"dark\", got %q", theme.Name)
	}
}

func TestAllThemeNames(t *testing.T) {
	names := AllThemeNames()
	if len(names) < 4 {
		t.Errorf("expected at least 4 theme names, got %d", len(names))
	}
}

func TestThemeColorsNotZero(t *testing.T) {
	for _, name := range AllThemeNames() {
		theme := GetTheme(name)
		if string(theme.Background) == "" {
			t.Errorf("theme %q: Background is empty", name)
		}
		if string(theme.Foreground) == "" {
			t.Errorf("theme %q: Foreground is empty", name)
		}
		if string(theme.Primary) == "" {
			t.Errorf("theme %q: Primary is empty", name)
		}
	}
}

func TestThemeBackgroundForegroundDiffer(t *testing.T) {
	for _, name := range AllThemeNames() {
		theme := GetTheme(name)
		if theme.Background == theme.Foreground {
			t.Errorf("theme %q: Background and Foreground are identical (%s)", name, string(theme.Background))
		}
	}
}
