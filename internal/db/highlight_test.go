package db

import (
	"testing"
)

func setupHighlightDB(t *testing.T) (*DB, int64) {
	t.Helper()
	db, err := OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	vID, err := db.InsertVersion("GAE", "개역개정", "ko")
	if err != nil {
		t.Fatalf("InsertVersion: %v", err)
	}
	bookID, err := db.InsertBook(vID, "gen", "창세기", "창", "old", 50, 0)
	if err != nil {
		t.Fatalf("InsertBook: %v", err)
	}
	verseID, err := db.InsertVerse(bookID, 1, 1, "태초에 하나님이 천지를 창조하시니라", "천지 창조", false)
	if err != nil {
		t.Fatalf("InsertVerse: %v", err)
	}
	return db, verseID
}

func TestAddAndGetHighlight(t *testing.T) {
	db, verseID := setupHighlightDB(t)

	if err := db.AddHighlight(verseID, "yellow"); err != nil {
		t.Fatalf("AddHighlight: %v", err)
	}

	color, err := db.GetHighlightColor(verseID)
	if err != nil {
		t.Fatalf("GetHighlightColor: %v", err)
	}
	if color != "yellow" {
		t.Errorf("expected color 'yellow', got %q", color)
	}
}

func TestHighlightUpsert(t *testing.T) {
	db, verseID := setupHighlightDB(t)

	if err := db.AddHighlight(verseID, "yellow"); err != nil {
		t.Fatalf("AddHighlight yellow: %v", err)
	}

	color, err := db.GetHighlightColor(verseID)
	if err != nil {
		t.Fatalf("GetHighlightColor after first add: %v", err)
	}
	if color != "yellow" {
		t.Errorf("expected color 'yellow', got %q", color)
	}

	if err := db.AddHighlight(verseID, "green"); err != nil {
		t.Fatalf("AddHighlight green: %v", err)
	}

	color, err = db.GetHighlightColor(verseID)
	if err != nil {
		t.Fatalf("GetHighlightColor after second add: %v", err)
	}
	if color != "green" {
		t.Errorf("expected color 'green', got %q", color)
	}
}

func TestRemoveHighlight(t *testing.T) {
	db, verseID := setupHighlightDB(t)

	if err := db.AddHighlight(verseID, "yellow"); err != nil {
		t.Fatalf("AddHighlight: %v", err)
	}

	color, err := db.GetHighlightColor(verseID)
	if err != nil {
		t.Fatalf("GetHighlightColor before remove: %v", err)
	}
	if color != "yellow" {
		t.Errorf("expected color 'yellow' before remove, got %q", color)
	}

	if err := db.RemoveHighlight(verseID); err != nil {
		t.Fatalf("RemoveHighlight: %v", err)
	}

	color, err = db.GetHighlightColor(verseID)
	if err != nil {
		t.Fatalf("GetHighlightColor after remove: %v", err)
	}
	if color != "" {
		t.Errorf("expected empty color after remove, got %q", color)
	}
}

func TestListHighlights(t *testing.T) {
	db, verseID := setupHighlightDB(t)

	if err := db.AddHighlight(verseID, "yellow"); err != nil {
		t.Fatalf("AddHighlight: %v", err)
	}

	vID, err := db.InsertVersion("HAN", "개역한글", "ko")
	if err != nil {
		t.Fatalf("InsertVersion 2: %v", err)
	}
	bookID, err := db.InsertBook(vID, "gen", "창세기", "창", "old", 50, 0)
	if err != nil {
		t.Fatalf("InsertBook 2: %v", err)
	}
	verseID2, err := db.InsertVerse(bookID, 1, 2, "땅이 혼돈하고 공허하며", "", false)
	if err != nil {
		t.Fatalf("InsertVerse 2: %v", err)
	}

	if err := db.AddHighlight(verseID2, "blue"); err != nil {
		t.Fatalf("AddHighlight 2: %v", err)
	}

	highlights, err := db.ListHighlights(10, 0)
	if err != nil {
		t.Fatalf("ListHighlights: %v", err)
	}
	if len(highlights) != 2 {
		t.Fatalf("expected 2 highlights, got %d", len(highlights))
	}

	if highlights[0].Color != "yellow" {
		t.Errorf("expected first highlight color 'yellow', got %q", highlights[0].Color)
	}
	if highlights[1].Color != "blue" {
		t.Errorf("expected second highlight color 'blue', got %q", highlights[1].Color)
	}

	if highlights[0].VerseText != "태초에 하나님이 천지를 창조하시니라" {
		t.Errorf("expected first verse text '태초에 하나님이 천지를 창조하시니라', got %q", highlights[0].VerseText)
	}
	if highlights[1].VerseText != "땅이 혼돈하고 공허하며" {
		t.Errorf("expected second verse text '땅이 혼돈하고 공허하며', got %q", highlights[1].VerseText)
	}
}

func TestGetHighlightColor_NoHighlight(t *testing.T) {
	db, verseID := setupHighlightDB(t)

	color, err := db.GetHighlightColor(verseID)
	if err != nil {
		t.Fatalf("GetHighlightColor: %v", err)
	}
	if color != "" {
		t.Errorf("expected empty color for non-highlighted verse, got %q", color)
	}
}
