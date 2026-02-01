package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yangsijun/bible-tui/internal/bible"
	"github.com/yangsijun/bible-tui/internal/tui/styles"
)

func TestChapterListInit(t *testing.T) {
	book := bible.BookInfo{Code: "gen", NameKo: "창세기", ChapterCount: 50}
	m := NewChapterList(book, styles.DefaultDarkTheme(), 80, 24)
	if m.selected != 1 {
		t.Errorf("expected selected=1, got %d", m.selected)
	}
}

func TestChapterListNavigation(t *testing.T) {
	book := bible.BookInfo{Code: "gen", NameKo: "창세기", ChapterCount: 50}
	m := NewChapterList(book, styles.DefaultDarkTheme(), 80, 24)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	if updated.selected != 2 {
		t.Errorf("right: expected 2, got %d", updated.selected)
	}

	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	if updated.selected != 1 {
		t.Errorf("left: expected 1, got %d", updated.selected)
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if updated.selected != 11 {
		t.Errorf("down: expected 11, got %d", updated.selected)
	}

	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if updated.selected != 1 {
		t.Errorf("up: expected 1, got %d", updated.selected)
	}
}

func TestChapterListBoundary(t *testing.T) {
	book := bible.BookInfo{Code: "gen", NameKo: "창세기", ChapterCount: 3}
	m := NewChapterList(book, styles.DefaultDarkTheme(), 80, 24)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	if updated.selected != 1 {
		t.Errorf("left at min: expected 1, got %d", updated.selected)
	}

	m.selected = 3
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	if updated.selected != 3 {
		t.Errorf("right at max: expected 3, got %d", updated.selected)
	}
}

func TestChapterListView(t *testing.T) {
	book := bible.BookInfo{Code: "gen", NameKo: "창세기", ChapterCount: 50}
	m := NewChapterList(book, styles.DefaultDarkTheme(), 80, 24)
	v := m.View()
	if v == "" {
		t.Error("expected non-empty view")
	}
}
