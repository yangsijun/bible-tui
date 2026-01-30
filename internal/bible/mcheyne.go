package bible

// McCheyneEntry represents one reading section for a day in the M'Cheyne plan.
type McCheyneEntry struct {
	BookCode     string
	ChapterStart int
	ChapterEnd   int
}

// chapterRef is a single chapter reference used internally.
type chapterRef struct {
	bookCode string
	chapter  int
}

// McCheyneSchedule returns 365 days of readings following Robert Murray
// M'Cheyne's Bible reading calendar. Each day has 2–4 reading sections.
//
// The plan has 4 parallel reading tracks:
//   - Track 1 (Family AM): Genesis→Deuteronomy, Joshua→Esther
//   - Track 2 (Secret AM): Acts→Revelation, Job→Song of Solomon
//   - Track 3 (Family PM): Matthew→John, Isaiah→Malachi
//   - Track 4 (Secret PM): Psalms, Romans→Revelation
//
// This covers the OT once, and reads the NT epistles and Psalms twice.
func McCheyneSchedule() [][]McCheyneEntry {
	tracks := [4][]string{
		// Track 1 — Family worship morning: Genesis through Esther
		{
			"gen", "exo", "lev", "num", "deu",
			"jos", "jdg", "rut", "1sa", "2sa",
			"1ki", "2ki", "1ch", "2ch", "ezr",
			"neh", "est",
		},
		// Track 2 — Secret morning: Acts through Revelation, then Job through Song of Solomon
		{
			"act", "rom", "1co", "2co", "gal",
			"eph", "php", "col", "1th", "2th",
			"1ti", "2ti", "tit", "phm", "heb",
			"jas", "1pe", "2pe", "1jn", "2jn",
			"3jn", "jud", "rev",
			"job", "psa", "pro", "ecc", "sng",
		},
		// Track 3 — Family worship evening: Matthew through John, then Isaiah through Malachi
		{
			"mat", "mrk", "luk", "jhn",
			"isa", "jer", "lam", "ezk", "dan",
			"hos", "jol", "amo", "oba", "jnh",
			"mic", "nam", "hab", "zep", "hag",
			"zec", "mal",
		},
		// Track 4 — Secret evening: Psalms, then Romans through Revelation
		{
			"psa",
			"rom", "1co", "2co", "gal", "eph",
			"php", "col", "1th", "2th", "1ti",
			"2ti", "tit", "phm", "heb", "jas",
			"1pe", "2pe", "1jn", "2jn", "3jn",
			"jud", "rev",
		},
	}

	var trackChapters [4][]chapterRef
	for t, bookCodes := range tracks {
		for _, code := range bookCodes {
			book, ok := GetBookByCode(code)
			if !ok {
				continue
			}
			for ch := 1; ch <= book.ChapterCount; ch++ {
				trackChapters[t] = append(trackChapters[t], chapterRef{code, ch})
			}
		}
	}

	days := make([][]McCheyneEntry, 365)
	for i := range days {
		days[i] = make([]McCheyneEntry, 0, 4)
	}

	for t := 0; t < 4; t++ {
		chs := trackChapters[t]
		n := len(chs)
		if n == 0 {
			continue
		}
		for day := 0; day < 365; day++ {
			startIdx := day * n / 365
			endIdx := (day + 1) * n / 365
			if startIdx >= endIdx {
				continue
			}
			entries := groupChapterRefs(chs[startIdx:endIdx])
			days[day] = append(days[day], entries...)
		}
	}

	return days
}

// groupChapterRefs merges consecutive chapters from the same book into
// single McCheyneEntry values.
func groupChapterRefs(refs []chapterRef) []McCheyneEntry {
	if len(refs) == 0 {
		return nil
	}
	var entries []McCheyneEntry
	cur := McCheyneEntry{
		BookCode:     refs[0].bookCode,
		ChapterStart: refs[0].chapter,
		ChapterEnd:   refs[0].chapter,
	}
	for i := 1; i < len(refs); i++ {
		r := refs[i]
		if r.bookCode == cur.BookCode && r.chapter == cur.ChapterEnd+1 {
			cur.ChapterEnd = r.chapter
		} else {
			entries = append(entries, cur)
			cur = McCheyneEntry{
				BookCode:     r.bookCode,
				ChapterStart: r.chapter,
				ChapterEnd:   r.chapter,
			}
		}
	}
	entries = append(entries, cur)
	return entries
}
