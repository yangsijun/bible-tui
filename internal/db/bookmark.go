package db

import (
	"fmt"
)

type BookmarkWithVerse struct {
	Bookmark
	VerseText string
	BookName  string
	BookCode  string
	Chapter   int
	VerseNum  int
}

// AddBookmark adds a bookmark for a verse with an optional note.
// Returns the bookmark ID.
func (d *DB) AddBookmark(verseID int64, note string) (int64, error) {
	var notePtr interface{}
	if note != "" {
		notePtr = note
	}
	res, err := d.conn.Exec(
		"INSERT INTO bookmarks (verse_id, note) VALUES (?, ?)",
		verseID, notePtr,
	)
	if err != nil {
		return 0, fmt.Errorf("add bookmark: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("add bookmark last id: %w", err)
	}
	return id, nil
}

// RemoveBookmark removes a bookmark by its ID.
func (d *DB) RemoveBookmark(id int64) error {
	_, err := d.conn.Exec("DELETE FROM bookmarks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("remove bookmark: %w", err)
	}
	return nil
}

// ListBookmarks returns bookmarks ordered by newest first.
// Includes verse text and book info via joins.
func (d *DB) ListBookmarks(limit, offset int) ([]BookmarkWithVerse, error) {
	rows, err := d.conn.Query(
		`SELECT bm.id, bm.verse_id, COALESCE(bm.note, ''), bm.created_at,
		        v.text, b.name_ko, b.code, v.chapter, v.verse_num
		 FROM bookmarks bm
		 JOIN verses v ON v.id = bm.verse_id
		 JOIN books b ON b.id = v.book_id
		 ORDER BY bm.created_at DESC
		 LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list bookmarks: %w", err)
	}
	defer rows.Close()

	var bookmarks []BookmarkWithVerse
	for rows.Next() {
		var bm BookmarkWithVerse
		if err := rows.Scan(
			&bm.ID, &bm.VerseID, &bm.Note, &bm.CreatedAt,
			&bm.VerseText, &bm.BookName, &bm.BookCode, &bm.Chapter, &bm.VerseNum,
		); err != nil {
			return nil, fmt.Errorf("scan bookmark: %w", err)
		}
		bookmarks = append(bookmarks, bm)
	}
	return bookmarks, rows.Err()
}

// IsBookmarked checks if a verse is bookmarked.
func (d *DB) IsBookmarked(verseID int64) (bool, error) {
	var count int
	err := d.conn.QueryRow(
		"SELECT COUNT(*) FROM bookmarks WHERE verse_id = ?",
		verseID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("is bookmarked: %w", err)
	}
	return count > 0, nil
}
