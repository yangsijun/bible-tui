package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/sijun-dong/bible-tui/internal/db"
	"github.com/sijun-dong/bible-tui/internal/tui/styles"
)

func TestSearchModel_Init(t *testing.T) {
	m := NewSearch(nil, styles.DefaultDarkTheme(), 80, 24)
	if m.input.Value() != "" {
		t.Error("expected empty input")
	}
	if len(m.results) != 0 {
		t.Error("expected no results")
	}
}

func TestSearchModel_Results(t *testing.T) {
	m := NewSearch(nil, styles.DefaultDarkTheme(), 80, 24)
	results := []db.SearchResult{
		{Verse: db.Verse{BookName: "창세기", Chapter: 1, VerseNum: 1, Text: "태초에 하나님이", BookCode: "gen"}},
	}
	updated, _ := m.Update(SearchResultsMsg{Results: results, Query: "하나님"})
	if len(updated.results) != 1 {
		t.Errorf("expected 1, got %d", len(updated.results))
	}
	if updated.noResults {
		t.Error("expected noResults=false")
	}
}

func TestSearchModel_NoResults(t *testing.T) {
	m := NewSearch(nil, styles.DefaultDarkTheme(), 80, 24)
	m.query = "없는단어"
	updated, _ := m.Update(SearchResultsMsg{Results: []db.SearchResult{}, Query: "없는단어"})
	if !updated.noResults {
		t.Error("expected noResults=true")
	}
}

func TestSearchModel_Navigation(t *testing.T) {
	m := NewSearch(nil, styles.DefaultDarkTheme(), 80, 24)
	m.results = []db.SearchResult{
		{Verse: db.Verse{BookName: "창세기", Chapter: 1, VerseNum: 1, Text: "t1", BookCode: "gen"}},
		{Verse: db.Verse{BookName: "창세기", Chapter: 1, VerseNum: 2, Text: "t2", BookCode: "gen"}},
	}
	m.input.Blur()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if updated.selected != 1 {
		t.Errorf("expected 1, got %d", updated.selected)
	}

	updated2, _ := updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if updated2.selected != 0 {
		t.Errorf("expected 0, got %d", updated2.selected)
	}
}

func TestSearchModel_View(t *testing.T) {
	m := NewSearch(nil, styles.DefaultDarkTheme(), 80, 24)
	v := m.View()
	if v == "" {
		t.Error("expected non-empty view")
	}
}

func TestSearchModel_ViewWithResults(t *testing.T) {
	m := NewSearch(nil, styles.DefaultDarkTheme(), 80, 24)
	m.results = []db.SearchResult{
		{Verse: db.Verse{BookName: "창세기", Chapter: 1, VerseNum: 1, Text: "태초에 하나님이", BookCode: "gen"}},
	}
	v := m.View()
	if v == "" {
		t.Error("expected non-empty view")
	}
}
