package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sijun-dong/bible-tui/internal/db"
	"github.com/sijun-dong/bible-tui/internal/tui/styles"
)

type SearchResultsMsg struct {
	Results []db.SearchResult
	Query   string
	Err     error
}

type GoToVerseMsg struct {
	BookCode string
	Chapter  int
	Verse    int
}

type SearchModel struct {
	input     textinput.Model
	results   []db.SearchResult
	selected  int
	database  *db.DB
	theme     *styles.Theme
	query     string
	loading   bool
	noResults bool
	width     int
	height    int
}

func NewSearch(database *db.DB, theme *styles.Theme, width, height int) SearchModel {
	ti := textinput.New()
	ti.Placeholder = "검색어를 입력하세요..."
	ti.Focus()
	ti.CharLimit = 100
	if width > 4 {
		ti.Width = width - 4
	}
	return SearchModel{
		input:    ti,
		database: database,
		theme:    theme,
		width:    width,
		height:   height,
	}
}

func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.input.Focused() {
			switch msg.String() {
			case "enter":
				query := strings.TrimSpace(m.input.Value())
				if query != "" {
					m.query = query
					m.loading = true
					return m, searchVerses(m.database, "GAE", query, 20)
				}
				return m, nil
			case "down", "tab":
				if len(m.results) > 0 {
					m.input.Blur()
				}
				return m, nil
			}
		} else {
			switch msg.String() {
			case "up", "k":
				if m.selected > 0 {
					m.selected--
				} else {
					m.input.Focus()
				}
				return m, nil
			case "down", "j":
				if m.selected < len(m.results)-1 {
					m.selected++
				}
				return m, nil
			case "enter":
				if m.selected < len(m.results) {
					r := m.results[m.selected]
					return m, func() tea.Msg {
						return GoToVerseMsg{
							BookCode: r.Verse.BookCode,
							Chapter:  r.Verse.Chapter,
							Verse:    r.Verse.VerseNum,
						}
					}
				}
				return m, nil
			case "/":
				m.input.Focus()
				return m, nil
			}
		}

	case SearchResultsMsg:
		m.loading = false
		if msg.Err != nil {
			m.results = nil
			m.noResults = true
			return m, nil
		}
		m.results = msg.Results
		m.noResults = len(msg.Results) == 0
		m.selected = 0
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m SearchModel) View() string {
	var b strings.Builder

	inputStyle := lipgloss.NewStyle().Padding(0, 1)
	b.WriteString(inputStyle.Render(m.input.View()))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString("  검색 중...")
		return b.String()
	}
	if m.noResults && m.query != "" {
		b.WriteString(fmt.Sprintf("  검색 결과가 없습니다: %q", m.query))
		return b.String()
	}

	for i, r := range m.results {
		cursor := "  "
		if i == m.selected && !m.input.Focused() {
			cursor = "▸ "
		}

		ref := fmt.Sprintf("%s %d:%d", r.Verse.BookName, r.Verse.Chapter, r.Verse.VerseNum)
		refStyle := lipgloss.NewStyle().Foreground(m.theme.Primary).Bold(true)

		text := r.Verse.Text
		maxRunes := m.width - 10
		runes := []rune(text)
		if maxRunes > 0 && len(runes) > maxRunes {
			text = string(runes[:maxRunes]) + "..."
		}

		b.WriteString(fmt.Sprintf("%s%s\n    %s\n", cursor, refStyle.Render(ref), text))
	}

	return b.String()
}

func searchVerses(database *db.DB, versionCode, query string, limit int) tea.Cmd {
	return func() tea.Msg {
		if database == nil {
			return SearchResultsMsg{Err: fmt.Errorf("no database")}
		}
		results, err := database.SearchVerses(versionCode, query, limit)
		return SearchResultsMsg{Results: results, Query: query, Err: err}
	}
}
