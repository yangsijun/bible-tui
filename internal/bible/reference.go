package bible

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Reference struct {
	BookCode   string
	Chapter    int
	VerseStart int
	VerseEnd   int
}

var referencePattern = regexp.MustCompile(`^(.+?)\s+(\d+)(?::(\d+)(?:-(\d+))?)?$`)

func ParseReference(input string) (*Reference, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return nil, fmt.Errorf("empty reference")
	}

	matches := referencePattern.FindStringSubmatch(trimmed)
	if matches == nil {
		return nil, fmt.Errorf("invalid reference format: %s", input)
	}

	bookPart := strings.TrimSpace(matches[1])
	chapterStr := matches[2]
	verseStartStr := matches[3]
	verseEndStr := matches[4]

	book := findBook(bookPart)
	if book == nil {
		return nil, fmt.Errorf("unknown book: %s", bookPart)
	}

	chapter, err := strconv.Atoi(chapterStr)
	if err != nil {
		return nil, fmt.Errorf("invalid chapter number: %s", chapterStr)
	}

	if chapter < 1 || chapter > book.ChapterCount {
		return nil, fmt.Errorf("chapter %d out of range for %s (max: %d)", chapter, book.NameKo, book.ChapterCount)
	}

	ref := &Reference{
		BookCode:   book.Code,
		Chapter:    chapter,
		VerseStart: 0,
		VerseEnd:   0,
	}

	if verseStartStr != "" {
		verseStart, err := strconv.Atoi(verseStartStr)
		if err != nil {
			return nil, fmt.Errorf("invalid verse number: %s", verseStartStr)
		}
		if verseStart < 1 {
			return nil, fmt.Errorf("verse number must be positive: %d", verseStart)
		}
		ref.VerseStart = verseStart

		if verseEndStr != "" {
			verseEnd, err := strconv.Atoi(verseEndStr)
			if err != nil {
				return nil, fmt.Errorf("invalid verse number: %s", verseEndStr)
			}
			if verseEnd < verseStart {
				return nil, fmt.Errorf("end verse %d cannot be less than start verse %d", verseEnd, verseStart)
			}
			ref.VerseEnd = verseEnd
		}
	}

	return ref, nil
}

func findBook(name string) *BookInfo {
	if book, ok := GetBookByName(name); ok {
		return book
	}

	if book, ok := GetBookByAbbrev(name); ok {
		return book
	}

	if book, ok := GetBookByCode(strings.ToLower(name)); ok {
		return book
	}

	return nil
}

func GetBookName(code string) string {
	book, ok := GetBookByCode(code)
	if !ok {
		return ""
	}
	return book.NameKo
}

func GetBookAbbrev(code string) string {
	book, ok := GetBookByCode(code)
	if !ok {
		return ""
	}
	return book.AbbrevKo
}
