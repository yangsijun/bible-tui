package bible

import (
	"testing"
)

func TestParseReference_FullName(t *testing.T) {
	ref, err := ParseReference("창세기 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.BookCode != "gen" {
		t.Errorf("expected BookCode 'gen', got '%s'", ref.BookCode)
	}
	if ref.Chapter != 1 {
		t.Errorf("expected Chapter 1, got %d", ref.Chapter)
	}
	if ref.VerseStart != 0 {
		t.Errorf("expected VerseStart 0, got %d", ref.VerseStart)
	}
	if ref.VerseEnd != 0 {
		t.Errorf("expected VerseEnd 0, got %d", ref.VerseEnd)
	}
}

func TestParseReference_Abbrev(t *testing.T) {
	ref, err := ParseReference("창 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.BookCode != "gen" {
		t.Errorf("expected BookCode 'gen', got '%s'", ref.BookCode)
	}
	if ref.Chapter != 1 {
		t.Errorf("expected Chapter 1, got %d", ref.Chapter)
	}
}

func TestParseReference_English(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"gen 1", "gen"},
		{"Gen 1", "gen"},
		{"GEN 1", "gen"},
	}

	for _, tt := range tests {
		ref, err := ParseReference(tt.input)
		if err != nil {
			t.Fatalf("unexpected error for '%s': %v", tt.input, err)
		}
		if ref.BookCode != tt.expected {
			t.Errorf("for input '%s', expected BookCode '%s', got '%s'", tt.input, tt.expected, ref.BookCode)
		}
		if ref.Chapter != 1 {
			t.Errorf("for input '%s', expected Chapter 1, got %d", tt.input, ref.Chapter)
		}
	}
}

func TestParseReference_WithVerse(t *testing.T) {
	ref, err := ParseReference("창세기 1:3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.BookCode != "gen" {
		t.Errorf("expected BookCode 'gen', got '%s'", ref.BookCode)
	}
	if ref.Chapter != 1 {
		t.Errorf("expected Chapter 1, got %d", ref.Chapter)
	}
	if ref.VerseStart != 3 {
		t.Errorf("expected VerseStart 3, got %d", ref.VerseStart)
	}
	if ref.VerseEnd != 0 {
		t.Errorf("expected VerseEnd 0, got %d", ref.VerseEnd)
	}
}

func TestParseReference_Range(t *testing.T) {
	ref, err := ParseReference("창 1:3-5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.BookCode != "gen" {
		t.Errorf("expected BookCode 'gen', got '%s'", ref.BookCode)
	}
	if ref.Chapter != 1 {
		t.Errorf("expected Chapter 1, got %d", ref.Chapter)
	}
	if ref.VerseStart != 3 {
		t.Errorf("expected VerseStart 3, got %d", ref.VerseStart)
	}
	if ref.VerseEnd != 5 {
		t.Errorf("expected VerseEnd 5, got %d", ref.VerseEnd)
	}
}

func TestParseReference_Invalid(t *testing.T) {
	tests := []string{
		"없는책 1",
		"",
		"   ",
		"창세기",
		"창세기 abc",
		"창세기 1:abc",
		"창세기 1:3-abc",
	}

	for _, input := range tests {
		_, err := ParseReference(input)
		if err == nil {
			t.Errorf("expected error for input '%s', got nil", input)
		}
	}
}

func TestParseReference_NumberedBook(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"삼상 1", "1sa"},
		{"사무엘상 1", "1sa"},
		{"1sa 1", "1sa"},
		{"고전 1", "1co"},
		{"고린도전서 1", "1co"},
		{"1co 1", "1co"},
		{"요일 1", "1jn"},
		{"요한1서 1", "1jn"},
		{"1jn 1", "1jn"},
	}

	for _, tt := range tests {
		ref, err := ParseReference(tt.input)
		if err != nil {
			t.Fatalf("unexpected error for '%s': %v", tt.input, err)
		}
		if ref.BookCode != tt.expected {
			t.Errorf("for input '%s', expected BookCode '%s', got '%s'", tt.input, tt.expected, ref.BookCode)
		}
	}
}

func TestParseReference_ExtraWhitespace(t *testing.T) {
	tests := []string{
		"  창세기   1  ",
		"창세기  1",
		"  창 1",
		"창 1  ",
	}

	for _, input := range tests {
		ref, err := ParseReference(input)
		if err != nil {
			t.Fatalf("unexpected error for '%s': %v", input, err)
		}
		if ref.BookCode != "gen" {
			t.Errorf("for input '%s', expected BookCode 'gen', got '%s'", input, ref.BookCode)
		}
		if ref.Chapter != 1 {
			t.Errorf("for input '%s', expected Chapter 1, got %d", input, ref.Chapter)
		}
	}
}

func TestParseReference_ChapterOutOfRange(t *testing.T) {
	_, err := ParseReference("창세기 51")
	if err == nil {
		t.Error("expected error for chapter out of range, got nil")
	}

	_, err = ParseReference("창세기 0")
	if err == nil {
		t.Error("expected error for chapter 0, got nil")
	}
}

func TestParseReference_InvalidVerseRange(t *testing.T) {
	_, err := ParseReference("창세기 1:5-3")
	if err == nil {
		t.Error("expected error for invalid verse range (end < start), got nil")
	}

	_, err = ParseReference("창세기 1:0")
	if err == nil {
		t.Error("expected error for verse 0, got nil")
	}
}

func TestAllBooks(t *testing.T) {
	books := AllBooks()
	if len(books) != 66 {
		t.Errorf("expected 66 books, got %d", len(books))
	}

	oldTestamentCount := 0
	newTestamentCount := 0
	for _, book := range books {
		if book.Testament == "old" {
			oldTestamentCount++
		} else if book.Testament == "new" {
			newTestamentCount++
		}
	}

	if oldTestamentCount != 39 {
		t.Errorf("expected 39 Old Testament books, got %d", oldTestamentCount)
	}
	if newTestamentCount != 27 {
		t.Errorf("expected 27 New Testament books, got %d", newTestamentCount)
	}
}

func TestGetBookName(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{"gen", "창세기"},
		{"exo", "출애굽기"},
		{"mat", "마태복음"},
		{"rev", "요한계시록"},
		{"1sa", "사무엘상"},
		{"1co", "고린도전서"},
		{"invalid", ""},
	}

	for _, tt := range tests {
		result := GetBookName(tt.code)
		if result != tt.expected {
			t.Errorf("GetBookName(%s): expected '%s', got '%s'", tt.code, tt.expected, result)
		}
	}
}

func TestGetBookAbbrev(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{"gen", "창"},
		{"exo", "출"},
		{"mat", "마"},
		{"rev", "계"},
		{"1sa", "삼상"},
		{"1co", "고전"},
		{"invalid", ""},
	}

	for _, tt := range tests {
		result := GetBookAbbrev(tt.code)
		if result != tt.expected {
			t.Errorf("GetBookAbbrev(%s): expected '%s', got '%s'", tt.code, tt.expected, result)
		}
	}
}

func TestGetBookByCode(t *testing.T) {
	book, ok := GetBookByCode("gen")
	if !ok {
		t.Fatal("expected to find book 'gen'")
	}
	if book.Code != "gen" {
		t.Errorf("expected Code 'gen', got '%s'", book.Code)
	}
	if book.NameKo != "창세기" {
		t.Errorf("expected NameKo '창세기', got '%s'", book.NameKo)
	}

	book, ok = GetBookByCode("GEN")
	if !ok {
		t.Fatal("expected to find book 'GEN' (case-insensitive)")
	}
	if book.Code != "gen" {
		t.Errorf("expected Code 'gen', got '%s'", book.Code)
	}

	_, ok = GetBookByCode("invalid")
	if ok {
		t.Error("expected not to find book 'invalid'")
	}
}

func TestGetBookByName(t *testing.T) {
	book, ok := GetBookByName("창세기")
	if !ok {
		t.Fatal("expected to find book '창세기'")
	}
	if book.Code != "gen" {
		t.Errorf("expected Code 'gen', got '%s'", book.Code)
	}

	_, ok = GetBookByName("없는책")
	if ok {
		t.Error("expected not to find book '없는책'")
	}
}

func TestGetBookByAbbrev(t *testing.T) {
	book, ok := GetBookByAbbrev("창")
	if !ok {
		t.Fatal("expected to find book '창'")
	}
	if book.Code != "gen" {
		t.Errorf("expected Code 'gen', got '%s'", book.Code)
	}

	_, ok = GetBookByAbbrev("없")
	if ok {
		t.Error("expected not to find book '없'")
	}
}

func TestAllBooksData(t *testing.T) {
	expectedBooks := map[string]struct {
		nameKo       string
		abbrevKo     string
		testament    string
		chapterCount int
	}{
		"gen": {"창세기", "창", "old", 50},
		"exo": {"출애굽기", "출", "old", 40},
		"mat": {"마태복음", "마", "new", 28},
		"rev": {"요한계시록", "계", "new", 22},
		"1sa": {"사무엘상", "삼상", "old", 31},
		"1co": {"고린도전서", "고전", "new", 16},
		"psa": {"시편", "시", "old", 150},
	}

	for code, expected := range expectedBooks {
		book, ok := GetBookByCode(code)
		if !ok {
			t.Errorf("expected to find book '%s'", code)
			continue
		}
		if book.NameKo != expected.nameKo {
			t.Errorf("book '%s': expected NameKo '%s', got '%s'", code, expected.nameKo, book.NameKo)
		}
		if book.AbbrevKo != expected.abbrevKo {
			t.Errorf("book '%s': expected AbbrevKo '%s', got '%s'", code, expected.abbrevKo, book.AbbrevKo)
		}
		if book.Testament != expected.testament {
			t.Errorf("book '%s': expected Testament '%s', got '%s'", code, expected.testament, book.Testament)
		}
		if book.ChapterCount != expected.chapterCount {
			t.Errorf("book '%s': expected ChapterCount %d, got %d", code, expected.chapterCount, book.ChapterCount)
		}
	}
}
