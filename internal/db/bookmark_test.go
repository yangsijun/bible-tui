package db

import (
	"testing"
)

func setupBookmarkDB(t *testing.T) (*DB, int64) {
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

func TestAddAndListBookmarks(t *testing.T) {
	db, verseID := setupBookmarkDB(t)

	id, err := db.AddBookmark(verseID, "")
	if err != nil {
		t.Fatalf("AddBookmark: %v", err)
	}
	if id == 0 {
		t.Error("expected non-zero bookmark ID")
	}

	bookmarks, err := db.ListBookmarks(10, 0)
	if err != nil {
		t.Fatalf("ListBookmarks: %v", err)
	}
	if len(bookmarks) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(bookmarks))
	}

	bm := bookmarks[0]
	if bm.ID != id {
		t.Errorf("expected ID %d, got %d", id, bm.ID)
	}
	if bm.VerseID != verseID {
		t.Errorf("expected verse_id %d, got %d", verseID, bm.VerseID)
	}
	if bm.VerseText != "태초에 하나님이 천지를 창조하시니라" {
		t.Errorf("expected verse text '태초에 하나님이 천지를 창조하시니라', got %q", bm.VerseText)
	}
	if bm.BookName != "창세기" {
		t.Errorf("expected book name '창세기', got %q", bm.BookName)
	}
	if bm.BookCode != "gen" {
		t.Errorf("expected book code 'gen', got %q", bm.BookCode)
	}
	if bm.Chapter != 1 {
		t.Errorf("expected chapter 1, got %d", bm.Chapter)
	}
	if bm.VerseNum != 1 {
		t.Errorf("expected verse_num 1, got %d", bm.VerseNum)
	}
}

func TestRemoveBookmark(t *testing.T) {
	db, verseID := setupBookmarkDB(t)

	id, err := db.AddBookmark(verseID, "")
	if err != nil {
		t.Fatalf("AddBookmark: %v", err)
	}

	bookmarks, err := db.ListBookmarks(10, 0)
	if err != nil {
		t.Fatalf("ListBookmarks before remove: %v", err)
	}
	if len(bookmarks) != 1 {
		t.Fatalf("expected 1 bookmark before remove, got %d", len(bookmarks))
	}

	if err := db.RemoveBookmark(id); err != nil {
		t.Fatalf("RemoveBookmark: %v", err)
	}

	bookmarks, err = db.ListBookmarks(10, 0)
	if err != nil {
		t.Fatalf("ListBookmarks after remove: %v", err)
	}
	if len(bookmarks) != 0 {
		t.Fatalf("expected 0 bookmarks after remove, got %d", len(bookmarks))
	}
}

func TestIsBookmarked(t *testing.T) {
	db, verseID := setupBookmarkDB(t)

	isBookmarked, err := db.IsBookmarked(verseID)
	if err != nil {
		t.Fatalf("IsBookmarked before add: %v", err)
	}
	if isBookmarked {
		t.Error("expected false before adding bookmark")
	}

	_, err = db.AddBookmark(verseID, "")
	if err != nil {
		t.Fatalf("AddBookmark: %v", err)
	}

	isBookmarked, err = db.IsBookmarked(verseID)
	if err != nil {
		t.Fatalf("IsBookmarked after add: %v", err)
	}
	if !isBookmarked {
		t.Error("expected true after adding bookmark")
	}

	bookmarks, err := db.ListBookmarks(10, 0)
	if err != nil {
		t.Fatalf("ListBookmarks: %v", err)
	}
	if len(bookmarks) == 0 {
		t.Fatal("no bookmarks found")
	}

	if err := db.RemoveBookmark(bookmarks[0].ID); err != nil {
		t.Fatalf("RemoveBookmark: %v", err)
	}

	isBookmarked, err = db.IsBookmarked(verseID)
	if err != nil {
		t.Fatalf("IsBookmarked after remove: %v", err)
	}
	if isBookmarked {
		t.Error("expected false after removing bookmark")
	}
}

func TestBookmarkWithNote(t *testing.T) {
	db, verseID := setupBookmarkDB(t)

	note := "중요한 구절"
	id, err := db.AddBookmark(verseID, note)
	if err != nil {
		t.Fatalf("AddBookmark: %v", err)
	}

	bookmarks, err := db.ListBookmarks(10, 0)
	if err != nil {
		t.Fatalf("ListBookmarks: %v", err)
	}
	if len(bookmarks) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(bookmarks))
	}

	bm := bookmarks[0]
	if bm.ID != id {
		t.Errorf("expected ID %d, got %d", id, bm.ID)
	}
	if bm.Note != note {
		t.Errorf("expected note %q, got %q", note, bm.Note)
	}
}
