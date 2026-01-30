package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sijun-dong/bible-tui/internal/db"
	"github.com/sijun-dong/bible-tui/internal/tui/styles"
)

type AppState int

const (
	StateBookList AppState = iota
	StateChapterList
	StateReading
	StateSearch
	StateBookmarks
	StateSettings
	StatePlans
	StateHelp
)

type AppModel struct {
	state       AppState
	prevState   AppState
	db          *db.DB
	theme       *styles.Theme
	width       int
	height      int
	ready       bool
	statusMsg   string
	bookList    BookListModel
	chapterList ChapterListModel
	reading     ReadingModel
}

func New(database *db.DB) AppModel {
	return AppModel{
		state:    StateBookList,
		db:       database,
		theme:    styles.DefaultDarkTheme(),
		bookList: NewBookList(80, 24),
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		contentHeight := m.height - 3
		if contentHeight < 1 {
			contentHeight = 1
		}
		m.bookList.list.SetSize(msg.Width, contentHeight)
		return m, nil

	case BookSelectedMsg:
		contentHeight := m.height - 3
		if contentHeight < 1 {
			contentHeight = 1
		}
		m.chapterList = NewChapterList(msg.Book, m.theme, m.width, contentHeight)
		m.state = StateChapterList
		return m, nil

	case ChapterSelectedMsg:
		contentHeight := m.height - 3
		if contentHeight < 1 {
			contentHeight = 1
		}
		m.reading = NewReading(msg.Book, msg.Chapter, m.theme, m.width, contentHeight)
		m.state = StateReading
		if m.db != nil {
			return m, LoadVerses(m.db, msg.Book.Code, msg.Chapter)
		}
		return m, nil

	case VersesLoadedMsg:
		var cmd tea.Cmd
		m.reading, cmd = m.reading.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if m.state == StateBookList && m.bookList.list.SettingFilter() {
			var cmd tea.Cmd
			m.bookList, cmd = m.bookList.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.state != StateSearch {
				return m, tea.Quit
			}
		case "?":
			if m.state != StateHelp {
				m.prevState = m.state
				m.state = StateHelp
			}
			return m, nil
		case "esc":
			if m.state == StateHelp {
				m.state = m.prevState
			} else if m.state == StateReading {
				m.state = StateChapterList
			} else if m.state != StateBookList {
				m.state = StateBookList
			}
			return m, nil
		case "b":
			if m.state != StateBookList {
				m.state = StateBookList
			}
			return m, nil
		case "/":
			if m.state != StateSearch {
				m.prevState = m.state
				m.state = StateSearch
			}
			return m, nil
		case "m":
			if m.state != StateBookmarks {
				m.prevState = m.state
				m.state = StateBookmarks
			}
			return m, nil
		}

		switch m.state {
		case StateBookList:
			var cmd tea.Cmd
			m.bookList, cmd = m.bookList.Update(msg)
			return m, cmd
		case StateChapterList:
			var cmd tea.Cmd
			m.chapterList, cmd = m.chapterList.Update(msg)
			return m, cmd
		case StateReading:
			var cmd tea.Cmd
			m.reading, cmd = m.reading.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m AppModel) View() string {
	if !m.ready {
		return "Loading..."
	}

	var content string
	switch m.state {
	case StateBookList:
		content = m.bookList.View()
	case StateChapterList:
		content = m.chapterList.View()
	case StateReading:
		content = m.reading.View()
	case StateHelp:
		content = "도움말\n\nq, Ctrl+C  종료\n?          도움말\nb          책 목록\n/          검색\nm          책갈피\nEsc        이전 화면"
	case StateSearch:
		content = "검색 (구현 예정)"
	case StateBookmarks:
		content = "책갈피 (구현 예정)"
	default:
		content = m.bookList.View()
	}

	return m.renderLayout(content)
}

func (m AppModel) renderLayout(content string) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.theme.Primary).
		Padding(0, 1)
	header := headerStyle.Render("성경 Bible TUI")

	statusStyle := lipgloss.NewStyle().
		Background(m.theme.StatusBarBg).
		Foreground(m.theme.StatusBarFg).
		Width(m.width).
		Padding(0, 1)

	label := m.stateLabel()
	hints := "q:종료 ?:도움말 b:책목록 /:검색 m:책갈피"
	statusBar := statusStyle.Render(fmt.Sprintf("%s  │  %s", label, hints))

	contentHeight := m.height - 3
	if contentHeight < 1 {
		contentHeight = 1
	}
	contentStyle := lipgloss.NewStyle().
		Height(contentHeight).
		Width(m.width)

	return header + "\n" + contentStyle.Render(content) + "\n" + statusBar
}

func (m AppModel) stateLabel() string {
	switch m.state {
	case StateBookList:
		return "책 목록"
	case StateChapterList:
		return "장 선택"
	case StateReading:
		return "읽기"
	case StateSearch:
		return "검색"
	case StateBookmarks:
		return "책갈피"
	case StateSettings:
		return "설정"
	case StatePlans:
		return "읽기 계획"
	case StateHelp:
		return "도움말"
	default:
		return ""
	}
}
