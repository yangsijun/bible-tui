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
	StatusBarBg    lipgloss.Color
	StatusBarFg    lipgloss.Color
	VerseNumber    lipgloss.Color
	SectionTitle   lipgloss.Color
	FootnoteMarker lipgloss.Color
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
		StatusBarBg:    lipgloss.Color("#1f2335"),
		StatusBarFg:    lipgloss.Color("#737aa2"),
		VerseNumber:    lipgloss.Color("#e0af68"),
		SectionTitle:   lipgloss.Color("#9ece6a"),
		FootnoteMarker: lipgloss.Color("#7dcfff"),
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
		StatusBarBg:    lipgloss.Color("#d5d6db"),
		StatusBarFg:    lipgloss.Color("#8990b3"),
		VerseNumber:    lipgloss.Color("#965027"),
		SectionTitle:   lipgloss.Color("#33635c"),
		FootnoteMarker: lipgloss.Color("#166775"),
	}
}

// SolarizedTheme returns a theme based on the Solarized Dark palette.
func SolarizedTheme() *Theme {
	return &Theme{
		Name:           "solarized",
		Background:     lipgloss.Color("#002b36"),
		Foreground:     lipgloss.Color("#839496"),
		Primary:        lipgloss.Color("#268bd2"),
		Secondary:      lipgloss.Color("#2aa198"),
		Muted:          lipgloss.Color("#586e75"),
		Error:          lipgloss.Color("#dc322f"),
		HighlightBg:    lipgloss.Color("#073642"),
		StatusBarBg:    lipgloss.Color("#073642"),
		StatusBarFg:    lipgloss.Color("#657b83"),
		VerseNumber:    lipgloss.Color("#b58900"),
		SectionTitle:   lipgloss.Color("#cb4b16"),
		FootnoteMarker: lipgloss.Color("#6c71c4"),
	}
}

// NordTheme returns a theme based on the Nord palette.
func NordTheme() *Theme {
	return &Theme{
		Name:           "nord",
		Background:     lipgloss.Color("#2E3440"),
		Foreground:     lipgloss.Color("#D8DEE9"),
		Primary:        lipgloss.Color("#88C0D0"),
		Secondary:      lipgloss.Color("#B48EAD"),
		Muted:          lipgloss.Color("#4C566A"),
		Error:          lipgloss.Color("#BF616A"),
		HighlightBg:    lipgloss.Color("#3B4252"),
		StatusBarBg:    lipgloss.Color("#3B4252"),
		StatusBarFg:    lipgloss.Color("#616E88"),
		VerseNumber:    lipgloss.Color("#EBCB8B"),
		SectionTitle:   lipgloss.Color("#A3BE8C"),
		FootnoteMarker: lipgloss.Color("#81A1C1"),
	}
}

// GetTheme returns a theme by name. Falls back to dark theme for unknown names.
func GetTheme(name string) *Theme {
	switch name {
	case "dark":
		return DefaultDarkTheme()
	case "light":
		return DefaultLightTheme()
	case "solarized":
		return SolarizedTheme()
	case "nord":
		return NordTheme()
	default:
		return DefaultDarkTheme()
	}
}

// AllThemeNames returns the names of all available preset themes.
func AllThemeNames() []string {
	return []string{"dark", "light", "solarized", "nord"}
}
