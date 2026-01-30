package db

import (
	"strings"
	"testing"
)

func setupSearchDB(t *testing.T) *DB {
	db, err := OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	vID, err := db.InsertVersion("GAE", "개역개정", "ko")
	if err != nil {
		t.Fatalf("InsertVersion: %v", err)
	}

	genID, err := db.InsertBook(vID, "gen", "창세기", "창", "old", 50, 0)
	if err != nil {
		t.Fatalf("InsertBook gen: %v", err)
	}

	jhnID, err := db.InsertBook(vID, "jhn", "요한복음", "요", "new", 21, 39)
	if err != nil {
		t.Fatalf("InsertBook jhn: %v", err)
	}

	verses := []struct {
		bookID       int64
		chapter      int
		verseNum     int
		text         string
		sectionTitle string
	}{
		{genID, 1, 1, "태초에 하나님이 천지를 창조하시니라", "천지 창조"},
		{genID, 1, 2, "땅이 혼돈하고 공허하며 흑암이 깊음 위에 있고 하나님의 영은 수면 위에 운행하시니라", ""},
		{genID, 1, 3, "하나님이 이르시되 빛이 있으라 하시니 빛이 있었고", ""},
		{genID, 1, 4, "하나님이 빛을 보시니 좋았더라 하나님이 빛과 어둠을 나누사", ""},
		{jhnID, 3, 16, "하나님이 세상을 이처럼 사랑하사 독생자를 주셨으니 이는 그를 믿는 자마다 멸망하지 않고 영생을 얻게 하려 하심이라", ""},
	}

	for _, v := range verses {
		_, err := db.InsertVerse(v.bookID, v.chapter, v.verseNum, v.text, v.sectionTitle, false)
		if err != nil {
			t.Fatalf("InsertVerse: %v", err)
		}
	}

	return db
}

func TestSearchVerses_FTS(t *testing.T) {
	db := setupSearchDB(t)
	defer db.Close()

	results, err := db.SearchVerses("GAE", "하나님", 10)
	if err != nil {
		t.Fatalf("SearchVerses: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected results, got none")
	}

	expectedCount := 5
	if len(results) != expectedCount {
		t.Errorf("expected %d results, got %d", expectedCount, len(results))
	}

	for _, r := range results {
		if !strings.Contains(r.Verse.Text, "하나님") {
			t.Errorf("result text does not contain '하나님': %s", r.Verse.Text)
		}
		if r.Snippet == "" {
			t.Errorf("snippet is empty for verse %d", r.Verse.ID)
		}
		if r.MatchCount == 0 {
			t.Errorf("match count is 0 for verse %d", r.Verse.ID)
		}
	}
}

func TestSearchVerses_LIKE(t *testing.T) {
	db := setupSearchDB(t)
	defer db.Close()

	results, err := db.SearchVerses("GAE", "빛", 10)
	if err != nil {
		t.Fatalf("SearchVerses: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected results, got none")
	}

	for _, r := range results {
		if !strings.Contains(r.Verse.Text, "빛") {
			t.Errorf("result text does not contain '빛': %s", r.Verse.Text)
		}
		if r.Snippet == "" {
			t.Errorf("snippet is empty for verse %d", r.Verse.ID)
		}
	}
}

func TestSearchVerses_NoResult(t *testing.T) {
	db := setupSearchDB(t)
	defer db.Close()

	results, err := db.SearchVerses("GAE", "존재하지않는단어", 10)
	if err != nil {
		t.Fatalf("SearchVerses: %v", err)
	}

	if results == nil {
		t.Fatal("expected empty slice, got nil")
	}

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchVerses_EmptyQuery(t *testing.T) {
	db := setupSearchDB(t)
	defer db.Close()

	results, err := db.SearchVerses("GAE", "", 10)
	if err != nil {
		t.Fatalf("SearchVerses: %v", err)
	}

	if results == nil {
		t.Fatal("expected empty slice, got nil")
	}

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchVerses_Limit(t *testing.T) {
	db := setupSearchDB(t)
	defer db.Close()

	results, err := db.SearchVerses("GAE", "하나님", 2)
	if err != nil {
		t.Fatalf("SearchVerses: %v", err)
	}

	if len(results) > 2 {
		t.Errorf("expected max 2 results, got %d", len(results))
	}
}

func TestSearchVerses_Snippet(t *testing.T) {
	db := setupSearchDB(t)
	defer db.Close()

	results, err := db.SearchVerses("GAE", "하나님", 10)
	if err != nil {
		t.Fatalf("SearchVerses: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected results, got none")
	}

	for _, r := range results {
		if r.Snippet == "" {
			t.Errorf("snippet is empty for verse %d", r.Verse.ID)
		}
		lowerSnippet := strings.ToLower(r.Snippet)
		lowerQuery := strings.ToLower("하나님")
		if !strings.Contains(lowerSnippet, lowerQuery) && !strings.Contains(r.Snippet, "**") {
			t.Errorf("snippet does not contain search term or markers: %s", r.Snippet)
		}
	}
}

func TestSearchVerses_MatchCount(t *testing.T) {
	db := setupSearchDB(t)
	defer db.Close()

	results, err := db.SearchVerses("GAE", "하나님", 10)
	if err != nil {
		t.Fatalf("SearchVerses: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected results, got none")
	}

	for _, r := range results {
		actualCount := strings.Count(strings.ToLower(r.Verse.Text), strings.ToLower("하나님"))
		if r.MatchCount != actualCount {
			t.Errorf("verse %d: expected match count %d, got %d (text: %s)",
				r.Verse.ID, actualCount, r.MatchCount, r.Verse.Text)
		}
	}

	for _, r := range results {
		if r.Verse.Chapter == 1 && r.Verse.VerseNum == 4 {
			if r.MatchCount != 2 {
				t.Errorf("Genesis 1:4 should have 2 matches, got %d", r.MatchCount)
			}
		}
	}
}

func TestSearchVerses_Deduplication(t *testing.T) {
	db := setupSearchDB(t)
	defer db.Close()

	results, err := db.SearchVerses("GAE", "하나님", 10)
	if err != nil {
		t.Fatalf("SearchVerses: %v", err)
	}

	seen := make(map[int64]bool)
	for _, r := range results {
		if seen[r.Verse.ID] {
			t.Errorf("duplicate verse ID %d in results", r.Verse.ID)
		}
		seen[r.Verse.ID] = true
	}
}

func TestGenerateSnippet(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		query        string
		contextChars int
		wantContains string
	}{
		{
			name:         "basic snippet",
			text:         "태초에 하나님이 천지를 창조하시니라",
			query:        "하나님",
			contextChars: 10,
			wantContains: "하나님",
		},
		{
			name:         "snippet at start",
			text:         "하나님이 세상을 이처럼 사랑하사",
			query:        "하나님",
			contextChars: 10,
			wantContains: "하나님",
		},
		{
			name:         "snippet at end",
			text:         "태초에 하나님",
			query:        "하나님",
			contextChars: 10,
			wantContains: "하나님",
		},
		{
			name:         "no match",
			text:         "태초에 천지를 창조하시니라",
			query:        "하나님",
			contextChars: 10,
			wantContains: "태초",
		},
		{
			name:         "case insensitive",
			text:         "The quick brown fox",
			query:        "QUICK",
			contextChars: 5,
			wantContains: "quick",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snippet := generateSnippet(tt.text, tt.query, tt.contextChars)
			if snippet == "" {
				t.Error("snippet is empty")
			}
			if !strings.Contains(strings.ToLower(snippet), strings.ToLower(tt.wantContains)) {
				t.Errorf("snippet %q does not contain %q", snippet, tt.wantContains)
			}
		})
	}
}

func TestGenerateSnippet_Truncation(t *testing.T) {
	text := "이것은 매우 긴 텍스트입니다. 하나님이 세상을 이처럼 사랑하사 독생자를 주셨으니 이는 그를 믿는 자마다 멸망하지 않고 영생을 얻게 하려 하심이라"
	snippet := generateSnippet(text, "하나님", 10)

	if !strings.Contains(snippet, "...") {
		t.Error("expected ellipsis in truncated snippet")
	}
	if !strings.Contains(snippet, "하나님") {
		t.Error("snippet does not contain search term")
	}
}

func TestCountMatches(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		query string
		want  int
	}{
		{
			name:  "single match",
			text:  "태초에 하나님이 천지를 창조하시니라",
			query: "하나님",
			want:  1,
		},
		{
			name:  "multiple matches",
			text:  "하나님이 빛을 보시니 좋았더라 하나님이 빛과 어둠을 나누사",
			query: "하나님",
			want:  2,
		},
		{
			name:  "no match",
			text:  "태초에 천지를 창조하시니라",
			query: "하나님",
			want:  0,
		},
		{
			name:  "case insensitive",
			text:  "The Quick Brown Fox",
			query: "quick",
			want:  1,
		},
		{
			name:  "overlapping not counted twice",
			text:  "aaa",
			query: "aa",
			want:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countMatches(tt.text, tt.query)
			if got != tt.want {
				t.Errorf("countMatches(%q, %q) = %d, want %d", tt.text, tt.query, got, tt.want)
			}
		})
	}
}

func TestSearchVerses_SpecialCharacters(t *testing.T) {
	db := setupSearchDB(t)
	defer db.Close()

	results, err := db.SearchVerses("GAE", "!", 10)
	if err != nil {
		t.Fatalf("SearchVerses with special char: %v", err)
	}

	if results == nil {
		t.Fatal("expected empty slice, got nil")
	}
}

func TestSearchVerses_BookMetadata(t *testing.T) {
	db := setupSearchDB(t)
	defer db.Close()

	results, err := db.SearchVerses("GAE", "하나님", 10)
	if err != nil {
		t.Fatalf("SearchVerses: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected results, got none")
	}

	for _, r := range results {
		if r.Verse.BookName == "" {
			t.Errorf("verse %d has empty BookName", r.Verse.ID)
		}
		if r.Verse.BookCode == "" {
			t.Errorf("verse %d has empty BookCode", r.Verse.ID)
		}
		if r.Verse.Chapter == 0 {
			t.Errorf("verse %d has zero Chapter", r.Verse.ID)
		}
		if r.Verse.VerseNum == 0 {
			t.Errorf("verse %d has zero VerseNum", r.Verse.ID)
		}
	}
}
