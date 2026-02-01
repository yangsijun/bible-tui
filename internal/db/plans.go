package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yangsijun/bible-tui/internal/bible"
)

type ReadingPlan struct {
	ID        int64
	Name      string
	PlanType  string // "sequential", "mcheyne", "custom"
	VersionID int64
	TotalDays int
	CreatedAt time.Time
}

type PlanEntry struct {
	ID           int64
	PlanID       int64
	DayNumber    int
	BookCode     string
	ChapterStart int
	ChapterEnd   int
	Completed    bool
	CompletedAt  *time.Time
}

func (d *DB) CreateSequentialPlan(versionID int64, name string) (int64, error) {
	books := bible.AllBooks()

	type chapRef struct {
		bookCode string
		chapter  int
	}
	var chapters []chapRef
	for _, b := range books {
		for ch := 1; ch <= b.ChapterCount; ch++ {
			chapters = append(chapters, chapRef{b.Code, ch})
		}
	}

	chaptersPerDay := 3
	totalDays := (len(chapters) + chaptersPerDay - 1) / chaptersPerDay

	res, err := d.conn.Exec(
		"INSERT INTO reading_plans (name, plan_type, version_id, total_days) VALUES (?, 'sequential', ?, ?)",
		name, versionID, totalDays,
	)
	if err != nil {
		return 0, fmt.Errorf("create sequential plan: %w", err)
	}
	planID, _ := res.LastInsertId()

	tx, err := d.conn.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		"INSERT INTO reading_plan_entries (plan_id, day_number, book_code, chapter_start, chapter_end) VALUES (?, ?, ?, ?, ?)",
	)
	if err != nil {
		return 0, fmt.Errorf("prepare stmt: %w", err)
	}
	defer stmt.Close()

	for day := 0; day < totalDays; day++ {
		start := day * chaptersPerDay
		end := start + chaptersPerDay
		if end > len(chapters) {
			end = len(chapters)
		}

		var curBook string
		var curStart, curEnd int
		for i := start; i < end; i++ {
			ref := chapters[i]
			if ref.bookCode == curBook && ref.chapter == curEnd+1 {
				curEnd = ref.chapter
			} else {
				if curBook != "" {
					if _, err := stmt.Exec(planID, day+1, curBook, curStart, curEnd); err != nil {
						return 0, fmt.Errorf("insert entry: %w", err)
					}
				}
				curBook = ref.bookCode
				curStart = ref.chapter
				curEnd = ref.chapter
			}
		}
		if curBook != "" {
			if _, err := stmt.Exec(planID, day+1, curBook, curStart, curEnd); err != nil {
				return 0, fmt.Errorf("insert entry: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return planID, nil
}

func (d *DB) CreateMcCheynePlan(versionID int64, name string) (int64, error) {
	schedule := bible.McCheyneSchedule()

	res, err := d.conn.Exec(
		"INSERT INTO reading_plans (name, plan_type, version_id, total_days) VALUES (?, 'mcheyne', ?, ?)",
		name, versionID, 365,
	)
	if err != nil {
		return 0, fmt.Errorf("create mcheyne plan: %w", err)
	}
	planID, _ := res.LastInsertId()

	tx, err := d.conn.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		"INSERT INTO reading_plan_entries (plan_id, day_number, book_code, chapter_start, chapter_end) VALUES (?, ?, ?, ?, ?)",
	)
	if err != nil {
		return 0, fmt.Errorf("prepare stmt: %w", err)
	}
	defer stmt.Close()

	for day, entries := range schedule {
		for _, e := range entries {
			if _, err := stmt.Exec(planID, day+1, e.BookCode, e.ChapterStart, e.ChapterEnd); err != nil {
				return 0, fmt.Errorf("insert mcheyne entry: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return planID, nil
}

func (d *DB) CreateCustomPlan(versionID int64, name string, entries []PlanEntry) (int64, error) {
	maxDay := 0
	for _, e := range entries {
		if e.DayNumber > maxDay {
			maxDay = e.DayNumber
		}
	}

	res, err := d.conn.Exec(
		"INSERT INTO reading_plans (name, plan_type, version_id, total_days) VALUES (?, 'custom', ?, ?)",
		name, versionID, maxDay,
	)
	if err != nil {
		return 0, fmt.Errorf("create custom plan: %w", err)
	}
	planID, _ := res.LastInsertId()

	tx, err := d.conn.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		"INSERT INTO reading_plan_entries (plan_id, day_number, book_code, chapter_start, chapter_end) VALUES (?, ?, ?, ?, ?)",
	)
	if err != nil {
		return 0, fmt.Errorf("prepare stmt: %w", err)
	}
	defer stmt.Close()

	for _, e := range entries {
		if _, err := stmt.Exec(planID, e.DayNumber, e.BookCode, e.ChapterStart, e.ChapterEnd); err != nil {
			return 0, fmt.Errorf("insert custom entry: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return planID, nil
}

func (d *DB) GetActivePlan(versionID int64) (*ReadingPlan, error) {
	p := &ReadingPlan{}
	err := d.conn.QueryRow(
		"SELECT id, name, plan_type, version_id, total_days, created_at FROM reading_plans WHERE version_id = ? ORDER BY created_at DESC LIMIT 1",
		versionID,
	).Scan(&p.ID, &p.Name, &p.PlanType, &p.VersionID, &p.TotalDays, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get active plan: %w", err)
	}
	return p, nil
}

func (d *DB) GetTodayEntries(planID int64) ([]PlanEntry, error) {
	var createdAt time.Time
	var totalDays int
	err := d.conn.QueryRow(
		"SELECT created_at, total_days FROM reading_plans WHERE id = ?", planID,
	).Scan(&createdAt, &totalDays)
	if err != nil {
		return nil, fmt.Errorf("get plan for today: %w", err)
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	created := time.Date(createdAt.Year(), createdAt.Month(), createdAt.Day(), 0, 0, 0, 0, createdAt.Location())
	dayNumber := int(today.Sub(created).Hours()/24) + 1
	if dayNumber < 1 {
		dayNumber = 1
	}
	if dayNumber > totalDays {
		dayNumber = totalDays
	}

	rows, err := d.conn.Query(
		`SELECT id, plan_id, day_number, book_code, chapter_start, chapter_end, completed, completed_at
		 FROM reading_plan_entries WHERE plan_id = ? AND day_number = ?
		 ORDER BY id`,
		planID, dayNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("get today entries: %w", err)
	}
	defer rows.Close()

	var entries []PlanEntry
	for rows.Next() {
		var e PlanEntry
		if err := rows.Scan(&e.ID, &e.PlanID, &e.DayNumber, &e.BookCode,
			&e.ChapterStart, &e.ChapterEnd, &e.Completed, &e.CompletedAt); err != nil {
			return nil, fmt.Errorf("scan plan entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (d *DB) MarkEntryCompleted(entryID int64) error {
	_, err := d.conn.Exec(
		"UPDATE reading_plan_entries SET completed = 1, completed_at = CURRENT_TIMESTAMP WHERE id = ?",
		entryID,
	)
	if err != nil {
		return fmt.Errorf("mark entry completed: %w", err)
	}
	return nil
}

func (d *DB) GetPlanProgress(planID int64) (completed, total int, err error) {
	err = d.conn.QueryRow(
		"SELECT COUNT(*), COALESCE(SUM(CASE WHEN completed = 1 THEN 1 ELSE 0 END), 0) FROM reading_plan_entries WHERE plan_id = ?",
		planID,
	).Scan(&total, &completed)
	if err != nil {
		return 0, 0, fmt.Errorf("get plan progress: %w", err)
	}
	return completed, total, nil
}

func (d *DB) ListPlans() ([]ReadingPlan, error) {
	rows, err := d.conn.Query(
		"SELECT id, name, plan_type, version_id, total_days, created_at FROM reading_plans ORDER BY created_at DESC, id DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("list plans: %w", err)
	}
	defer rows.Close()

	var plans []ReadingPlan
	for rows.Next() {
		var p ReadingPlan
		if err := rows.Scan(&p.ID, &p.Name, &p.PlanType, &p.VersionID, &p.TotalDays, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan plan: %w", err)
		}
		plans = append(plans, p)
	}
	return plans, rows.Err()
}

func (d *DB) DeletePlan(planID int64) error {
	_, err := d.conn.Exec("DELETE FROM reading_plan_entries WHERE plan_id = ?", planID)
	if err != nil {
		return fmt.Errorf("delete plan entries: %w", err)
	}
	_, err = d.conn.Exec("DELETE FROM reading_plans WHERE id = ?", planID)
	if err != nil {
		return fmt.Errorf("delete plan: %w", err)
	}
	return nil
}
