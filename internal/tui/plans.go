package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/yangsijun/bible-tui/internal/db"
	"github.com/yangsijun/bible-tui/internal/tui/styles"
)

type PlansLoadedMsg struct {
	Plans []db.ReadingPlan
	Err   error
}

type PlanEntriesLoadedMsg struct {
	Entries  []db.PlanEntry
	Progress PlanProgress
	Err      error
}

type PlanProgress struct {
	Completed int
	Total     int
}

type PlanCreatedMsg struct {
	PlanID int64
	Err    error
}

type EntryCompletedMsg struct {
	Err error
}

type PlanDeletedMsg struct {
	Err error
}

type PlanViewState int

const (
	PlanViewList   PlanViewState = iota
	PlanViewToday
	PlanViewCreate
)

type PlanModel struct {
	database  *db.DB
	theme     *styles.Theme
	width     int
	height    int

	viewState PlanViewState
	plans     []db.ReadingPlan
	entries   []db.PlanEntry
	progress  PlanProgress
	selected  int
	loaded    bool

	createIdx int
	versionID int64
}

func NewPlans(database *db.DB, theme *styles.Theme, width, height int) PlanModel {
	return PlanModel{
		database:  database,
		theme:     theme,
		width:     width,
		height:    height,
		viewState: PlanViewList,
		versionID: 1,
	}
}

func LoadPlans(database *db.DB) tea.Cmd {
	return func() tea.Msg {
		if database == nil {
			return PlansLoadedMsg{Err: fmt.Errorf("no database")}
		}
		plans, err := database.ListPlans()
		return PlansLoadedMsg{Plans: plans, Err: err}
	}
}

func LoadPlanEntries(database *db.DB, planID int64) tea.Cmd {
	return func() tea.Msg {
		if database == nil {
			return PlanEntriesLoadedMsg{Err: fmt.Errorf("no database")}
		}
		entries, err := database.GetTodayEntries(planID)
		if err != nil {
			return PlanEntriesLoadedMsg{Err: err}
		}
		completed, total, err := database.GetPlanProgress(planID)
		if err != nil {
			return PlanEntriesLoadedMsg{Err: err}
		}
		return PlanEntriesLoadedMsg{
			Entries:  entries,
			Progress: PlanProgress{Completed: completed, Total: total},
		}
	}
}

func createPlan(database *db.DB, planType string, versionID int64) tea.Cmd {
	return func() tea.Msg {
		if database == nil {
			return PlanCreatedMsg{Err: fmt.Errorf("no database")}
		}
		var id int64
		var err error
		switch planType {
		case "sequential":
			id, err = database.CreateSequentialPlan(versionID, "통독 계획")
		case "mcheyne":
			id, err = database.CreateMcCheynePlan(versionID, "매쿠인 계획")
		default:
			err = fmt.Errorf("unknown plan type: %s", planType)
		}
		return PlanCreatedMsg{PlanID: id, Err: err}
	}
}

func deletePlan(database *db.DB, planID int64) tea.Cmd {
	return func() tea.Msg {
		if database == nil {
			return PlanDeletedMsg{Err: fmt.Errorf("no database")}
		}
		err := database.DeletePlan(planID)
		return PlanDeletedMsg{Err: err}
	}
}

func markEntryCompleted(database *db.DB, entryID int64) tea.Cmd {
	return func() tea.Msg {
		if database == nil {
			return EntryCompletedMsg{Err: fmt.Errorf("no database")}
		}
		err := database.MarkEntryCompleted(entryID)
		return EntryCompletedMsg{Err: err}
	}
}

func (m PlanModel) Update(msg tea.Msg) (PlanModel, tea.Cmd) {
	switch msg := msg.(type) {

	case PlansLoadedMsg:
		m.loaded = true
		if msg.Err == nil {
			m.plans = msg.Plans
		}
		m.selected = 0
		return m, nil

	case PlanEntriesLoadedMsg:
		if msg.Err == nil {
			m.entries = msg.Entries
			m.progress = msg.Progress
		}
		m.selected = 0
		return m, nil

	case PlanCreatedMsg:
		if msg.Err == nil {
			return m, LoadPlans(m.database)
		}
		return m, nil

	case EntryCompletedMsg:
		if msg.Err == nil && m.viewState == PlanViewToday && len(m.plans) > 0 {
			planID := m.activePlanID()
			if planID > 0 {
				return m, LoadPlanEntries(m.database, planID)
			}
		}
		return m, nil

	case PlanDeletedMsg:
		if msg.Err == nil {
			return m, LoadPlans(m.database)
		}
		return m, nil

	case tea.KeyMsg:
		switch m.viewState {
		case PlanViewList:
			return m.updateList(msg)
		case PlanViewToday:
			return m.updateToday(msg)
		case PlanViewCreate:
			return m.updateCreate(msg)
		}
	}
	return m, nil
}

func (m PlanModel) updateList(msg tea.KeyMsg) (PlanModel, tea.Cmd) {
	switch msg.String() {
	case "down", "j":
		if m.selected < len(m.plans)-1 {
			m.selected++
		}
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
	case "enter":
		if m.selected < len(m.plans) && len(m.plans) > 0 {
			m.viewState = PlanViewToday
			planID := m.plans[m.selected].ID
			m.selected = 0
			return m, LoadPlanEntries(m.database, planID)
		}
	case "n":
		m.viewState = PlanViewCreate
		m.createIdx = 0
	case "d":
		if m.selected < len(m.plans) && len(m.plans) > 0 {
			planID := m.plans[m.selected].ID
			return m, deletePlan(m.database, planID)
		}
	}
	return m, nil
}

func (m PlanModel) updateToday(msg tea.KeyMsg) (PlanModel, tea.Cmd) {
	switch msg.String() {
	case "down", "j":
		if m.selected < len(m.entries)-1 {
			m.selected++
		}
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
	case " ":
		return m, m.toggleEntry()
	case "enter":
		return m, m.goToEntry()
	case "esc":
		m.viewState = PlanViewList
		m.selected = 0
	}
	return m, nil
}

func (m PlanModel) updateCreate(msg tea.KeyMsg) (PlanModel, tea.Cmd) {
	switch msg.String() {
	case "down", "j":
		if m.createIdx < 1 {
			m.createIdx++
		}
	case "up", "k":
		if m.createIdx > 0 {
			m.createIdx--
		}
	case "enter":
		m.viewState = PlanViewList
		planType := "sequential"
		if m.createIdx == 1 {
			planType = "mcheyne"
		}
		return m, createPlan(m.database, planType, m.versionID)
	case "esc":
		m.viewState = PlanViewList
		m.selected = 0
	}
	return m, nil
}

func (m PlanModel) activePlanID() int64 {
	if len(m.entries) > 0 {
		return m.entries[0].PlanID
	}
	return 0
}

func (m PlanModel) toggleEntry() tea.Cmd {
	if m.selected < len(m.entries) {
		e := m.entries[m.selected]
		if !e.Completed {
			return markEntryCompleted(m.database, e.ID)
		}
	}
	return nil
}

func (m PlanModel) goToEntry() tea.Cmd {
	if m.selected < len(m.entries) {
		e := m.entries[m.selected]
		return func() tea.Msg {
			return GoToVerseMsg{BookCode: e.BookCode, Chapter: e.ChapterStart, Verse: 1}
		}
	}
	return nil
}

func (m PlanModel) View() string {
	switch m.viewState {
	case PlanViewToday:
		return m.viewToday()
	case PlanViewCreate:
		return m.viewCreate()
	default:
		return m.viewList()
	}
}

func (m PlanModel) viewList() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Primary)
	b.WriteString("  " + titleStyle.Render("읽기 계획") + "\n\n")

	if !m.loaded {
		b.WriteString("  로딩 중...")
		return b.String()
	}

	if len(m.plans) == 0 {
		mutedStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
		b.WriteString("  " + mutedStyle.Render("읽기 계획이 없습니다.") + "\n")
	} else {
		for i, p := range m.plans {
			cursor := "  "
			if i == m.selected {
				cursor = "▸ "
			}
			nameStyle := lipgloss.NewStyle().Foreground(m.theme.Foreground)
			if i == m.selected {
				nameStyle = nameStyle.Bold(true)
			}
			line := fmt.Sprintf("%s%s (%d일)", cursor, nameStyle.Render(p.Name), p.TotalDays)
			b.WriteString(line + "\n")
		}
	}

	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
	b.WriteString("  " + helpStyle.Render("n:새 계획  Enter:오늘 읽기  d:삭제"))

	return b.String()
}

func (m PlanModel) viewToday() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Primary)

	planName := "읽기 계획"
	dayNumber := 0
	if len(m.entries) > 0 {
		for _, p := range m.plans {
			if p.ID == m.entries[0].PlanID {
				planName = p.Name
				break
			}
		}
		dayNumber = m.entries[0].DayNumber
	}

	pct := 0
	if m.progress.Total > 0 {
		pct = m.progress.Completed * 100 / m.progress.Total
	}

	header := fmt.Sprintf("%s — %d일차 (%d%%)", planName, dayNumber, pct)
	b.WriteString("  " + titleStyle.Render(header) + "\n")

	barWidth := 20
	if m.width > 30 {
		barWidth = m.width / 4
	}
	bar := renderProgressBar(m.progress.Completed, m.progress.Total, barWidth)
	if bar != "" {
		barStyle := lipgloss.NewStyle().Foreground(m.theme.Secondary)
		b.WriteString("  " + barStyle.Render(bar) + "\n")
	}
	b.WriteString("\n")

	if len(m.entries) == 0 {
		mutedStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
		b.WriteString("  " + mutedStyle.Render("오늘의 읽기가 없습니다.") + "\n")
	} else {
		for i, e := range m.entries {
			cursor := "  "
			if i == m.selected {
				cursor = "▸ "
			}
			check := "[ ]"
			if e.Completed {
				check = "[✓]"
			}

			ref := formatEntryRef(e)
			checkStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
			if e.Completed {
				checkStyle = lipgloss.NewStyle().Foreground(m.theme.SectionTitle)
			}
			refStyle := lipgloss.NewStyle().Foreground(m.theme.Foreground)
			if i == m.selected {
				refStyle = refStyle.Bold(true)
			}

			b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, checkStyle.Render(check), refStyle.Render(ref)))
		}
	}

	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
	b.WriteString("  " + helpStyle.Render("Space:완료  Enter:읽기  Esc:목록"))

	return b.String()
}

func (m PlanModel) viewCreate() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Primary)
	b.WriteString("  " + titleStyle.Render("새 읽기 계획 만들기") + "\n\n")

	options := []struct {
		name string
		desc string
	}{
		{"통독", "창세기→요한계시록, 1일 3장"},
		{"매쿠인", "1일 4구간, 1년 완독"},
	}

	for i, opt := range options {
		cursor := "  "
		if i == m.createIdx {
			cursor = "▸ "
		}
		nameStyle := lipgloss.NewStyle().Foreground(m.theme.Foreground).Bold(i == m.createIdx)
		descStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
		b.WriteString(fmt.Sprintf("%s%s (%s)\n", cursor, nameStyle.Render(opt.name), descStyle.Render(opt.desc)))
	}

	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(m.theme.Muted)
	b.WriteString("  " + helpStyle.Render("Enter:생성  Esc:취소"))

	return b.String()
}

func renderProgressBar(completed, total, width int) string {
	if total == 0 {
		return ""
	}
	if width <= 0 {
		width = 10
	}
	pct := completed * 100 / total
	filled := completed * width / total
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("%s %d/%d (%d%%)", bar, completed, total, pct)
}

func formatEntryRef(e db.PlanEntry) string {
	if e.ChapterStart == e.ChapterEnd {
		return fmt.Sprintf("%s %d장", e.BookCode, e.ChapterStart)
	}
	return fmt.Sprintf("%s %d-%d장", e.BookCode, e.ChapterStart, e.ChapterEnd)
}
