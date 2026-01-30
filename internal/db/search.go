package db

import (
	"fmt"
	"strings"
)

// SearchResult represents a single search result with context.
type SearchResult struct {
	Verse      Verse
	Snippet    string // text with search term context (surrounding text)
	MatchCount int    // number of matches in this verse
}

// SearchVerses searches for verses matching the query.
// Strategy: Try FTS5 MATCH first. If results < minResults, supplement with LIKE fallback.
// Results are deduplicated (FTS5 and LIKE may return same verses).
func (d *DB) SearchVerses(versionCode, query string, limit int) ([]SearchResult, error) {
	if query == "" {
		return []SearchResult{}, nil
	}

	seen := make(map[int64]bool)
	results := []SearchResult{}

	// Try FTS5 search first
	ftsResults, err := d.searchFTS(versionCode, query, limit, seen)
	if err != nil {
		// FTS5 might fail on special characters, fall back to LIKE only
		likeResults, likeErr := d.searchLIKE(versionCode, query, limit, seen)
		if likeErr != nil {
			return []SearchResult{}, nil
		}
		return likeResults, nil
	}
	results = append(results, ftsResults...)

	// If we have fewer than 5 results, supplement with LIKE
	if len(results) < 5 && len(results) < limit {
		remaining := limit - len(results)
		likeResults, err := d.searchLIKE(versionCode, query, remaining, seen)
		if err != nil {
			return results, nil // return what we have from FTS
		}
		results = append(results, likeResults...)
	}

	return results, nil
}

// searchFTS performs FTS5 full-text search.
func (d *DB) searchFTS(versionCode, query string, limit int, seen map[int64]bool) ([]SearchResult, error) {
	rows, err := d.conn.Query(
		`SELECT v.id, v.book_id, v.chapter, v.verse_num, v.text,
		        COALESCE(v.section_title, ''), v.has_footnote,
		        b.name_ko, b.code,
		        snippet(verses_fts, 0, '**', '**', '...', 32) as snippet
		 FROM verses_fts
		 JOIN verses v ON v.id = verses_fts.rowid
		 JOIN books b ON b.id = v.book_id
		 JOIN versions ver ON ver.id = b.version_id
		 WHERE ver.code = ? AND verses_fts MATCH ?
		 ORDER BY rank
		 LIMIT ?`,
		versionCode, query, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []SearchResult{}
	for rows.Next() {
		var v Verse
		var snippet string
		if err := rows.Scan(
			&v.ID, &v.BookID, &v.Chapter, &v.VerseNum, &v.Text,
			&v.SectionTitle, &v.HasFootnote,
			&v.BookName, &v.BookCode,
			&snippet,
		); err != nil {
			return nil, fmt.Errorf("scan fts result: %w", err)
		}

		if seen[v.ID] {
			continue
		}
		seen[v.ID] = true

		matchCount := countMatches(v.Text, query)
		results = append(results, SearchResult{
			Verse:      v,
			Snippet:    snippet,
			MatchCount: matchCount,
		})
	}

	return results, rows.Err()
}

// searchLIKE performs LIKE-based search as fallback.
func (d *DB) searchLIKE(versionCode, query string, limit int, seen map[int64]bool) ([]SearchResult, error) {
	rows, err := d.conn.Query(
		`SELECT v.id, v.book_id, v.chapter, v.verse_num, v.text,
		        COALESCE(v.section_title, ''), v.has_footnote,
		        b.name_ko, b.code
		 FROM verses v
		 JOIN books b ON b.id = v.book_id
		 JOIN versions ver ON ver.id = b.version_id
		 WHERE ver.code = ? AND v.text LIKE '%' || ? || '%'
		 ORDER BY b.sort_order, v.chapter, v.verse_num
		 LIMIT ?`,
		versionCode, query, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("like search: %w", err)
	}
	defer rows.Close()

	results := []SearchResult{}
	for rows.Next() {
		var v Verse
		if err := rows.Scan(
			&v.ID, &v.BookID, &v.Chapter, &v.VerseNum, &v.Text,
			&v.SectionTitle, &v.HasFootnote,
			&v.BookName, &v.BookCode,
		); err != nil {
			return nil, fmt.Errorf("scan like result: %w", err)
		}

		if seen[v.ID] {
			continue
		}
		seen[v.ID] = true

		snippet := generateSnippet(v.Text, query, 30)
		matchCount := countMatches(v.Text, query)
		results = append(results, SearchResult{
			Verse:      v,
			Snippet:    snippet,
			MatchCount: matchCount,
		})
	}

	return results, rows.Err()
}

// generateSnippet creates a text snippet showing context around the search term.
func generateSnippet(text, query string, contextChars int) string {
	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(query)

	pos := strings.Index(lowerText, lowerQuery)
	if pos == -1 {
		// No match found, return beginning of text
		if len(text) > contextChars*2 {
			return text[:contextChars*2] + "..."
		}
		return text
	}

	// Calculate start and end positions
	start := pos - contextChars
	if start < 0 {
		start = 0
	}

	end := pos + len(query) + contextChars
	if end > len(text) {
		end = len(text)
	}

	// Extract snippet
	snippet := text[start:end]

	// Add ellipsis if truncated
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(text) {
		snippet = snippet + "..."
	}

	return snippet
}

// countMatches counts how many times query appears in text (case-insensitive).
func countMatches(text, query string) int {
	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(query)

	count := 0
	pos := 0
	for {
		idx := strings.Index(lowerText[pos:], lowerQuery)
		if idx == -1 {
			break
		}
		count++
		pos += idx + len(lowerQuery)
	}
	return count
}
