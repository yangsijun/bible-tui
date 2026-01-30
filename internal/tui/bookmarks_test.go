package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/sijun-dong/bible-tui/internal/db"
	"github.com/sijun-dong/bible-tui/internal/tui/styles"
)

func TestBookmarkModel_Init(t *testing.T) {
	m := NewBookmarks(nil, styles.DefaultDarkTheme(), 80, 24)
	if m.tab != TabBookmarks {
		t.Errorf("expected TabBookmarks")
	}
	if m.loaded {
		t.Error("expected not loaded")
	}
}

func TestBookmarkModel_Loaded(t *testing.T) {
	m := NewBookmarks(nil, styles.DefaultDarkTheme(), 80, 24)
	bookmarks := []db.BookmarkWithVerse{
		{Bookmark: db.Bookmark{ID: 1, VerseID: 1}, VerseText: "태초에", BookName: "창세기", Chapter: 1, VerseNum: 1},
	}
	updated, _ := m.Update(BookmarksLoadedMsg{Bookmarks: bookmarks, Highlights: nil})
	if !updated.loaded {
		t.Error("expected loaded")
	}
	if len(updated.bookmarks) != 1 {
		t.Errorf("expected 1 bookmark, got %d", len(updated.bookmarks))
	}
}

func TestBookmarkModel_TabSwitch(t *testing.T) {
	m := NewBookmarks(nil, styles.DefaultDarkTheme(), 80, 24)
	m.loaded = true
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	if updated.tab != TabHighlights {
		t.Errorf("expected TabHighlights")
	}
	updated2, _ := updated.Update(tea.KeyMsg{Type: tea.KeyTab})
	if updated2.tab != TabBookmarks {
		t.Errorf("expected TabBookmarks")
	}
}

func TestBookmarkModel_Navigation(t *testing.T) {
	m := NewBookmarks(nil, styles.DefaultDarkTheme(), 80, 24)
	m.loaded = true
	m.bookmarks = []db.BookmarkWithVerse{
		{Bookmark: db.Bookmark{ID: 1}, VerseText: "v1", BookName: "창세기", Chapter: 1, VerseNum: 1},
		{Bookmark: db.Bookmark{ID: 2}, VerseText: "v2", BookName: "창세기", Chapter: 1, VerseNum: 2},
	}
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if updated.selected != 1 {
		t.Errorf("expected 1, got %d", updated.selected)
	}
	updated2, _ := updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if updated2.selected != 0 {
		t.Errorf("expected 0, got %d", updated2.selected)
	}
}

func TestBookmarkModel_View(t *testing.T) {
	m := NewBookmarks(nil, styles.DefaultDarkTheme(), 80, 24)
	m.loaded = true
	v := m.View()
	if v == "" {
		t.Error("expected non-empty view")
	}
}

func TestBookmarkModel_ViewEmpty(t *testing.T) {
	m := NewBookmarks(nil, styles.DefaultDarkTheme(), 80, 24)
	m.loaded = true
	v := m.View()
	if !strings.Contains(v, "책갈피가 없습니다") {
		t.Error("expected empty bookmark message")
	}
}
