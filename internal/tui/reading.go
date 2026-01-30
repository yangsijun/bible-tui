package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sijun-dong/bible-tui/internal/bible"
	"github.com/sijun-dong/bible-tui/internal/db"
	"github.com/sijun-dong/bible-tui/internal/tui/styles"
)

type VersesLoadedMsg struct {
	Verses []db.Verse
	Err    error
}

type ReadingModel struct {
	viewport viewport.Model
	book     bible.BookInfo
	chapter  int
	verses   []db.Verse
	loading  bool
	theme    *styles.Theme
	width    int
	height   int
}

func NewReading(book bible.BookInfo, chapter int, theme *styles.Theme, width, height int) ReadingModel {
	vp := viewport.New(width, height-4)
	vp.SetContent("로딩 중...")
	return ReadingModel{
		viewport: vp,
		book:     book,
		chapter:  chapter,
		loading:  true,
		theme:    theme,
		width:    width,
		height:   height,
	}
}

func LoadVerses(database *db.DB, bookCode string, chapter int) tea.Cmd {
	return func() tea.Msg {
		verses, err := database.GetVerses("GAE", bookCode, chapter)
		return VersesLoadedMsg{Verses: verses, Err: err}
	}
}

func (m ReadingModel) Update(msg tea.Msg) (ReadingModel, tea.Cmd) {
	switch msg := msg.(type) {
	case VersesLoadedMsg:
		m.loading = false
		if msg.Err != nil {
			m.viewport.SetContent(fmt.Sprintf("오류: %v", msg.Err))
			return m, nil
		}
		m.verses = msg.Verses
		m.viewport.SetContent(m.renderVerses())
		m.viewport.GotoTop()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			if m.chapter > 1 {
				return m, func() tea.Msg {
					return ChapterSelectedMsg{Book: m.book, Chapter: m.chapter - 1}
				}
			}
			return m, nil
		case "l", "right":
			if m.chapter < m.book.ChapterCount {
				return m, func() tea.Msg {
					return ChapterSelectedMsg{Book: m.book, Chapter: m.chapter + 1}
				}
			}
			return m, nil
		case "g":
			m.viewport.GotoTop()
			return m, nil
		case "G":
			m.viewport.GotoBottom()
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m ReadingModel) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Primary).Padding(0, 1)
	title := titleStyle.Render(fmt.Sprintf("%s %d장", m.book.NameKo, m.chapter))

	navHint := lipgloss.NewStyle().Foreground(m.theme.Muted).Render("  ←/h:이전장  →/l:다음장  j/k:스크롤  Esc:돌아가기")

	return title + navHint + "\n" + m.viewport.View()
}

func (m ReadingModel) renderVerses() string {
	var b strings.Builder
	numStyle := lipgloss.NewStyle().Foreground(m.theme.Muted).Width(4).Align(lipgloss.Right)
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Secondary)

	for _, v := range m.verses {
		if v.SectionTitle != "" {
			b.WriteString("\n")
			b.WriteString(titleStyle.Render(v.SectionTitle))
			b.WriteString("\n\n")
		}
		num := numStyle.Render(fmt.Sprintf("%d", v.VerseNum))
		b.WriteString(fmt.Sprintf("%s  %s\n", num, v.Text))
	}
	return b.String()
}
