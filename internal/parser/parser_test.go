package parser

import (
	"os"
	"strings"
	"testing"
)

func loadFixture(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile("../../testdata/genesis_1.html")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}
	return string(data)
}

func parseFixture(t *testing.T) *ChapterData {
	t.Helper()
	htmlStr := loadFixture(t)
	result, err := ParseChapterHTML(htmlStr)
	if err != nil {
		t.Fatalf("ParseChapterHTML failed: %v", err)
	}
	return result
}

func findVerse(t *testing.T, data *ChapterData, num int) VerseData {
	t.Helper()
	for _, v := range data.Verses {
		if v.Number == num {
			return v
		}
	}
	t.Fatalf("verse %d not found", num)
	return VerseData{}
}

func TestParseGenesis1_VerseCount(t *testing.T) {
	data := parseFixture(t)
	if len(data.Verses) != 31 {
		t.Errorf("expected 31 verses, got %d", len(data.Verses))
		for i, v := range data.Verses {
			t.Logf("  verse[%d]: Number=%d Text=%q", i, v.Number, v.Text[:min(len(v.Text), 40)])
		}
	}
}

func TestParseGenesis1_Verse1Text(t *testing.T) {
	data := parseFixture(t)
	v := findVerse(t, data, 1)
	expected := "태초에 하나님이 천지를 창조하시니라"
	if !strings.Contains(v.Text, expected) {
		t.Errorf("verse 1 text does not contain %q\ngot: %q", expected, v.Text)
	}
}

func TestParseGenesis1_Verse2Text(t *testing.T) {
	data := parseFixture(t)
	v := findVerse(t, data, 2)
	if !strings.Contains(v.Text, "혼돈하고 공허하며") {
		t.Errorf("verse 2 text should contain '혼돈하고 공허하며'\ngot: %q", v.Text)
	}
	if strings.Contains(v.Text, "1)") {
		t.Errorf("verse 2 text should not contain footnote marker '1)'\ngot: %q", v.Text)
	}
}

func TestParseGenesis1_SectionTitle(t *testing.T) {
	data := parseFixture(t)
	v := findVerse(t, data, 1)
	if v.SectionTitle != "천지 창조" {
		t.Errorf("verse 1 section title: expected %q, got %q", "천지 창조", v.SectionTitle)
	}
}

func TestParseGenesis1_Footnotes(t *testing.T) {
	data := parseFixture(t)
	v := findVerse(t, data, 2)
	if len(v.Footnotes) != 1 {
		t.Fatalf("verse 2: expected 1 footnote, got %d", len(v.Footnotes))
	}
	fn := v.Footnotes[0]
	if fn.Marker != "1)" {
		t.Errorf("footnote marker: expected %q, got %q", "1)", fn.Marker)
	}
	if !strings.Contains(fn.Content, "또는 형체가 없는") {
		t.Errorf("footnote content should contain '또는 형체가 없는'\ngot: %q", fn.Content)
	}
}

func TestParseGenesis1_FootnoteVerse14(t *testing.T) {
	data := parseFixture(t)
	v := findVerse(t, data, 14)
	if len(v.Footnotes) == 0 {
		t.Fatal("verse 14: expected at least 1 footnote")
	}
	found := false
	for _, fn := range v.Footnotes {
		if strings.Contains(fn.Content, "또는 발광체") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("verse 14: no footnote contains '또는 발광체'")
	}
}

func TestParseGenesis1_FootnoteVerse26(t *testing.T) {
	data := parseFixture(t)
	v := findVerse(t, data, 26)
	if len(v.Footnotes) == 0 {
		t.Fatal("verse 26: expected at least 1 footnote")
	}
	found := false
	for _, fn := range v.Footnotes {
		if strings.Contains(fn.Content, "온 땅의 짐승과") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("verse 26: no footnote contains '온 땅의 짐승과'")
	}
}

func TestParseGenesis1_Verse16SmallText(t *testing.T) {
	data := parseFixture(t)
	v := findVerse(t, data, 16)
	if !strings.Contains(v.Text, "만드시고") {
		t.Errorf("verse 16 text should contain '만드시고'\ngot: %q", v.Text)
	}
}

func TestParseGenesis1_AllVersesNonEmpty(t *testing.T) {
	data := parseFixture(t)
	for _, v := range data.Verses {
		if strings.TrimSpace(v.Text) == "" {
			t.Errorf("verse %d has empty text", v.Number)
		}
	}
}

func TestParseEmptyHTML(t *testing.T) {
	_, err := ParseChapterHTML("")
	if err == nil {
		t.Error("expected error for empty HTML")
	}
}

func TestParseNoContainer(t *testing.T) {
	htmlStr := "<html><body><div>no bible content</div></body></html>"
	_, err := ParseChapterHTML(htmlStr)
	if err == nil {
		t.Error("expected error for HTML without div#tdBible1")
	}
}
