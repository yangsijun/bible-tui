package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sijun-dong/bible-tui/internal/config"
	"github.com/sijun-dong/bible-tui/internal/db"
	"github.com/sijun-dong/bible-tui/internal/tui/styles"
)

type SettingsLoadedMsg struct {
	Config *config.Config
	Err    error
}

type SettingsSavedMsg struct {
	Err error
}

type ThemeChangedMsg struct {
	Theme *styles.Theme
}

var (
	fontSizeLabels   = []string{"작게", "보통", "크게"}
	versionLabels    = []string{"개역개정 (GAE)"}
	versionCodes     = []string{"GAE"}
	settingsRowCount = 3
)

type SettingsModel struct {
	database    *db.DB
	theme       *styles.Theme
	width       int
	height      int
	focusRow    int
	themeIdx    int
	fontSizeIdx int
	versionIdx  int
	loaded      bool
	saved       bool
}

func NewSettings(database *db.DB, theme *styles.Theme, width, height int) SettingsModel {
	return SettingsModel{
		database: database,
		theme:    theme,
		width:    width,
		height:   height,
	}
}

func LoadSettings(database *db.DB) tea.Cmd {
	return func() tea.Msg {
		if database == nil {
			return SettingsLoadedMsg{Err: fmt.Errorf("no database")}
		}
		cfg, err := config.LoadConfig(database)
		return SettingsLoadedMsg{Config: cfg, Err: err}
	}
}

func (m SettingsModel) Update(msg tea.Msg) (SettingsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case SettingsLoadedMsg:
		m.loaded = true
		if msg.Err != nil || msg.Config == nil {
			return m, nil
		}
		m.themeIdx = themeNameToIdx(msg.Config.ThemeName)
		// FontSize is 1-based in config, 0-based as index
		m.fontSizeIdx = msg.Config.FontSize - 1
		if m.fontSizeIdx < 0 {
			m.fontSizeIdx = 0
		}
		if m.fontSizeIdx >= len(fontSizeLabels) {
			m.fontSizeIdx = len(fontSizeLabels) - 1
		}
		m.versionIdx = versionCodeToIdx(msg.Config.VersionCode)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.focusRow > 0 {
				m.focusRow--
			}
			return m, nil
		case "down", "j":
			if m.focusRow < settingsRowCount-1 {
				m.focusRow++
			}
			return m, nil
		case "left", "h":
			m.decrementOption()
			return m, nil
		case "right", "l":
			m.incrementOption()
			return m, nil
		case "enter", "s":
			m.saved = true
			return m, tea.Batch(m.saveConfig(), m.emitThemeChange())
		}
	}
	return m, nil
}

func (m SettingsModel) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Primary)
	b.WriteString("\n  " + titleStyle.Render("설정") + "\n\n")

	themeNames := styles.AllThemeNames()

	rows := []struct {
		label string
		value string
	}{
		{"테마", themeNames[m.themeIdx]},
		{"글자크기", fontSizeLabels[m.fontSizeIdx]},
		{"기본역본", versionLabels[m.versionIdx]},
	}

	for i, row := range rows {
		cursor := "  "
		labelStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
		valueStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
		if i == m.focusRow {
			cursor = "▸ "
			labelStyle = lipgloss.NewStyle().Foreground(m.theme.Primary).Bold(true)
			valueStyle = lipgloss.NewStyle().Foreground(m.theme.Primary)
		}

		label := labelStyle.Render(fmt.Sprintf("%-8s", row.label))
		value := valueStyle.Render(row.value)
		arrowStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
		left := arrowStyle.Render("◀")
		right := arrowStyle.Render("▶")

		b.WriteString(fmt.Sprintf("  %s%s %s %s %s\n", cursor, label, left, value, right))
	}

	footerStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
	b.WriteString("\n  " + footerStyle.Render("s:저장  Esc:취소") + "\n")

	return b.String()
}

func (m *SettingsModel) incrementOption() {
	switch m.focusRow {
	case 0:
		themeNames := styles.AllThemeNames()
		m.themeIdx = (m.themeIdx + 1) % len(themeNames)
	case 1:
		m.fontSizeIdx = (m.fontSizeIdx + 1) % len(fontSizeLabels)
	case 2:
		m.versionIdx = (m.versionIdx + 1) % len(versionLabels)
	}
}

func (m *SettingsModel) decrementOption() {
	switch m.focusRow {
	case 0:
		themeNames := styles.AllThemeNames()
		m.themeIdx = (m.themeIdx - 1 + len(themeNames)) % len(themeNames)
	case 1:
		m.fontSizeIdx = (m.fontSizeIdx - 1 + len(fontSizeLabels)) % len(fontSizeLabels)
	case 2:
		m.versionIdx = (m.versionIdx - 1 + len(versionLabels)) % len(versionLabels)
	}
}

func (m SettingsModel) saveConfig() tea.Cmd {
	return func() tea.Msg {
		themeNames := styles.AllThemeNames()
		cfg := &config.Config{
			ThemeName:   themeNames[m.themeIdx],
			FontSize:    m.fontSizeIdx + 1,
			VersionCode: versionCodes[m.versionIdx],
		}
		if m.database == nil {
			return SettingsSavedMsg{Err: fmt.Errorf("no database")}
		}
		err := config.SaveConfig(m.database, cfg)
		return SettingsSavedMsg{Err: err}
	}
}

func (m SettingsModel) emitThemeChange() tea.Cmd {
	return func() tea.Msg {
		themeNames := styles.AllThemeNames()
		return ThemeChangedMsg{Theme: styles.GetTheme(themeNames[m.themeIdx])}
	}
}

func themeNameToIdx(name string) int {
	for i, n := range styles.AllThemeNames() {
		if n == name {
			return i
		}
	}
	return 0
}

func versionCodeToIdx(code string) int {
	for i, c := range versionCodes {
		if c == code {
			return i
		}
	}
	return 0
}
