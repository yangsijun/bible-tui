package bible

import "strings"

// BookInfo contains metadata for a single Bible book
type BookInfo struct {
	Code         string // English code: "gen", "exo", etc.
	NameKo       string // Korean full name: "창세기"
	AbbrevKo     string // Korean abbreviation: "창"
	Testament    string // "old" or "new"
	ChapterCount int    // Number of chapters
}

var allBooks = []BookInfo{
	// Old Testament (39 books)
	{Code: "gen", NameKo: "창세기", AbbrevKo: "창", Testament: "old", ChapterCount: 50},
	{Code: "exo", NameKo: "출애굽기", AbbrevKo: "출", Testament: "old", ChapterCount: 40},
	{Code: "lev", NameKo: "레위기", AbbrevKo: "레", Testament: "old", ChapterCount: 27},
	{Code: "num", NameKo: "민수기", AbbrevKo: "민", Testament: "old", ChapterCount: 36},
	{Code: "deu", NameKo: "신명기", AbbrevKo: "신", Testament: "old", ChapterCount: 34},
	{Code: "jos", NameKo: "여호수아", AbbrevKo: "수", Testament: "old", ChapterCount: 24},
	{Code: "jdg", NameKo: "사사기", AbbrevKo: "삿", Testament: "old", ChapterCount: 21},
	{Code: "rut", NameKo: "룻기", AbbrevKo: "룻", Testament: "old", ChapterCount: 4},
	{Code: "1sa", NameKo: "사무엘상", AbbrevKo: "삼상", Testament: "old", ChapterCount: 31},
	{Code: "2sa", NameKo: "사무엘하", AbbrevKo: "삼하", Testament: "old", ChapterCount: 24},
	{Code: "1ki", NameKo: "열왕기상", AbbrevKo: "왕상", Testament: "old", ChapterCount: 22},
	{Code: "2ki", NameKo: "열왕기하", AbbrevKo: "왕하", Testament: "old", ChapterCount: 25},
	{Code: "1ch", NameKo: "역대상", AbbrevKo: "대상", Testament: "old", ChapterCount: 29},
	{Code: "2ch", NameKo: "역대하", AbbrevKo: "대하", Testament: "old", ChapterCount: 36},
	{Code: "ezr", NameKo: "에스라", AbbrevKo: "스", Testament: "old", ChapterCount: 10},
	{Code: "neh", NameKo: "느헤미야", AbbrevKo: "느", Testament: "old", ChapterCount: 13},
	{Code: "est", NameKo: "에스더", AbbrevKo: "에", Testament: "old", ChapterCount: 10},
	{Code: "job", NameKo: "욥기", AbbrevKo: "욥", Testament: "old", ChapterCount: 42},
	{Code: "psa", NameKo: "시편", AbbrevKo: "시", Testament: "old", ChapterCount: 150},
	{Code: "pro", NameKo: "잠언", AbbrevKo: "잠", Testament: "old", ChapterCount: 31},
	{Code: "ecc", NameKo: "전도서", AbbrevKo: "전", Testament: "old", ChapterCount: 12},
	{Code: "sng", NameKo: "아가", AbbrevKo: "아", Testament: "old", ChapterCount: 8},
	{Code: "isa", NameKo: "이사야", AbbrevKo: "사", Testament: "old", ChapterCount: 66},
	{Code: "jer", NameKo: "예레미야", AbbrevKo: "렘", Testament: "old", ChapterCount: 52},
	{Code: "lam", NameKo: "예레미야애가", AbbrevKo: "애", Testament: "old", ChapterCount: 5},
	{Code: "ezk", NameKo: "에스겔", AbbrevKo: "겔", Testament: "old", ChapterCount: 48},
	{Code: "dan", NameKo: "다니엘", AbbrevKo: "단", Testament: "old", ChapterCount: 12},
	{Code: "hos", NameKo: "호세아", AbbrevKo: "호", Testament: "old", ChapterCount: 14},
	{Code: "jol", NameKo: "요엘", AbbrevKo: "욜", Testament: "old", ChapterCount: 3},
	{Code: "amo", NameKo: "아모스", AbbrevKo: "암", Testament: "old", ChapterCount: 9},
	{Code: "oba", NameKo: "오바댜", AbbrevKo: "옵", Testament: "old", ChapterCount: 1},
	{Code: "jnh", NameKo: "요나", AbbrevKo: "욘", Testament: "old", ChapterCount: 4},
	{Code: "mic", NameKo: "미가", AbbrevKo: "미", Testament: "old", ChapterCount: 7},
	{Code: "nam", NameKo: "나훔", AbbrevKo: "나", Testament: "old", ChapterCount: 3},
	{Code: "hab", NameKo: "하박국", AbbrevKo: "합", Testament: "old", ChapterCount: 3},
	{Code: "zep", NameKo: "스바냐", AbbrevKo: "습", Testament: "old", ChapterCount: 3},
	{Code: "hag", NameKo: "학개", AbbrevKo: "학", Testament: "old", ChapterCount: 2},
	{Code: "zec", NameKo: "스가랴", AbbrevKo: "슥", Testament: "old", ChapterCount: 14},
	{Code: "mal", NameKo: "말라기", AbbrevKo: "말", Testament: "old", ChapterCount: 4},

	// New Testament (27 books)
	{Code: "mat", NameKo: "마태복음", AbbrevKo: "마", Testament: "new", ChapterCount: 28},
	{Code: "mrk", NameKo: "마가복음", AbbrevKo: "막", Testament: "new", ChapterCount: 16},
	{Code: "luk", NameKo: "누가복음", AbbrevKo: "눅", Testament: "new", ChapterCount: 24},
	{Code: "jhn", NameKo: "요한복음", AbbrevKo: "요", Testament: "new", ChapterCount: 21},
	{Code: "act", NameKo: "사도행전", AbbrevKo: "행", Testament: "new", ChapterCount: 28},
	{Code: "rom", NameKo: "로마서", AbbrevKo: "롬", Testament: "new", ChapterCount: 16},
	{Code: "1co", NameKo: "고린도전서", AbbrevKo: "고전", Testament: "new", ChapterCount: 16},
	{Code: "2co", NameKo: "고린도후서", AbbrevKo: "고후", Testament: "new", ChapterCount: 13},
	{Code: "gal", NameKo: "갈라디아서", AbbrevKo: "갈", Testament: "new", ChapterCount: 6},
	{Code: "eph", NameKo: "에베소서", AbbrevKo: "엡", Testament: "new", ChapterCount: 6},
	{Code: "php", NameKo: "빌립보서", AbbrevKo: "빌", Testament: "new", ChapterCount: 4},
	{Code: "col", NameKo: "골로새서", AbbrevKo: "골", Testament: "new", ChapterCount: 4},
	{Code: "1th", NameKo: "데살로니가전서", AbbrevKo: "살전", Testament: "new", ChapterCount: 5},
	{Code: "2th", NameKo: "데살로니가후서", AbbrevKo: "살후", Testament: "new", ChapterCount: 3},
	{Code: "1ti", NameKo: "디모데전서", AbbrevKo: "딤전", Testament: "new", ChapterCount: 6},
	{Code: "2ti", NameKo: "디모데후서", AbbrevKo: "딤후", Testament: "new", ChapterCount: 4},
	{Code: "tit", NameKo: "디도서", AbbrevKo: "딛", Testament: "new", ChapterCount: 3},
	{Code: "phm", NameKo: "빌레몬서", AbbrevKo: "몬", Testament: "new", ChapterCount: 1},
	{Code: "heb", NameKo: "히브리서", AbbrevKo: "히", Testament: "new", ChapterCount: 13},
	{Code: "jas", NameKo: "야고보서", AbbrevKo: "약", Testament: "new", ChapterCount: 5},
	{Code: "1pe", NameKo: "베드로전서", AbbrevKo: "벧전", Testament: "new", ChapterCount: 5},
	{Code: "2pe", NameKo: "베드로후서", AbbrevKo: "벧후", Testament: "new", ChapterCount: 3},
	{Code: "1jn", NameKo: "요한1서", AbbrevKo: "요일", Testament: "new", ChapterCount: 5},
	{Code: "2jn", NameKo: "요한2서", AbbrevKo: "요이", Testament: "new", ChapterCount: 1},
	{Code: "3jn", NameKo: "요한3서", AbbrevKo: "요삼", Testament: "new", ChapterCount: 1},
	{Code: "jud", NameKo: "유다서", AbbrevKo: "유", Testament: "new", ChapterCount: 1},
	{Code: "rev", NameKo: "요한계시록", AbbrevKo: "계", Testament: "new", ChapterCount: 22},
}

// AllBooks returns a copy of all 66 Bible books
func AllBooks() []BookInfo {
	result := make([]BookInfo, len(allBooks))
	copy(result, allBooks)
	return result
}

// GetBookByCode looks up a book by its English code (case-insensitive)
func GetBookByCode(code string) (*BookInfo, bool) {
	lowerCode := strings.ToLower(code)
	for i := range allBooks {
		if allBooks[i].Code == lowerCode {
			return &allBooks[i], true
		}
	}
	return nil, false
}

// GetBookByName looks up a book by its Korean full name
func GetBookByName(name string) (*BookInfo, bool) {
	for i := range allBooks {
		if allBooks[i].NameKo == name {
			return &allBooks[i], true
		}
	}
	return nil, false
}

// GetBookByAbbrev looks up a book by its Korean abbreviation
func GetBookByAbbrev(abbrev string) (*BookInfo, bool) {
	for i := range allBooks {
		if allBooks[i].AbbrevKo == abbrev {
			return &allBooks[i], true
		}
	}
	return nil, false
}
