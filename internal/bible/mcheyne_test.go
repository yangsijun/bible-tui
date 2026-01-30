package bible

import "testing"

func TestMcCheyneSchedule(t *testing.T) {
	schedule := McCheyneSchedule()

	if len(schedule) != 365 {
		t.Fatalf("expected 365 days, got %d", len(schedule))
	}

	for day, entries := range schedule {
		if len(entries) < 2 {
			t.Errorf("day %d: expected at least 2 entries, got %d", day+1, len(entries))
		}
	}

	day1 := schedule[0]
	if day1[0].BookCode != "gen" {
		t.Errorf("day 1 first entry: expected 'gen', got %q", day1[0].BookCode)
	}
	if day1[0].ChapterStart != 1 {
		t.Errorf("day 1 first entry: expected chapter 1, got %d", day1[0].ChapterStart)
	}
}

func TestMcCheyneBookCodes(t *testing.T) {
	schedule := McCheyneSchedule()

	validCodes := make(map[string]bool)
	for _, b := range AllBooks() {
		validCodes[b.Code] = true
	}

	seen := make(map[string]bool)
	for day, entries := range schedule {
		for _, e := range entries {
			if !validCodes[e.BookCode] {
				t.Errorf("day %d: invalid book code %q", day+1, e.BookCode)
			}
			seen[e.BookCode] = true
			if e.ChapterStart < 1 {
				t.Errorf("day %d: chapter_start %d < 1 for %s", day+1, e.ChapterStart, e.BookCode)
			}
			if e.ChapterEnd < e.ChapterStart {
				t.Errorf("day %d: chapter_end %d < chapter_start %d for %s",
					day+1, e.ChapterEnd, e.ChapterStart, e.BookCode)
			}
		}
	}

	if len(seen) == 0 {
		t.Error("no book codes seen in schedule")
	}
}
