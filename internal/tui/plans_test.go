package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/sijun-dong/bible-tui/internal/db"
	"github.com/sijun-dong/bible-tui/internal/tui/styles"
)

func TestPlanModel_Init(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	if m.viewState != PlanViewList {
		t.Errorf("expected PlanViewList, got %d", m.viewState)
	}
	if m.loaded {
		t.Error("expected not loaded")
	}
	if m.selected != 0 {
		t.Errorf("expected selected=0, got %d", m.selected)
	}
}

func TestPlanModel_PlansLoaded(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m, _ = m.Update(PlansLoadedMsg{
		Plans: []db.ReadingPlan{
			{ID: 1, Name: "통독 계획", PlanType: "sequential", TotalDays: 397},
			{ID: 2, Name: "매쿠인 계획", PlanType: "mcheyne", TotalDays: 365},
		},
	})
	if !m.loaded {
		t.Error("expected loaded")
	}
	if len(m.plans) != 2 {
		t.Errorf("expected 2 plans, got %d", len(m.plans))
	}
	if m.selected != 0 {
		t.Errorf("expected selected reset to 0, got %d", m.selected)
	}
}

func TestPlanModel_PlansLoadedError(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m, _ = m.Update(PlansLoadedMsg{Err: fmt.Errorf("db error")})
	if !m.loaded {
		t.Error("expected loaded even on error")
	}
	if len(m.plans) != 0 {
		t.Errorf("expected 0 plans on error, got %d", len(m.plans))
	}
}

func TestPlanModel_Navigation(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m, _ = m.Update(PlansLoadedMsg{
		Plans: []db.ReadingPlan{
			{ID: 1, Name: "plan1", TotalDays: 10},
			{ID: 2, Name: "plan2", TotalDays: 20},
			{ID: 3, Name: "plan3", TotalDays: 30},
		},
	})

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.selected != 1 {
		t.Errorf("after j: expected 1, got %d", m.selected)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.selected != 2 {
		t.Errorf("after j: expected 2, got %d", m.selected)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.selected != 2 {
		t.Errorf("after j at bottom: expected 2, got %d", m.selected)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.selected != 1 {
		t.Errorf("after k: expected 1, got %d", m.selected)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.selected != 0 {
		t.Errorf("after k: expected 0, got %d", m.selected)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.selected != 0 {
		t.Errorf("after k at top: expected 0, got %d", m.selected)
	}
}

func TestPlanModel_SwitchToCreate(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m.loaded = true
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if m.viewState != PlanViewCreate {
		t.Errorf("expected PlanViewCreate, got %d", m.viewState)
	}
	if m.createIdx != 0 {
		t.Errorf("expected createIdx=0, got %d", m.createIdx)
	}
}

func TestPlanModel_CreateNavigation(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m.viewState = PlanViewCreate
	m.createIdx = 0

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.createIdx != 1 {
		t.Errorf("after j: expected createIdx=1, got %d", m.createIdx)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.createIdx != 1 {
		t.Errorf("after j at bottom: expected createIdx=1, got %d", m.createIdx)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.createIdx != 0 {
		t.Errorf("after k: expected createIdx=0, got %d", m.createIdx)
	}
}

func TestPlanModel_BackFromCreate(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m.viewState = PlanViewCreate

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if m.viewState != PlanViewList {
		t.Errorf("expected PlanViewList, got %d", m.viewState)
	}
}

func TestPlanModel_ViewList(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m, _ = m.Update(PlansLoadedMsg{
		Plans: []db.ReadingPlan{
			{ID: 1, Name: "통독 계획", TotalDays: 397},
		},
	})
	v := m.View()
	if !strings.Contains(v, "읽기 계획") {
		t.Error("expected view to contain '읽기 계획'")
	}
	if !strings.Contains(v, "통독 계획") {
		t.Error("expected view to contain plan name")
	}
	if !strings.Contains(v, "397일") {
		t.Error("expected view to contain total days")
	}
}

func TestPlanModel_ViewListEmpty(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m, _ = m.Update(PlansLoadedMsg{Plans: nil})
	v := m.View()
	if !strings.Contains(v, "읽기 계획이 없습니다") {
		t.Error("expected empty message")
	}
}

func TestPlanModel_ViewToday(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m, _ = m.Update(PlansLoadedMsg{
		Plans: []db.ReadingPlan{
			{ID: 1, Name: "통독 계획", TotalDays: 397},
		},
	})
	m.viewState = PlanViewToday
	m.entries = []db.PlanEntry{
		{ID: 1, PlanID: 1, DayNumber: 42, BookCode: "gen", ChapterStart: 10, ChapterEnd: 12, Completed: true},
		{ID: 2, PlanID: 1, DayNumber: 42, BookCode: "gen", ChapterStart: 13, ChapterEnd: 15, Completed: false},
	}
	m.progress = PlanProgress{Completed: 310, Total: 397}

	v := m.View()
	if !strings.Contains(v, "█") && !strings.Contains(v, "░") {
		t.Error("expected progress bar with █ or ░")
	}
	if !strings.Contains(v, "310/397") {
		t.Error("expected progress numbers")
	}
	if !strings.Contains(v, "[✓]") {
		t.Error("expected completed check mark")
	}
	if !strings.Contains(v, "[ ]") {
		t.Error("expected uncompleted check")
	}
}

func TestPlanModel_ViewCreate(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m.viewState = PlanViewCreate
	v := m.View()
	if !strings.Contains(v, "새 읽기 계획 만들기") {
		t.Error("expected create title")
	}
	if !strings.Contains(v, "통독") {
		t.Error("expected sequential option")
	}
	if !strings.Contains(v, "매쿠인") {
		t.Error("expected mcheyne option")
	}
}

func TestPlanModel_SpaceToggle(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m.viewState = PlanViewToday
	m.entries = []db.PlanEntry{
		{ID: 10, PlanID: 1, DayNumber: 1, BookCode: "gen", ChapterStart: 1, ChapterEnd: 3, Completed: false},
		{ID: 11, PlanID: 1, DayNumber: 1, BookCode: "gen", ChapterStart: 4, ChapterEnd: 6, Completed: false},
	}
	m.selected = 0

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	if cmd == nil {
		t.Error("expected non-nil cmd from Space on uncompleted entry")
	}
}

func TestPlanModel_SpaceToggleCompleted(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m.viewState = PlanViewToday
	m.entries = []db.PlanEntry{
		{ID: 10, PlanID: 1, DayNumber: 1, BookCode: "gen", ChapterStart: 1, ChapterEnd: 3, Completed: true},
	}
	m.selected = 0

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	if cmd != nil {
		t.Error("expected nil cmd from Space on already completed entry")
	}
}

func TestPlanModel_EnterGoToVerse(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m.viewState = PlanViewToday
	m.entries = []db.PlanEntry{
		{ID: 10, PlanID: 1, DayNumber: 1, BookCode: "gen", ChapterStart: 5, ChapterEnd: 7},
	}
	m.selected = 0

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected non-nil cmd from Enter on entry")
	}
	msg := cmd()
	goTo, ok := msg.(GoToVerseMsg)
	if !ok {
		t.Fatalf("expected GoToVerseMsg, got %T", msg)
	}
	if goTo.BookCode != "gen" {
		t.Errorf("expected BookCode=gen, got %s", goTo.BookCode)
	}
	if goTo.Chapter != 5 {
		t.Errorf("expected Chapter=5, got %d", goTo.Chapter)
	}
}

func TestPlanModel_EscFromToday(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m.viewState = PlanViewToday
	m.entries = []db.PlanEntry{
		{ID: 1, PlanID: 1, DayNumber: 1, BookCode: "gen", ChapterStart: 1, ChapterEnd: 3},
	}
	m.selected = 0

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if m.viewState != PlanViewList {
		t.Errorf("expected PlanViewList, got %d", m.viewState)
	}
}

func TestPlanModel_TodayNavigation(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m.viewState = PlanViewToday
	m.entries = []db.PlanEntry{
		{ID: 1, PlanID: 1, DayNumber: 1, BookCode: "gen", ChapterStart: 1, ChapterEnd: 3},
		{ID: 2, PlanID: 1, DayNumber: 1, BookCode: "gen", ChapterStart: 4, ChapterEnd: 6},
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.selected != 1 {
		t.Errorf("expected selected=1, got %d", m.selected)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.selected != 0 {
		t.Errorf("expected selected=0, got %d", m.selected)
	}
}

func TestPlanModel_ProgressBar(t *testing.T) {
	bar := renderProgressBar(50, 100, 10)
	if !strings.Contains(bar, "█") {
		t.Error("expected filled blocks")
	}
	if !strings.Contains(bar, "░") {
		t.Error("expected empty blocks")
	}
	if !strings.Contains(bar, "50/100") {
		t.Error("expected progress numbers")
	}
	if !strings.Contains(bar, "50%") {
		t.Error("expected percentage")
	}
}

func TestPlanModel_ProgressBarEmpty(t *testing.T) {
	bar := renderProgressBar(0, 0, 10)
	if bar != "" {
		t.Errorf("expected empty string for zero total, got %q", bar)
	}
}

func TestPlanModel_ProgressBarFull(t *testing.T) {
	bar := renderProgressBar(100, 100, 10)
	if !strings.Contains(bar, "100%") {
		t.Error("expected 100%")
	}
	if strings.Contains(bar, "░") {
		t.Error("expected no empty blocks at 100%")
	}
}

func TestPlanModel_FormatEntryRef(t *testing.T) {
	e1 := db.PlanEntry{BookCode: "gen", ChapterStart: 5, ChapterEnd: 5}
	ref1 := formatEntryRef(e1)
	if ref1 != "gen 5장" {
		t.Errorf("expected 'gen 5장', got %q", ref1)
	}

	e2 := db.PlanEntry{BookCode: "gen", ChapterStart: 10, ChapterEnd: 12}
	ref2 := formatEntryRef(e2)
	if ref2 != "gen 10-12장" {
		t.Errorf("expected 'gen 10-12장', got %q", ref2)
	}
}

func TestPlanModel_EntriesLoaded(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	m.viewState = PlanViewToday
	m, _ = m.Update(PlanEntriesLoadedMsg{
		Entries: []db.PlanEntry{
			{ID: 1, PlanID: 1, DayNumber: 5, BookCode: "gen", ChapterStart: 13, ChapterEnd: 15},
		},
		Progress: PlanProgress{Completed: 12, Total: 397},
	})
	if len(m.entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(m.entries))
	}
	if m.progress.Completed != 12 {
		t.Errorf("expected completed=12, got %d", m.progress.Completed)
	}
	if m.progress.Total != 397 {
		t.Errorf("expected total=397, got %d", m.progress.Total)
	}
}

func TestPlanModel_ViewNotLoaded(t *testing.T) {
	m := NewPlans(nil, styles.DefaultDarkTheme(), 80, 24)
	v := m.View()
	if !strings.Contains(v, "로딩 중") {
		t.Error("expected loading message")
	}
}
