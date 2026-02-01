package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yangsijun/bible-tui/internal/db"
	"github.com/yangsijun/bible-tui/internal/tui/styles"
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

func TestTryParseReference_FullRef(t *testing.T) {
	msg := tryParseReference("창세기 2")
	if msg == nil {
		t.Fatal("expected GoToVerseMsg for '창세기 2'")
	}
	if msg.BookCode != "gen" || msg.Chapter != 2 || msg.Verse != 1 {
		t.Errorf("expected gen/2/1, got %s/%d/%d", msg.BookCode, msg.Chapter, msg.Verse)
	}
}

func TestTryParseReference_WithVerse(t *testing.T) {
	msg := tryParseReference("창 3:3")
	if msg == nil {
		t.Fatal("expected GoToVerseMsg for '창 3:3'")
	}
	if msg.BookCode != "gen" || msg.Chapter != 3 || msg.Verse != 3 {
		t.Errorf("expected gen/3/3, got %s/%d/%d", msg.BookCode, msg.Chapter, msg.Verse)
	}
}

func TestTryParseReference_BookOnly(t *testing.T) {
	msg := tryParseReference("창세기")
	if msg == nil {
		t.Fatal("expected GoToVerseMsg for '창세기'")
	}
	if msg.BookCode != "gen" || msg.Chapter != 1 {
		t.Errorf("expected gen/1, got %s/%d", msg.BookCode, msg.Chapter)
	}
}

func TestTryParseReference_AbbrevOnly(t *testing.T) {
	msg := tryParseReference("창")
	if msg == nil {
		t.Fatal("expected GoToVerseMsg for '창'")
	}
	if msg.BookCode != "gen" {
		t.Errorf("expected gen, got %s", msg.BookCode)
	}
}

func TestTryParseReference_TextQuery(t *testing.T) {
	msg := tryParseReference("사랑")
	if msg != nil {
		t.Error("expected nil for plain text query '사랑'")
	}
}

func TestTryParseReference_EnglishCode(t *testing.T) {
	msg := tryParseReference("gen 1")
	if msg == nil {
		t.Fatal("expected GoToVerseMsg for 'gen 1'")
	}
	if msg.BookCode != "gen" || msg.Chapter != 1 {
		t.Errorf("expected gen/1, got %s/%d", msg.BookCode, msg.Chapter)
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
