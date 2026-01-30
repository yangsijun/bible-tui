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
	state     AppState
	prevState AppState
	db        *db.DB
	theme     *styles.Theme
	width     int
	height    int
	ready     bool
	statusMsg string
}

func New(database *db.DB) AppModel {
	return AppModel{
		state: StateBookList,
		db:    database,
		theme: styles.DefaultDarkTheme(),
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
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "?":
			if m.state != StateHelp {
				m.prevState = m.state
				m.state = StateHelp
			}
			return m, nil
		case "esc":
			if m.state == StateHelp {
				m.state = m.prevState
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
		content = "책 목록 (구현 예정)"
	case StateHelp:
		content = "도움말\n\nq, Ctrl+C  종료\n?          도움말\nb          책 목록\n/          검색\nm          책갈피\nEsc        이전 화면"
	case StateSearch:
		content = "검색 (구현 예정)"
	case StateBookmarks:
		content = "책갈피 (구현 예정)"
	default:
		content = "책 목록 (구현 예정)"
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
