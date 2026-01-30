package db

import (
	"testing"
)

func setupTestDB(t *testing.T) *DB {
	t.Helper()
	d, err := OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	if err := d.Migrate(); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	t.Cleanup(func() { d.Close() })
	return d
}

func seedTestData(t *testing.T, d *DB) (versionID, bookID int64) {
	t.Helper()
	versionID, err := d.InsertVersion("GAE", "개역개정", "ko")
	if err != nil {
		t.Fatalf("InsertVersion: %v", err)
	}
	bookID, err = d.InsertBook(versionID, "gen", "창세기", "창", "old", 50, 1)
	if err != nil {
		t.Fatalf("InsertBook: %v", err)
	}
	return versionID, bookID
}

func TestMigrate(t *testing.T) {
	d := setupTestDB(t)

	expectedTables := []string{
		"versions", "books", "verses", "footnotes",
		"bookmarks", "highlights", "reading_plans", "reading_plan_entries",
		"settings", "crawl_status", "verses_fts",
	}

	for _, table := range expectedTables {
		var name string
		err := d.conn.QueryRow(
			"SELECT name FROM sqlite_master WHERE type IN ('table', 'shadow') AND name = ?",
			table,
		).Scan(&name)
		if err != nil {
			t.Errorf("table %q not found: %v", table, err)
		}
	}

	expectedTriggers := []string{"verses_ai", "verses_ad", "verses_au"}
	for _, trig := range expectedTriggers {
		var name string
		err := d.conn.QueryRow(
			"SELECT name FROM sqlite_master WHERE type = 'trigger' AND name = ?",
			trig,
		).Scan(&name)
		if err != nil {
			t.Errorf("trigger %q not found: %v", trig, err)
		}
	}
}

func TestInsertAndGetVerses(t *testing.T) {
	d := setupTestDB(t)
	_, bookID := seedTestData(t, d)

	texts := []string{
		"태초에 하나님이 천지를 창조하시니라",
		"땅이 혼돈하고 공허하며 흑암이 깊음 위에 있고",
		"하나님이 이르시되 빛이 있으라 하시니 빛이 있었고",
	}
	for i, text := range texts {
		_, err := d.InsertVerse(bookID, 1, i+1, text, "", false)
		if err != nil {
			t.Fatalf("InsertVerse(%d): %v", i+1, err)
		}
	}

	verses, err := d.GetVerses("GAE", "gen", 1)
	if err != nil {
		t.Fatalf("GetVerses: %v", err)
	}
	if len(verses) != 3 {
		t.Fatalf("expected 3 verses, got %d", len(verses))
	}
	for i, v := range verses {
		if v.Text != texts[i] {
			t.Errorf("verse %d: expected %q, got %q", i+1, texts[i], v.Text)
		}
		if v.VerseNum != i+1 {
			t.Errorf("verse %d: expected verse_num %d, got %d", i+1, i+1, v.VerseNum)
		}
		if v.BookName != "창세기" {
			t.Errorf("verse %d: expected book name 창세기, got %q", i+1, v.BookName)
		}
		if v.BookCode != "gen" {
			t.Errorf("verse %d: expected book code gen, got %q", i+1, v.BookCode)
		}
	}
}

func TestGetRandomVerse(t *testing.T) {
	d := setupTestDB(t)
	_, bookID := seedTestData(t, d)

	_, err := d.InsertVerse(bookID, 1, 1, "태초에 하나님이 천지를 창조하시니라", "", false)
	if err != nil {
		t.Fatalf("InsertVerse: %v", err)
	}
	_, err = d.InsertVerse(bookID, 1, 2, "땅이 혼돈하고 공허하며", "", false)
	if err != nil {
		t.Fatalf("InsertVerse: %v", err)
	}

	v, err := d.GetRandomVerse("GAE")
	if err != nil {
		t.Fatalf("GetRandomVerse: %v", err)
	}
	if v == nil {
		t.Fatal("expected non-nil verse")
	}
	if v.BookName != "창세기" {
		t.Errorf("expected book name 창세기, got %q", v.BookName)
	}
	if v.Text == "" {
		t.Error("expected non-empty text")
	}
}

func TestSettings(t *testing.T) {
	d := setupTestDB(t)

	val, err := d.GetSetting("nonexistent")
	if err != nil {
		t.Fatalf("GetSetting nonexistent: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty string for nonexistent key, got %q", val)
	}

	if err := d.SetSetting("theme", "dark"); err != nil {
		t.Fatalf("SetSetting: %v", err)
	}
	val, err = d.GetSetting("theme")
	if err != nil {
		t.Fatalf("GetSetting: %v", err)
	}
	if val != "dark" {
		t.Errorf("expected 'dark', got %q", val)
	}

	if err := d.SetSetting("theme", "light"); err != nil {
		t.Fatalf("SetSetting overwrite: %v", err)
	}
	val, err = d.GetSetting("theme")
	if err != nil {
		t.Fatalf("GetSetting after overwrite: %v", err)
	}
	if val != "light" {
		t.Errorf("expected 'light', got %q", val)
	}
}

func TestInsertFootnote(t *testing.T) {
	d := setupTestDB(t)
	_, bookID := seedTestData(t, d)

	verseID, err := d.InsertVerse(bookID, 1, 1, "태초에 하나님이 천지를 창조하시니라", "", true)
	if err != nil {
		t.Fatalf("InsertVerse: %v", err)
	}

	if err := d.InsertFootnote(verseID, "1)", "히브리어 원문 해석"); err != nil {
		t.Fatalf("InsertFootnote: %v", err)
	}

	var fn Footnote
	err = d.conn.QueryRow(
		"SELECT id, verse_id, marker, content FROM footnotes WHERE verse_id = ?",
		verseID,
	).Scan(&fn.ID, &fn.VerseID, &fn.Marker, &fn.Content)
	if err != nil {
		t.Fatalf("query footnote: %v", err)
	}
	if fn.VerseID != verseID {
		t.Errorf("expected verse_id %d, got %d", verseID, fn.VerseID)
	}
	if fn.Marker != "1)" {
		t.Errorf("expected marker '1)', got %q", fn.Marker)
	}
	if fn.Content != "히브리어 원문 해석" {
		t.Errorf("expected content '히브리어 원문 해석', got %q", fn.Content)
	}
}
