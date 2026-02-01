package tui

import (
	"fmt"
	"testing"

	"github.com/yangsijun/bible-tui/internal/bible"
	"github.com/yangsijun/bible-tui/internal/db"
	"github.com/yangsijun/bible-tui/internal/tui/styles"
)

func TestReadingModel_VersesLoaded(t *testing.T) {
	book := bible.BookInfo{Code: "gen", NameKo: "창세기", ChapterCount: 50}
	m := NewReading(book, 1, styles.DefaultDarkTheme(), 80, 24)
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
	m := NewReading(book, 1, styles.DefaultDarkTheme(), 80, 24)
	updated, _ := m.Update(VersesLoadedMsg{Err: fmt.Errorf("db error")})
	if updated.loading {
		t.Error("expected loading=false after error")
	}
}

func TestReadingModel_View(t *testing.T) {
	book := bible.BookInfo{Code: "gen", NameKo: "창세기", ChapterCount: 50}
	m := NewReading(book, 1, styles.DefaultDarkTheme(), 80, 24)
	v := m.View()
	if v == "" {
		t.Error("expected non-empty view")
	}
}

func TestReadingModel_InitialLoading(t *testing.T) {
	book := bible.BookInfo{Code: "gen", NameKo: "창세기", ChapterCount: 50}
	m := NewReading(book, 1, styles.DefaultDarkTheme(), 80, 24)
	if !m.loading {
		t.Error("expected loading=true initially")
	}
}
