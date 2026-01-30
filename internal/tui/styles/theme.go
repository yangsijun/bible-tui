package styles

import "github.com/charmbracelet/lipgloss"

// Theme defines the color scheme for the TUI
type Theme struct {
	Name        string
	Background  lipgloss.Color
	Foreground  lipgloss.Color
	Primary     lipgloss.Color
	Secondary   lipgloss.Color
	Muted       lipgloss.Color
	Error       lipgloss.Color
	HighlightBg lipgloss.Color
	StatusBarBg lipgloss.Color
	StatusBarFg lipgloss.Color
}

func DefaultDarkTheme() *Theme {
	return &Theme{
		Name:        "dark",
		Background:  lipgloss.Color("#1a1b26"),
		Foreground:  lipgloss.Color("#c0caf5"),
		Primary:     lipgloss.Color("#7aa2f7"),
		Secondary:   lipgloss.Color("#bb9af7"),
		Muted:       lipgloss.Color("#565f89"),
		Error:       lipgloss.Color("#f7768e"),
		HighlightBg: lipgloss.Color("#292e42"),
		StatusBarBg: lipgloss.Color("#1f2335"),
		StatusBarFg: lipgloss.Color("#737aa2"),
	}
}

func DefaultLightTheme() *Theme {
	return &Theme{
		Name:        "light",
		Background:  lipgloss.Color("#ffffff"),
		Foreground:  lipgloss.Color("#343b58"),
		Primary:     lipgloss.Color("#34548a"),
		Secondary:   lipgloss.Color("#5a4a78"),
		Muted:       lipgloss.Color("#9699a3"),
		Error:       lipgloss.Color("#8c4351"),
		HighlightBg: lipgloss.Color("#e9e9ed"),
		StatusBarBg: lipgloss.Color("#d5d6db"),
		StatusBarFg: lipgloss.Color("#8990b3"),
	}
}
