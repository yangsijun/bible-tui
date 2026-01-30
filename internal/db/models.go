package db

import "time"

type Version struct {
	ID   int64
	Code string // e.g. "GAE", "HAN"
	Name string // e.g. "개역개정", "개역한글"
	Lang string // e.g. "ko"
}

type Book struct {
	ID           int64
	VersionID    int64
	Code         string // e.g. "gen", "exo"
	NameKo       string // e.g. "창세기"
	AbbrevKo     string // e.g. "창"
	Testament    string // "old" or "new"
	ChapterCount int
	SortOrder    int
}

type Verse struct {
	ID           int64
	BookID       int64
	Chapter      int
	VerseNum     int
	Text         string
	SectionTitle string // nullable
	HasFootnote  bool

	// joined fields populated by queries
	BookName string
	BookCode string
}

type Footnote struct {
	ID      int64
	VerseID int64
	Marker  string // e.g. "1)", "2)"
	Content string
}

type Bookmark struct {
	ID        int64
	VerseID   int64
	Note      string
	CreatedAt time.Time
}

type Highlight struct {
	ID        int64
	VerseID   int64
	Color     string
	CreatedAt time.Time
}
