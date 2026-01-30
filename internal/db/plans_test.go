package db

import (
	"testing"
)

func setupPlanDB(t *testing.T) (*DB, int64) {
	t.Helper()
	d, err := OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	if err := d.Migrate(); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	t.Cleanup(func() { d.Close() })

	versionID, err := d.InsertVersion("GAE", "개역개정", "ko")
	if err != nil {
		t.Fatalf("InsertVersion: %v", err)
	}
	return d, versionID
}

func TestPlanCreateSequential(t *testing.T) {
	d, vID := setupPlanDB(t)

	planID, err := d.CreateSequentialPlan(vID, "통독")
	if err != nil {
		t.Fatalf("CreateSequentialPlan: %v", err)
	}
	if planID == 0 {
		t.Fatal("expected non-zero plan ID")
	}

	plan, err := d.GetActivePlan(vID)
	if err != nil {
		t.Fatalf("GetActivePlan: %v", err)
	}
	if plan == nil {
		t.Fatal("expected non-nil plan")
	}
	if plan.TotalDays < 300 {
		t.Errorf("expected > 300 days, got %d", plan.TotalDays)
	}
	if plan.PlanType != "sequential" {
		t.Errorf("expected plan_type 'sequential', got %q", plan.PlanType)
	}

	var firstBookCode string
	var firstChStart int
	err = d.conn.QueryRow(
		"SELECT book_code, chapter_start FROM reading_plan_entries WHERE plan_id = ? ORDER BY day_number, id LIMIT 1",
		planID,
	).Scan(&firstBookCode, &firstChStart)
	if err != nil {
		t.Fatalf("query first entry: %v", err)
	}
	if firstBookCode != "gen" {
		t.Errorf("first entry book: expected 'gen', got %q", firstBookCode)
	}
	if firstChStart != 1 {
		t.Errorf("first entry chapter_start: expected 1, got %d", firstChStart)
	}
}

func TestPlanCreateMcCheyne(t *testing.T) {
	d, vID := setupPlanDB(t)

	planID, err := d.CreateMcCheynePlan(vID, "맥체인")
	if err != nil {
		t.Fatalf("CreateMcCheynePlan: %v", err)
	}
	if planID == 0 {
		t.Fatal("expected non-zero plan ID")
	}

	plan, err := d.GetActivePlan(vID)
	if err != nil {
		t.Fatalf("GetActivePlan: %v", err)
	}
	if plan == nil {
		t.Fatal("expected non-nil plan")
	}
	if plan.TotalDays != 365 {
		t.Errorf("expected 365 days, got %d", plan.TotalDays)
	}
	if plan.PlanType != "mcheyne" {
		t.Errorf("expected plan_type 'mcheyne', got %q", plan.PlanType)
	}

	var distinctDays int
	err = d.conn.QueryRow(
		"SELECT COUNT(DISTINCT day_number) FROM reading_plan_entries WHERE plan_id = ?", planID,
	).Scan(&distinctDays)
	if err != nil {
		t.Fatalf("count distinct days: %v", err)
	}
	if distinctDays != 365 {
		t.Errorf("expected 365 distinct days, got %d", distinctDays)
	}
}

func TestPlanMarkCompleted(t *testing.T) {
	d, vID := setupPlanDB(t)

	entries := []PlanEntry{
		{DayNumber: 1, BookCode: "gen", ChapterStart: 1, ChapterEnd: 3},
		{DayNumber: 1, BookCode: "psa", ChapterStart: 1, ChapterEnd: 1},
	}
	planID, err := d.CreateCustomPlan(vID, "custom", entries)
	if err != nil {
		t.Fatalf("CreateCustomPlan: %v", err)
	}

	todayEntries, err := d.GetTodayEntries(planID)
	if err != nil {
		t.Fatalf("GetTodayEntries: %v", err)
	}
	if len(todayEntries) == 0 {
		t.Fatal("expected entries for today")
	}

	entryID := todayEntries[0].ID
	if err := d.MarkEntryCompleted(entryID); err != nil {
		t.Fatalf("MarkEntryCompleted: %v", err)
	}

	var completed bool
	err = d.conn.QueryRow(
		"SELECT completed FROM reading_plan_entries WHERE id = ?", entryID,
	).Scan(&completed)
	if err != nil {
		t.Fatalf("query completed: %v", err)
	}
	if !completed {
		t.Error("expected completed = true")
	}
}

func TestPlanGetProgress(t *testing.T) {
	d, vID := setupPlanDB(t)

	entries := []PlanEntry{
		{DayNumber: 1, BookCode: "gen", ChapterStart: 1, ChapterEnd: 3},
		{DayNumber: 1, BookCode: "psa", ChapterStart: 1, ChapterEnd: 1},
		{DayNumber: 2, BookCode: "gen", ChapterStart: 4, ChapterEnd: 6},
	}
	planID, err := d.CreateCustomPlan(vID, "progress-test", entries)
	if err != nil {
		t.Fatalf("CreateCustomPlan: %v", err)
	}

	completed, total, err := d.GetPlanProgress(planID)
	if err != nil {
		t.Fatalf("GetPlanProgress: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if completed != 0 {
		t.Errorf("expected completed 0, got %d", completed)
	}

	todayEntries, err := d.GetTodayEntries(planID)
	if err != nil {
		t.Fatalf("GetTodayEntries: %v", err)
	}
	if len(todayEntries) > 0 {
		if err := d.MarkEntryCompleted(todayEntries[0].ID); err != nil {
			t.Fatalf("MarkEntryCompleted: %v", err)
		}
	}

	completed, total, err = d.GetPlanProgress(planID)
	if err != nil {
		t.Fatalf("GetPlanProgress after mark: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if completed != 1 {
		t.Errorf("expected completed 1, got %d", completed)
	}
}

func TestPlanGetTodayEntries(t *testing.T) {
	d, vID := setupPlanDB(t)

	entries := []PlanEntry{
		{DayNumber: 1, BookCode: "gen", ChapterStart: 1, ChapterEnd: 3},
		{DayNumber: 1, BookCode: "mat", ChapterStart: 1, ChapterEnd: 1},
		{DayNumber: 2, BookCode: "gen", ChapterStart: 4, ChapterEnd: 6},
	}
	planID, err := d.CreateCustomPlan(vID, "today-test", entries)
	if err != nil {
		t.Fatalf("CreateCustomPlan: %v", err)
	}

	todayEntries, err := d.GetTodayEntries(planID)
	if err != nil {
		t.Fatalf("GetTodayEntries: %v", err)
	}
	if len(todayEntries) != 2 {
		t.Fatalf("expected 2 entries for day 1, got %d", len(todayEntries))
	}
	if todayEntries[0].BookCode != "gen" {
		t.Errorf("expected first entry 'gen', got %q", todayEntries[0].BookCode)
	}
	if todayEntries[1].BookCode != "mat" {
		t.Errorf("expected second entry 'mat', got %q", todayEntries[1].BookCode)
	}
}

func TestPlanListPlans(t *testing.T) {
	d, vID := setupPlanDB(t)

	_, err := d.CreateCustomPlan(vID, "plan-A", []PlanEntry{
		{DayNumber: 1, BookCode: "gen", ChapterStart: 1, ChapterEnd: 1},
	})
	if err != nil {
		t.Fatalf("CreateCustomPlan A: %v", err)
	}
	_, err = d.CreateCustomPlan(vID, "plan-B", []PlanEntry{
		{DayNumber: 1, BookCode: "exo", ChapterStart: 1, ChapterEnd: 1},
	})
	if err != nil {
		t.Fatalf("CreateCustomPlan B: %v", err)
	}

	plans, err := d.ListPlans()
	if err != nil {
		t.Fatalf("ListPlans: %v", err)
	}
	if len(plans) != 2 {
		t.Fatalf("expected 2 plans, got %d", len(plans))
	}
	if plans[0].Name != "plan-B" {
		t.Errorf("expected first plan 'plan-B' (newest), got %q", plans[0].Name)
	}
	if plans[1].Name != "plan-A" {
		t.Errorf("expected second plan 'plan-A' (oldest), got %q", plans[1].Name)
	}
}

func TestPlanDeletePlan(t *testing.T) {
	d, vID := setupPlanDB(t)

	planID, err := d.CreateCustomPlan(vID, "to-delete", []PlanEntry{
		{DayNumber: 1, BookCode: "gen", ChapterStart: 1, ChapterEnd: 3},
		{DayNumber: 2, BookCode: "gen", ChapterStart: 4, ChapterEnd: 6},
	})
	if err != nil {
		t.Fatalf("CreateCustomPlan: %v", err)
	}

	if err := d.DeletePlan(planID); err != nil {
		t.Fatalf("DeletePlan: %v", err)
	}

	var count int
	err = d.conn.QueryRow("SELECT COUNT(*) FROM reading_plans WHERE id = ?", planID).Scan(&count)
	if err != nil {
		t.Fatalf("count plans: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 plans after delete, got %d", count)
	}

	err = d.conn.QueryRow("SELECT COUNT(*) FROM reading_plan_entries WHERE plan_id = ?", planID).Scan(&count)
	if err != nil {
		t.Fatalf("count entries: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 entries after delete, got %d", count)
	}
}
