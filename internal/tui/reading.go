package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/yangsijun/bible-tui/internal/bible"
	"github.com/yangsijun/bible-tui/internal/db"
	"github.com/yangsijun/bible-tui/internal/tui/styles"
)

type VersesLoadedMsg struct {
	Verses []db.Verse
	Err    error
}

type ReadingModel struct {
	viewport    viewport.Model
	book        bible.BookInfo
	chapter     int
	verses      []db.Verse
	loading     bool
	theme       *styles.Theme
	width       int
	height      int
	cursorIdx   int
	database    *db.DB
	statusMsg   string
	statusTimer int
	lineOffsets []int // line offset for each verse in rendered content
}

func NewReading(book bible.BookInfo, chapter int, database *db.DB, theme *styles.Theme, width, height int) ReadingModel {
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
		database: database,
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
		m.cursorIdx = 0
		m.viewport.SetContent(m.renderVerses())
		m.viewport.GotoTop()
		return m, nil
	case tea.KeyMsg:
		if m.statusMsg != "" {
			m.statusMsg = ""
			m.statusTimer = 0
		}
		switch msg.String() {
		case "j", "down":
			if len(m.verses) > 0 && m.cursorIdx < len(m.verses)-1 {
				m.cursorIdx++
				m.viewport.SetContent(m.renderVerses())
				m.ensureCursorVisible()
			}
			return m, nil
		case "k", "up":
			if len(m.verses) > 0 && m.cursorIdx > 0 {
				m.cursorIdx--
				m.viewport.SetContent(m.renderVerses())
				m.ensureCursorVisible()
			}
			return m, nil
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
			m.cursorIdx = 0
			m.viewport.SetContent(m.renderVerses())
			m.viewport.GotoTop()
			return m, nil
		case "G":
			if len(m.verses) > 0 {
				m.cursorIdx = len(m.verses) - 1
				m.viewport.SetContent(m.renderVerses())
				m.viewport.GotoBottom()
			}
			return m, nil
		case "B":
			if len(m.verses) > 0 && m.database != nil {
				v := m.verses[m.cursorIdx]
				_, err := m.database.AddBookmark(v.ID, "")
				if err == nil {
					m.statusMsg = fmt.Sprintf("책갈피 추가: %s %d:%d", m.book.NameKo, v.Chapter, v.VerseNum)
				} else {
					m.statusMsg = fmt.Sprintf("오류: %v", err)
				}
			}
			return m, nil
		case "H":
			if len(m.verses) > 0 && m.database != nil {
				v := m.verses[m.cursorIdx]
				err := m.database.AddHighlight(v.ID, "yellow")
				if err == nil {
					m.statusMsg = fmt.Sprintf("하이라이트 추가: %s %d:%d", m.book.NameKo, v.Chapter, v.VerseNum)
				} else {
					m.statusMsg = fmt.Sprintf("오류: %v", err)
				}
			}
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *ReadingModel) ensureCursorVisible() {
	if m.cursorIdx < 0 || m.cursorIdx >= len(m.lineOffsets) {
		return
	}
	cursorLine := m.lineOffsets[m.cursorIdx]
	vpHeight := m.viewport.Height
	yOffset := m.viewport.YOffset

	if cursorLine < yOffset {
		m.viewport.SetYOffset(cursorLine)
	} else if cursorLine >= yOffset+vpHeight {
		m.viewport.SetYOffset(cursorLine - vpHeight + 1)
	}
}

func (m ReadingModel) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Primary).Padding(0, 1)
	title := titleStyle.Render(fmt.Sprintf("%s %d장", m.book.NameKo, m.chapter))

	navHint := lipgloss.NewStyle().Foreground(m.theme.Muted).Render("  ←/h:이전장  →/l:다음장  j/k:구절이동  B:책갈피  H:하이라이트  Esc:돌아가기")

	header := title + navHint
	if m.statusMsg != "" {
		statusStyle := lipgloss.NewStyle().Foreground(m.theme.Secondary).Bold(true).PaddingLeft(2)
		header += statusStyle.Render("  " + m.statusMsg)
	}

	return header + "\n" + m.viewport.View()
}

func (m *ReadingModel) renderVerses() string {
	var b strings.Builder
	numStyle := lipgloss.NewStyle().Foreground(m.theme.Muted).Width(4).Align(lipgloss.Right)
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Secondary)
	cursorStyle := lipgloss.NewStyle().Foreground(m.theme.Primary).Bold(true)

	m.lineOffsets = make([]int, len(m.verses))
	lineCount := 0

	for i, v := range m.verses {
		if v.SectionTitle != "" {
			b.WriteString("\n")
			lineCount++
			b.WriteString(titleStyle.Render(v.SectionTitle))
			b.WriteString("\n\n")
			lineCount += 2
		}

		m.lineOffsets[i] = lineCount

		marker := "  "
		if i == m.cursorIdx {
			marker = cursorStyle.Render("▸ ")
		}

		num := numStyle.Render(fmt.Sprintf("%d", v.VerseNum))
		b.WriteString(fmt.Sprintf("%s%s  %s\n", marker, num, v.Text))
		lineCount++
	}
	return b.String()
}
