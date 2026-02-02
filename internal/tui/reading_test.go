package tui

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yangsijun/bible-tui/internal/bible"
	"github.com/yangsijun/bible-tui/internal/db"
	"github.com/yangsijun/bible-tui/internal/tui/styles"
)

func TestReadingModel_VersesLoaded(t *testing.T) {
	book := bible.BookInfo{Code: "gen", NameKo: "창세기", ChapterCount: 50}
	m := NewReading(book, 1, nil, styles.DefaultDarkTheme(), 80, 24)
	verses := []db.Verse{
		{VerseNum: 1, Text: "태초에 하나님이 천지를 창조하시니라", SectionTitle: "천지 창조"},
	}
	updated, _ := m.Update(VersesLoadedMsg{Verses: verses})
	if updated.loading {
		t.Error("expected loading=false")
	}
	if len(updated.verses) != 1 {
		t.Errorf("expected 1 verse, got %d", len(updated.verses))
	}
}

func TestReadingModel_VersesLoadedError(t *testing.T) {
	book := bible.BookInfo{Code: "gen", NameKo: "창세기", ChapterCount: 50}
	m := NewReading(book, 1, nil, styles.DefaultDarkTheme(), 80, 24)
	updated, _ := m.Update(VersesLoadedMsg{Err: fmt.Errorf("db error")})
	if updated.loading {
		t.Error("expected loading=false after error")
	}
}

func TestReadingModel_View(t *testing.T) {
	book := bible.BookInfo{Code: "gen", NameKo: "창세기", ChapterCount: 50}
	m := NewReading(book, 1, nil, styles.DefaultDarkTheme(), 80, 24)
	v := m.View()
	if v == "" {
		t.Error("expected non-empty view")
	}
}

func TestReadingModel_InitialLoading(t *testing.T) {
	book := bible.BookInfo{Code: "gen", NameKo: "창세기", ChapterCount: 50}
	m := NewReading(book, 1, nil, styles.DefaultDarkTheme(), 80, 24)
	if !m.loading {
		t.Error("expected loading=true initially")
	}
}

func TestReadingModel_CursorMovement(t *testing.T) {
	book := bible.BookInfo{Code: "gen", NameKo: "창세기", ChapterCount: 50}
	m := NewReading(book, 1, nil, styles.DefaultDarkTheme(), 80, 24)

	verses := []db.Verse{
		{ID: 1, VerseNum: 1, Text: "첫째 구절", Chapter: 1},
		{ID: 2, VerseNum: 2, Text: "둘째 구절", Chapter: 1},
		{ID: 3, VerseNum: 3, Text: "셋째 구절", Chapter: 1},
	}
	m, _ = m.Update(VersesLoadedMsg{Verses: verses})

	if m.cursorIdx != 0 {
		t.Errorf("expected cursorIdx=0 after load, got %d", m.cursorIdx)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.cursorIdx != 1 {
		t.Errorf("expected cursorIdx=1 after j, got %d", m.cursorIdx)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.cursorIdx != 2 {
		t.Errorf("expected cursorIdx=2 after j, got %d", m.cursorIdx)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.cursorIdx != 2 {
		t.Errorf("expected cursorIdx=2 (clamped) after j, got %d", m.cursorIdx)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.cursorIdx != 1 {
		t.Errorf("expected cursorIdx=1 after k, got %d", m.cursorIdx)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.cursorIdx != 0 {
		t.Errorf("expected cursorIdx=0 after k, got %d", m.cursorIdx)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.cursorIdx != 0 {
		t.Errorf("expected cursorIdx=0 (clamped) after k, got %d", m.cursorIdx)
	}
}

func TestReadingModel_CursorGAndShiftG(t *testing.T) {
	book := bible.BookInfo{Code: "gen", NameKo: "창세기", ChapterCount: 50}
	m := NewReading(book, 1, nil, styles.DefaultDarkTheme(), 80, 24)

	verses := []db.Verse{
		{ID: 1, VerseNum: 1, Text: "첫째 구절", Chapter: 1},
		{ID: 2, VerseNum: 2, Text: "둘째 구절", Chapter: 1},
		{ID: 3, VerseNum: 3, Text: "셋째 구절", Chapter: 1},
	}
	m, _ = m.Update(VersesLoadedMsg{Verses: verses})

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	if m.cursorIdx != 2 {
		t.Errorf("expected cursorIdx=2 after G, got %d", m.cursorIdx)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	if m.cursorIdx != 0 {
		t.Errorf("expected cursorIdx=0 after g, got %d", m.cursorIdx)
	}
}

func TestReadingModel_BookmarkWithoutDB(t *testing.T) {
	book := bible.BookInfo{Code: "gen", NameKo: "창세기", ChapterCount: 50}
	m := NewReading(book, 1, nil, styles.DefaultDarkTheme(), 80, 24)

	verses := []db.Verse{
		{ID: 1, VerseNum: 1, Text: "첫째 구절", Chapter: 1},
	}
	m, _ = m.Update(VersesLoadedMsg{Verses: verses})

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'B'}})
	if m.statusMsg != "" {
		t.Errorf("expected empty statusMsg without DB, got %q", m.statusMsg)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'H'}})
	if m.statusMsg != "" {
		t.Errorf("expected empty statusMsg without DB, got %q", m.statusMsg)
	}
}
