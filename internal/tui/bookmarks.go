package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sijun-dong/bible-tui/internal/db"
	"github.com/sijun-dong/bible-tui/internal/tui/styles"
)

type BookmarkTab int

const (
	TabBookmarks BookmarkTab = iota
	TabHighlights
)

type BookmarksLoadedMsg struct {
	Bookmarks  []db.BookmarkWithVerse
	Highlights []db.HighlightWithVerse
	Err        error
}

type BookmarkDeletedMsg struct {
	ID  int64
	Err error
}

type BookmarkModel struct {
	tab        BookmarkTab
	bookmarks  []db.BookmarkWithVerse
	highlights []db.HighlightWithVerse
	selected   int
	database   *db.DB
	theme      *styles.Theme
	loaded     bool
	width      int
	height     int
}

func NewBookmarks(database *db.DB, theme *styles.Theme, width, height int) BookmarkModel {
	return BookmarkModel{
		tab:      TabBookmarks,
		database: database,
		theme:    theme,
		width:    width,
		height:   height,
	}
}

func LoadBookmarks(database *db.DB) tea.Cmd {
	return func() tea.Msg {
		if database == nil {
			return BookmarksLoadedMsg{Err: fmt.Errorf("no database")}
		}
		bookmarks, err := database.ListBookmarks(50, 0)
		if err != nil {
			return BookmarksLoadedMsg{Err: err}
		}
		highlights, err := database.ListHighlights(50, 0)
		if err != nil {
			return BookmarksLoadedMsg{Err: err}
		}
		return BookmarksLoadedMsg{Bookmarks: bookmarks, Highlights: highlights}
	}
}

func (m BookmarkModel) Update(msg tea.Msg) (BookmarkModel, tea.Cmd) {
	switch msg := msg.(type) {
	case BookmarksLoadedMsg:
		m.loaded = true
		if msg.Err == nil {
			m.bookmarks = msg.Bookmarks
			m.highlights = msg.Highlights
		}
		m.selected = 0
		return m, nil
	case BookmarkDeletedMsg:
		if msg.Err == nil {
			return m, LoadBookmarks(m.database)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.tab == TabBookmarks {
				m.tab = TabHighlights
			} else {
				m.tab = TabBookmarks
			}
			m.selected = 0
			return m, nil
		case "down", "j":
			max := m.currentListLen() - 1
			if m.selected < max {
				m.selected++
			}
			return m, nil
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
			return m, nil
		case "d":
			return m, m.deleteSelected()
		case "enter":
			return m, m.goToSelected()
		}
	}
	return m, nil
}

func (m BookmarkModel) currentListLen() int {
	if m.tab == TabBookmarks {
		return len(m.bookmarks)
	}
	return len(m.highlights)
}

func (m BookmarkModel) deleteSelected() tea.Cmd {
	if m.database == nil {
		return nil
	}
	if m.tab == TabBookmarks && m.selected < len(m.bookmarks) {
		id := m.bookmarks[m.selected].ID
		return func() tea.Msg {
			err := m.database.RemoveBookmark(id)
			return BookmarkDeletedMsg{ID: id, Err: err}
		}
	}
	if m.tab == TabHighlights && m.selected < len(m.highlights) {
		verseID := m.highlights[m.selected].VerseID
		return func() tea.Msg {
			err := m.database.RemoveHighlight(verseID)
			return BookmarkDeletedMsg{ID: verseID, Err: err}
		}
	}
	return nil
}

func (m BookmarkModel) goToSelected() tea.Cmd {
	if m.tab == TabBookmarks && m.selected < len(m.bookmarks) {
		bm := m.bookmarks[m.selected]
		return func() tea.Msg {
			return GoToVerseMsg{BookCode: bm.BookCode, Chapter: bm.Chapter, Verse: bm.VerseNum}
		}
	}
	if m.tab == TabHighlights && m.selected < len(m.highlights) {
		hl := m.highlights[m.selected]
		return func() tea.Msg {
			return GoToVerseMsg{BookCode: hl.BookCode, Chapter: hl.Chapter, Verse: hl.VerseNum}
		}
	}
	return nil
}

func (m BookmarkModel) View() string {
	var b strings.Builder

	activeTab := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Primary).Underline(true)
	inactiveTab := lipgloss.NewStyle().Foreground(m.theme.Muted)

	bookmarkLabel := inactiveTab.Render("책갈피")
	highlightLabel := inactiveTab.Render("하이라이트")
	if m.tab == TabBookmarks {
		bookmarkLabel = activeTab.Render("책갈피")
	} else {
		highlightLabel = activeTab.Render("하이라이트")
	}
	b.WriteString(fmt.Sprintf("  %s  │  %s    (Tab:전환  d:삭제  Enter:이동)\n\n", bookmarkLabel, highlightLabel))

	if !m.loaded {
		b.WriteString("  로딩 중...")
		return b.String()
	}

	if m.tab == TabBookmarks {
		if len(m.bookmarks) == 0 {
			b.WriteString("  책갈피가 없습니다.")
			return b.String()
		}
		for i, bm := range m.bookmarks {
			cursor := "  "
			if i == m.selected {
				cursor = "▸ "
			}
			ref := fmt.Sprintf("%s %d:%d", bm.BookName, bm.Chapter, bm.VerseNum)
			refStyle := lipgloss.NewStyle().Foreground(m.theme.Primary).Bold(true)
			text := truncateRunes(bm.VerseText, m.width-15)
			line := fmt.Sprintf("%s%s — %s", cursor, refStyle.Render(ref), text)
			if bm.Note != "" {
				line += "\n    📝 " + bm.Note
			}
			b.WriteString(line + "\n")
		}
	} else {
		if len(m.highlights) == 0 {
			b.WriteString("  하이라이트가 없습니다.")
			return b.String()
		}
		for i, hl := range m.highlights {
			cursor := "  "
			if i == m.selected {
				cursor = "▸ "
			}
			ref := fmt.Sprintf("%s %d:%d", hl.BookName, hl.Chapter, hl.VerseNum)
			refStyle := lipgloss.NewStyle().Foreground(m.theme.Primary).Bold(true)
			colorTag := lipgloss.NewStyle().Foreground(lipgloss.Color(highlightColor(hl.Color))).Render("[" + hl.Color + "]")
			text := truncateRunes(hl.VerseText, m.width-20)
			b.WriteString(fmt.Sprintf("%s%s %s — %s\n", cursor, colorTag, refStyle.Render(ref), text))
		}
	}
	return b.String()
}

func truncateRunes(s string, maxLen int) string {
	if maxLen <= 0 {
		return s
	}
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + "..."
	}
	return s
}

func highlightColor(name string) string {
	switch name {
	case "yellow":
		return "#e0af68"
	case "green":
		return "#9ece6a"
	case "blue":
		return "#7aa2f7"
	case "pink":
		return "#f7768e"
	case "purple":
		return "#bb9af7"
	default:
		return "#e0af68"
	}
}
