package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sijun-dong/bible-tui/internal/bible"
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
	search      SearchModel
	bookmarks   BookmarkModel
	help        HelpModel
	settings    SettingsModel
	plans       PlanModel
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

	case SearchResultsMsg:
		var cmd tea.Cmd
		m.search, cmd = m.search.Update(msg)
		return m, cmd

	case BookmarksLoadedMsg:
		var cmd tea.Cmd
		m.bookmarks, cmd = m.bookmarks.Update(msg)
		return m, cmd

	case BookmarkDeletedMsg:
		var cmd tea.Cmd
		m.bookmarks, cmd = m.bookmarks.Update(msg)
		return m, cmd

	case SettingsLoadedMsg:
		var cmd tea.Cmd
		m.settings, cmd = m.settings.Update(msg)
		return m, cmd

	case SettingsSavedMsg:
		var cmd tea.Cmd
		m.settings, cmd = m.settings.Update(msg)
		m.state = m.prevState
		return m, cmd

	case ThemeChangedMsg:
		m.theme = msg.Theme
		return m, nil

	case PlansLoadedMsg:
		var cmd tea.Cmd
		m.plans, cmd = m.plans.Update(msg)
		return m, cmd

	case PlanEntriesLoadedMsg:
		var cmd tea.Cmd
		m.plans, cmd = m.plans.Update(msg)
		return m, cmd

	case PlanCreatedMsg:
		var cmd tea.Cmd
		m.plans, cmd = m.plans.Update(msg)
		return m, cmd

	case EntryCompletedMsg:
		var cmd tea.Cmd
		m.plans, cmd = m.plans.Update(msg)
		return m, cmd

	case GoToVerseMsg:
		book := findBookByCode(msg.BookCode)
		if book != nil {
			contentHeight := m.height - 3
			if contentHeight < 1 {
				contentHeight = 1
			}
			m.reading = NewReading(*book, msg.Chapter, m.theme, m.width, contentHeight)
			m.state = StateReading
			if m.db != nil {
				return m, LoadVerses(m.db, msg.BookCode, msg.Chapter)
			}
		}
		return m, nil

	case tea.KeyMsg:
		if m.state == StateBookList && m.bookList.list.SettingFilter() {
			var cmd tea.Cmd
			m.bookList, cmd = m.bookList.Update(msg)
			return m, cmd
		}

		if m.state == StateSearch && m.search.input.Focused() {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.state = m.prevState
				return m, nil
			default:
				var cmd tea.Cmd
				m.search, cmd = m.search.Update(msg)
				return m, cmd
			}
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
				contentHeight := m.height - 3
				if contentHeight < 1 {
					contentHeight = 1
				}
				m.help = NewHelp(m.theme, m.width, contentHeight)
			}
			return m, nil
		case "esc":
			switch m.state {
			case StateHelp, StateSearch, StateBookmarks, StateSettings:
				m.state = m.prevState
			case StatePlans:
				m.state = m.prevState
			case StateReading:
				m.state = StateChapterList
			case StateChapterList:
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
				contentHeight := m.height - 3
				if contentHeight < 1 {
					contentHeight = 1
				}
				m.search = NewSearch(m.db, m.theme, m.width, contentHeight)
				return m, m.search.input.Focus()
			}
			return m, nil
		case "m":
			if m.state != StateBookmarks {
				m.prevState = m.state
				m.state = StateBookmarks
				contentHeight := m.height - 3
				if contentHeight < 1 {
					contentHeight = 1
				}
				m.bookmarks = NewBookmarks(m.db, m.theme, m.width, contentHeight)
				return m, LoadBookmarks(m.db)
			}
			return m, nil
		case "s":
			if m.state != StateSettings && m.state != StateSearch {
				m.prevState = m.state
				m.state = StateSettings
				contentHeight := m.height - 3
				if contentHeight < 1 {
					contentHeight = 1
				}
				m.settings = NewSettings(m.db, m.theme, m.width, contentHeight)
				return m, LoadSettings(m.db)
			}
			return m, nil
		case "p":
			if m.state != StatePlans {
				m.prevState = m.state
				m.state = StatePlans
				contentHeight := m.height - 3
				if contentHeight < 1 {
					contentHeight = 1
				}
				m.plans = NewPlans(m.db, m.theme, m.width, contentHeight)
				return m, LoadPlans(m.db)
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
		case StateSearch:
			var cmd tea.Cmd
			m.search, cmd = m.search.Update(msg)
			return m, cmd
		case StateBookmarks:
			var cmd tea.Cmd
			m.bookmarks, cmd = m.bookmarks.Update(msg)
			return m, cmd
		case StateHelp:
			var cmd tea.Cmd
			m.help, cmd = m.help.Update(msg)
			return m, cmd
		case StateSettings:
			var cmd tea.Cmd
			m.settings, cmd = m.settings.Update(msg)
			return m, cmd
		case StatePlans:
			var cmd tea.Cmd
			m.plans, cmd = m.plans.Update(msg)
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
		content = m.help.View()
	case StateSearch:
		content = m.search.View()
	case StateBookmarks:
		content = m.bookmarks.View()
	case StateSettings:
		content = m.settings.View()
	case StatePlans:
		content = m.plans.View()
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
	hints := "q:종료 ?:도움말 b:책목록 /:검색 m:책갈피 s:설정 p:읽기계획"
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

func findBookByCode(code string) *bible.BookInfo {
	for _, b := range bible.AllBooks() {
		if b.Code == code {
			return &b
		}
	}
	return nil
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
