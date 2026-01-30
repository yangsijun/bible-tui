package tui

import (
	"strings"
	"testing"

	"github.com/sijun-dong/bible-tui/internal/tui/styles"
)

func TestHelpModel_Init(t *testing.T) {
	m := NewHelp(styles.DefaultDarkTheme(), 80, 24)
	v := m.View()
	if v == "" {
		t.Error("expected non-empty view")
	}
}

func TestHelpModel_ContainsKeybindings(t *testing.T) {
	content := renderHelpContent(styles.DefaultDarkTheme())
	checks := []string{"Ctrl+C", "Esc", "Enter", "책 목록", "읽기 화면", "검색"}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected help to contain %q", check)
		}
	}
}

func TestHelpModel_ContainsSections(t *testing.T) {
	m := NewHelp(styles.DefaultDarkTheme(), 80, 24)
	content := m.View()
	sections := []string{"전역", "장 선택", "책갈피"}
	for _, s := range sections {
		if !strings.Contains(content, s) {
			t.Errorf("expected help to contain section %q", s)
		}
	}
}
