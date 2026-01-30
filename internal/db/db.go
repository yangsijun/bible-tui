package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

func Open(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if _, err := conn.Exec("PRAGMA journal_mode=WAL"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("set WAL mode: %w", err)
	}
	if _, err := conn.Exec("PRAGMA foreign_keys=ON"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}
	return &DB{conn: conn}, nil
}

func OpenMemory() (*DB, error) {
	return Open(":memory:")
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) Migrate() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS versions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			lang TEXT NOT NULL DEFAULT 'ko'
		)`,
		`CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version_id INTEGER NOT NULL REFERENCES versions(id),
			code TEXT NOT NULL,
			name_ko TEXT NOT NULL,
			abbrev_ko TEXT NOT NULL,
			testament TEXT NOT NULL,
			chapter_count INTEGER NOT NULL,
			sort_order INTEGER NOT NULL,
			UNIQUE(version_id, code)
		)`,
		`CREATE TABLE IF NOT EXISTS verses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			book_id INTEGER NOT NULL REFERENCES books(id),
			chapter INTEGER NOT NULL,
			verse_num INTEGER NOT NULL,
			text TEXT NOT NULL,
			section_title TEXT,
			has_footnote BOOLEAN NOT NULL DEFAULT 0,
			UNIQUE(book_id, chapter, verse_num)
		)`,
		`CREATE TABLE IF NOT EXISTS footnotes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			verse_id INTEGER NOT NULL REFERENCES verses(id),
			marker TEXT,
			content TEXT NOT NULL
		)`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS verses_fts USING fts5(
			text,
			content=verses,
			content_rowid=id,
			tokenize='unicode61'
		)`,
		`CREATE TRIGGER IF NOT EXISTS verses_ai AFTER INSERT ON verses BEGIN
			INSERT INTO verses_fts(rowid, text) VALUES (new.id, new.text);
		END`,
		`CREATE TRIGGER IF NOT EXISTS verses_ad AFTER DELETE ON verses BEGIN
			INSERT INTO verses_fts(verses_fts, rowid, text) VALUES('delete', old.id, old.text);
		END`,
		`CREATE TRIGGER IF NOT EXISTS verses_au AFTER UPDATE ON verses BEGIN
			INSERT INTO verses_fts(verses_fts, rowid, text) VALUES('delete', old.id, old.text);
			INSERT INTO verses_fts(rowid, text) VALUES (new.id, new.text);
		END`,
		`CREATE TABLE IF NOT EXISTS bookmarks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			verse_id INTEGER NOT NULL REFERENCES verses(id),
			note TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS highlights (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			verse_id INTEGER NOT NULL REFERENCES verses(id),
			color TEXT NOT NULL DEFAULT 'yellow',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(verse_id)
		)`,
		`CREATE TABLE IF NOT EXISTS reading_plans (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			plan_type TEXT NOT NULL,
			version_id INTEGER NOT NULL REFERENCES versions(id),
			total_days INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS reading_plan_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			plan_id INTEGER NOT NULL REFERENCES reading_plans(id),
			day_number INTEGER NOT NULL,
			book_code TEXT NOT NULL,
			chapter_start INTEGER NOT NULL,
			chapter_end INTEGER NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT 0,
			completed_at DATETIME,
			UNIQUE(plan_id, day_number, book_code, chapter_start)
		)`,
		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS crawl_status (
			version_code TEXT NOT NULL,
			book_code TEXT NOT NULL,
			chapter INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			verse_count INTEGER,
			crawled_at DATETIME,
			error_msg TEXT,
			PRIMARY KEY(version_code, book_code, chapter)
		)`,
	}

	for _, stmt := range statements {
		if _, err := d.conn.Exec(stmt); err != nil {
			return fmt.Errorf("migrate: %w\nSQL: %s", err, stmt)
		}
	}
	return nil
}

func (d *DB) InsertVersion(code, name, lang string) (int64, error) {
	res, err := d.conn.Exec(
		"INSERT OR IGNORE INTO versions (code, name, lang) VALUES (?, ?, ?)",
		code, name, lang,
	)
	if err != nil {
		return 0, fmt.Errorf("insert version: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("insert version last id: %w", err)
	}
	if id == 0 {
		return d.getVersionID(code)
	}
	return id, nil
}

func (d *DB) getVersionID(code string) (int64, error) {
	var id int64
	err := d.conn.QueryRow("SELECT id FROM versions WHERE code = ?", code).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("get version id: %w", err)
	}
	return id, nil
}

func (d *DB) GetVersionByCode(code string) (*Version, error) {
	v := &Version{}
	err := d.conn.QueryRow(
		"SELECT id, code, name, lang FROM versions WHERE code = ?", code,
	).Scan(&v.ID, &v.Code, &v.Name, &v.Lang)
	if err != nil {
		return nil, fmt.Errorf("get version by code: %w", err)
	}
	return v, nil
}

func (d *DB) InsertBook(versionID int64, code, nameKo, abbrevKo, testament string, chapterCount, sortOrder int) (int64, error) {
	res, err := d.conn.Exec(
		`INSERT OR IGNORE INTO books (version_id, code, name_ko, abbrev_ko, testament, chapter_count, sort_order)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		versionID, code, nameKo, abbrevKo, testament, chapterCount, sortOrder,
	)
	if err != nil {
		return 0, fmt.Errorf("insert book: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("insert book last id: %w", err)
	}
	return id, nil
}

func (d *DB) GetBookByCode(versionCode, bookCode string) (*Book, error) {
	b := &Book{}
	err := d.conn.QueryRow(
		`SELECT b.id, b.version_id, b.code, b.name_ko, b.abbrev_ko, b.testament, b.chapter_count, b.sort_order
		 FROM books b
		 JOIN versions v ON v.id = b.version_id
		 WHERE v.code = ? AND b.code = ?`,
		versionCode, bookCode,
	).Scan(&b.ID, &b.VersionID, &b.Code, &b.NameKo, &b.AbbrevKo, &b.Testament, &b.ChapterCount, &b.SortOrder)
	if err != nil {
		return nil, fmt.Errorf("get book by code: %w", err)
	}
	return b, nil
}

func (d *DB) InsertVerse(bookID int64, chapter, verseNum int, text, sectionTitle string, hasFootnote bool) (int64, error) {
	var secTitle interface{}
	if sectionTitle != "" {
		secTitle = sectionTitle
	}
	res, err := d.conn.Exec(
		`INSERT INTO verses (book_id, chapter, verse_num, text, section_title, has_footnote)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		bookID, chapter, verseNum, text, secTitle, hasFootnote,
	)
	if err != nil {
		return 0, fmt.Errorf("insert verse: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("insert verse last id: %w", err)
	}
	return id, nil
}

func (d *DB) InsertFootnote(verseID int64, marker, content string) error {
	_, err := d.conn.Exec(
		"INSERT INTO footnotes (verse_id, marker, content) VALUES (?, ?, ?)",
		verseID, marker, content,
	)
	if err != nil {
		return fmt.Errorf("insert footnote: %w", err)
	}
	return nil
}

func (d *DB) GetVerses(versionCode, bookCode string, chapter int) ([]Verse, error) {
	rows, err := d.conn.Query(
		`SELECT v.id, v.book_id, v.chapter, v.verse_num, v.text,
		        COALESCE(v.section_title, ''), v.has_footnote,
		        b.name_ko, b.code
		 FROM verses v
		 JOIN books b ON b.id = v.book_id
		 JOIN versions ver ON ver.id = b.version_id
		 WHERE ver.code = ? AND b.code = ? AND v.chapter = ?
		 ORDER BY v.verse_num`,
		versionCode, bookCode, chapter,
	)
	if err != nil {
		return nil, fmt.Errorf("get verses: %w", err)
	}
	defer rows.Close()

	var verses []Verse
	for rows.Next() {
		var v Verse
		if err := rows.Scan(
			&v.ID, &v.BookID, &v.Chapter, &v.VerseNum, &v.Text,
			&v.SectionTitle, &v.HasFootnote,
			&v.BookName, &v.BookCode,
		); err != nil {
			return nil, fmt.Errorf("scan verse: %w", err)
		}
		verses = append(verses, v)
	}
	return verses, rows.Err()
}

func (d *DB) GetRandomVerse(versionCode string) (*Verse, error) {
	v := &Verse{}
	err := d.conn.QueryRow(
		`SELECT v.id, v.book_id, v.chapter, v.verse_num, v.text,
		        COALESCE(v.section_title, ''), v.has_footnote,
		        b.name_ko, b.code
		 FROM verses v
		 JOIN books b ON b.id = v.book_id
		 JOIN versions ver ON ver.id = b.version_id
		 WHERE ver.code = ?
		 ORDER BY RANDOM()
		 LIMIT 1`,
		versionCode,
	).Scan(
		&v.ID, &v.BookID, &v.Chapter, &v.VerseNum, &v.Text,
		&v.SectionTitle, &v.HasFootnote,
		&v.BookName, &v.BookCode,
	)
	if err != nil {
		return nil, fmt.Errorf("get random verse: %w", err)
	}
	return v, nil
}

func (d *DB) GetSetting(key string) (string, error) {
	var value string
	err := d.conn.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get setting: %w", err)
	}
	return value, nil
}

func (d *DB) SetSetting(key, value string) error {
	_, err := d.conn.Exec(
		"INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)",
		key, value,
	)
	if err != nil {
		return fmt.Errorf("set setting: %w", err)
	}
	return nil
}
