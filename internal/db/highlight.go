package db

import (
	"database/sql"
	"fmt"
)

type HighlightWithVerse struct {
	Highlight
	VerseText string
	BookName  string
	BookCode  string
	Chapter   int
	VerseNum  int
}

// AddHighlight adds or updates a highlight for a verse.
// Uses INSERT OR REPLACE since verse_id is UNIQUE.
func (d *DB) AddHighlight(verseID int64, color string) error {
	_, err := d.conn.Exec(
		"INSERT OR REPLACE INTO highlights (verse_id, color) VALUES (?, ?)",
		verseID, color,
	)
	if err != nil {
		return fmt.Errorf("add highlight: %w", err)
	}
	return nil
}

// RemoveHighlight removes a highlight by verse ID.
func (d *DB) RemoveHighlight(verseID int64) error {
	_, err := d.conn.Exec("DELETE FROM highlights WHERE verse_id = ?", verseID)
	if err != nil {
		return fmt.Errorf("remove highlight: %w", err)
	}
	return nil
}

// ListHighlights returns highlights ordered by newest first.
func (d *DB) ListHighlights(limit, offset int) ([]HighlightWithVerse, error) {
	rows, err := d.conn.Query(
		`SELECT h.id, h.verse_id, h.color, h.created_at,
		        v.text, b.name_ko, b.code, v.chapter, v.verse_num
		 FROM highlights h
		 JOIN verses v ON v.id = h.verse_id
		 JOIN books b ON b.id = v.book_id
		 ORDER BY h.created_at DESC
		 LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list highlights: %w", err)
	}
	defer rows.Close()

	var highlights []HighlightWithVerse
	for rows.Next() {
		var h HighlightWithVerse
		if err := rows.Scan(
			&h.ID, &h.VerseID, &h.Color, &h.CreatedAt,
			&h.VerseText, &h.BookName, &h.BookCode, &h.Chapter, &h.VerseNum,
		); err != nil {
			return nil, fmt.Errorf("scan highlight: %w", err)
		}
		highlights = append(highlights, h)
	}
	return highlights, rows.Err()
}

// GetHighlightColor returns the highlight color for a verse, or "" if not highlighted.
func (d *DB) GetHighlightColor(verseID int64) (string, error) {
	var color string
	err := d.conn.QueryRow(
		"SELECT color FROM highlights WHERE verse_id = ?",
		verseID,
	).Scan(&color)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get highlight color: %w", err)
	}
	return color, nil
}
