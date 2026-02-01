package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yangsijun/bible-tui/internal/config"
	"github.com/yangsijun/bible-tui/internal/tui/styles"
)

func newTestSettingsModel() SettingsModel {
	return NewSettings(nil, styles.DefaultDarkTheme(), 80, 24)
}

func TestSettingsModel_Init(t *testing.T) {
	m := newTestSettingsModel()

	if m.focusRow != 0 {
		t.Errorf("expected focusRow=0, got %d", m.focusRow)
	}
	if m.loaded {
		t.Error("expected loaded=false")
	}
	if m.themeIdx != 0 {
		t.Errorf("expected themeIdx=0, got %d", m.themeIdx)
	}
	if m.fontSizeIdx != 0 {
		t.Errorf("expected fontSizeIdx=0, got %d", m.fontSizeIdx)
	}
	if m.versionIdx != 0 {
		t.Errorf("expected versionIdx=0, got %d", m.versionIdx)
	}
}

func TestSettingsModel_Loaded(t *testing.T) {
	m := newTestSettingsModel()

	cfg := &config.Config{ThemeName: "dark", FontSize: 2, VersionCode: "GAE"}
	m, _ = m.Update(SettingsLoadedMsg{Config: cfg})

	if !m.loaded {
		t.Error("expected loaded=true")
	}
	if m.themeIdx != 0 {
		t.Errorf("expected themeIdx=0 for 'dark', got %d", m.themeIdx)
	}
	if m.fontSizeIdx != 1 {
		t.Errorf("expected fontSizeIdx=1 for FontSize=2, got %d", m.fontSizeIdx)
	}
	if m.versionIdx != 0 {
		t.Errorf("expected versionIdx=0 for 'GAE', got %d", m.versionIdx)
	}
}

func TestSettingsModel_Loaded_Solarized(t *testing.T) {
	m := newTestSettingsModel()

	cfg := &config.Config{ThemeName: "solarized", FontSize: 3, VersionCode: "GAE"}
	m, _ = m.Update(SettingsLoadedMsg{Config: cfg})

	if m.themeIdx != 2 {
		t.Errorf("expected themeIdx=2 for 'solarized', got %d", m.themeIdx)
	}
	if m.fontSizeIdx != 2 {
		t.Errorf("expected fontSizeIdx=2 for FontSize=3, got %d", m.fontSizeIdx)
	}
}

func TestSettingsModel_Loaded_Error(t *testing.T) {
	m := newTestSettingsModel()

	m, _ = m.Update(SettingsLoadedMsg{Err: fmt.Errorf("db error")})

	if !m.loaded {
		t.Error("expected loaded=true even on error")
	}
	if m.themeIdx != 0 {
		t.Errorf("expected themeIdx unchanged, got %d", m.themeIdx)
	}
}

func TestSettingsModel_Navigation(t *testing.T) {
	m := newTestSettingsModel()

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if m.focusRow != 1 {
		t.Errorf("expected focusRow=1 after down, got %d", m.focusRow)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if m.focusRow != 2 {
		t.Errorf("expected focusRow=2 after second down, got %d", m.focusRow)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if m.focusRow != 2 {
		t.Errorf("expected focusRow=2 (clamped), got %d", m.focusRow)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	if m.focusRow != 1 {
		t.Errorf("expected focusRow=1 after up, got %d", m.focusRow)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	if m.focusRow != 0 {
		t.Errorf("expected focusRow=0 (clamped), got %d", m.focusRow)
	}
}

func TestSettingsModel_Navigation_ArrowKeys(t *testing.T) {
	m := newTestSettingsModel()

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.focusRow != 1 {
		t.Errorf("expected focusRow=1 after KeyDown, got %d", m.focusRow)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.focusRow != 0 {
		t.Errorf("expected focusRow=0 after KeyUp, got %d", m.focusRow)
	}
}

func TestSettingsModel_ChangeTheme(t *testing.T) {
	m := newTestSettingsModel()
	m.focusRow = 0

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	if m.themeIdx != 1 {
		t.Errorf("expected themeIdx=1 after right, got %d", m.themeIdx)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	if m.themeIdx != 2 {
		t.Errorf("expected themeIdx=2 after second right, got %d", m.themeIdx)
	}
}

func TestSettingsModel_ChangeTheme_Wrap(t *testing.T) {
	m := newTestSettingsModel()
	m.focusRow = 0
	themeCount := len(styles.AllThemeNames())

	for i := 0; i < themeCount; i++ {
		m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	}
	if m.themeIdx != 0 {
		t.Errorf("expected themeIdx=0 after wrapping, got %d", m.themeIdx)
	}
}

func TestSettingsModel_ChangeTheme_WrapLeft(t *testing.T) {
	m := newTestSettingsModel()
	m.focusRow = 0

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	expected := len(styles.AllThemeNames()) - 1
	if m.themeIdx != expected {
		t.Errorf("expected themeIdx=%d after left wrap, got %d", expected, m.themeIdx)
	}
}

func TestSettingsModel_ChangeFontSize(t *testing.T) {
	m := newTestSettingsModel()
	m.focusRow = 1

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	if m.fontSizeIdx != 1 {
		t.Errorf("expected fontSizeIdx=1 after right, got %d", m.fontSizeIdx)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	if m.fontSizeIdx != 2 {
		t.Errorf("expected fontSizeIdx=2 after second right, got %d", m.fontSizeIdx)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	if m.fontSizeIdx != 0 {
		t.Errorf("expected fontSizeIdx=0 after wrap, got %d", m.fontSizeIdx)
	}
}

func TestSettingsModel_View(t *testing.T) {
	m := newTestSettingsModel()
	cfg := &config.Config{ThemeName: "dark", FontSize: 2, VersionCode: "GAE"}
	m, _ = m.Update(SettingsLoadedMsg{Config: cfg})

	view := m.View()

	if !strings.Contains(view, "설정") {
		t.Error("view should contain '설정'")
	}
	if !strings.Contains(view, "테마") {
		t.Error("view should contain '테마'")
	}
	if !strings.Contains(view, "글자크기") {
		t.Error("view should contain '글자크기'")
	}
	if !strings.Contains(view, "기본역본") {
		t.Error("view should contain '기본역본'")
	}
	if !strings.Contains(view, "s:저장") {
		t.Error("view should contain 's:저장'")
	}
	if !strings.Contains(view, "Esc:취소") {
		t.Error("view should contain 'Esc:취소'")
	}
	if !strings.Contains(view, "▸") {
		t.Error("view should contain cursor '▸'")
	}
}

func TestSettingsModel_Save(t *testing.T) {
	m := newTestSettingsModel()
	cfg := &config.Config{ThemeName: "dark", FontSize: 2, VersionCode: "GAE"}
	m, _ = m.Update(SettingsLoadedMsg{Config: cfg})

	var cmd tea.Cmd
	m, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})

	if cmd == nil {
		t.Error("expected non-nil cmd after save")
	}
	if !m.saved {
		t.Error("expected saved=true after pressing 's'")
	}
}

func TestSettingsModel_SaveEnter(t *testing.T) {
	m := newTestSettingsModel()
	cfg := &config.Config{ThemeName: "dark", FontSize: 2, VersionCode: "GAE"}
	m, _ = m.Update(SettingsLoadedMsg{Config: cfg})

	var cmd tea.Cmd
	m, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Error("expected non-nil cmd after enter")
	}
	if !m.saved {
		t.Error("expected saved=true after pressing enter")
	}
}
