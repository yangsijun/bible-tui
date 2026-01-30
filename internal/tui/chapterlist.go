package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sijun-dong/bible-tui/internal/bible"
	"github.com/sijun-dong/bible-tui/internal/tui/styles"
)

type ChapterSelectedMsg struct {
	Book    bible.BookInfo
	Chapter int
}

type ChapterListModel struct {
	book     bible.BookInfo
	selected int // 1-based current selection
	cols     int
	theme    *styles.Theme
	width    int
	height   int
}

func NewChapterList(book bible.BookInfo, theme *styles.Theme, width, height int) ChapterListModel {
	return ChapterListModel{
		book:     book,
		selected: 1,
		cols:     10,
		theme:    theme,
		width:    width,
		height:   height,
	}
}

func (m ChapterListModel) Update(msg tea.Msg) (ChapterListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right", "l":
			if m.selected < m.book.ChapterCount {
				m.selected++
			}
		case "left", "h":
			if m.selected > 1 {
				m.selected--
			}
		case "down", "j":
			next := m.selected + m.cols
			if next <= m.book.ChapterCount {
				m.selected = next
			}
		case "up", "k":
			prev := m.selected - m.cols
			if prev >= 1 {
				m.selected = prev
			}
		case "enter":
			return m, func() tea.Msg {
				return ChapterSelectedMsg{Book: m.book, Chapter: m.selected}
			}
		}
	}
	return m, nil
}

func (m ChapterListModel) View() string {
	var b strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Primary)
	b.WriteString(titleStyle.Render(fmt.Sprintf("%s — 장 선택", m.book.NameKo)))
	b.WriteString("\n\n")

	normalStyle := lipgloss.NewStyle().Width(5).Align(lipgloss.Center)
	selectedStyle := normalStyle.Background(m.theme.HighlightBg).Foreground(m.theme.Primary).Bold(true)

	for i := 1; i <= m.book.ChapterCount; i++ {
		style := normalStyle
		if i == m.selected {
			style = selectedStyle
		}
		b.WriteString(style.Render(fmt.Sprintf("%d", i)))
		if i%m.cols == 0 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
