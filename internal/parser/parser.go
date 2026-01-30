package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type ChapterData struct {
	Verses []VerseData
}

type VerseData struct {
	Number       int
	Text         string
	SectionTitle string
	Footnotes    []FootnoteData
}

type FootnoteData struct {
	Marker  string
	Content string
}

var spaceNormalizer = regexp.MustCompile(`\s+`)

func ParseChapterHTML(htmlStr string) (*ChapterData, error) {
	if strings.TrimSpace(htmlStr) == "" {
		return nil, fmt.Errorf("empty html input")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	container := doc.Find("div#tdBible1")
	if container.Length() == 0 {
		return nil, fmt.Errorf("container div#tdBible1 not found")
	}

	state := &parseState{
		currentVerseIdx: -1,
	}

	for _, node := range container.Nodes {
		walkNode(node, state)
	}

	for i := range state.verses {
		text := state.verses[i].Text
		text = spaceNormalizer.ReplaceAllString(text, " ")
		text = strings.TrimSpace(text)
		state.verses[i].Text = text
	}

	return &ChapterData{Verses: state.verses}, nil
}

type parseState struct {
	verses          []VerseData
	currentVerseIdx int
	pendingTitle    string
}

func (ps *parseState) startVerse(num int) {
	v := VerseData{
		Number:       num,
		SectionTitle: ps.pendingTitle,
	}
	ps.pendingTitle = ""
	ps.verses = append(ps.verses, v)
	ps.currentVerseIdx = len(ps.verses) - 1
}

func (ps *parseState) addText(text string) {
	if ps.currentVerseIdx >= 0 && ps.currentVerseIdx < len(ps.verses) {
		ps.verses[ps.currentVerseIdx].Text += text
	}
}

func (ps *parseState) addFootnoteMarker(marker string) {
	if ps.currentVerseIdx >= 0 && ps.currentVerseIdx < len(ps.verses) {
		ps.verses[ps.currentVerseIdx].Footnotes = append(
			ps.verses[ps.currentVerseIdx].Footnotes,
			FootnoteData{Marker: marker},
		)
	}
}

func (ps *parseState) setFootnoteContent(content string) {
	if ps.currentVerseIdx < 0 || ps.currentVerseIdx >= len(ps.verses) {
		return
	}
	v := &ps.verses[ps.currentVerseIdx]
	for i := len(v.Footnotes) - 1; i >= 0; i-- {
		if v.Footnotes[i].Content == "" {
			v.Footnotes[i].Content = content
			return
		}
	}
}

func walkNode(n *html.Node, state *parseState) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "font":
			if hasClass(n, "smallTitle") {
				state.pendingTitle = strings.TrimSpace(collectText(n))
				return
			}
			if hasStyleContaining(n, "display:none") {
				return
			}
		case "span":
			if hasClass(n, "number") {
				raw := collectText(n)
				raw = strings.ReplaceAll(raw, "\u00a0", "")
				raw = strings.TrimSpace(raw)
				if num, err := strconv.Atoi(raw); err == nil {
					state.startVerse(num)
				}
				return
			}
		case "a":
			if hasClass(n, "comment") {
				marker := strings.TrimSpace(collectText(n))
				if marker != "" {
					state.addFootnoteMarker(marker)
				}
				return
			}
		case "div":
			if hasClass(n, "D2") {
				content := strings.TrimSpace(collectText(n))
				if content != "" {
					state.setFootnoteContent(content)
				}
				return
			}
		case "img", "script":
			return
		}
	}

	if n.Type == html.TextNode {
		state.addText(n.Data)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkNode(c, state)
	}
}

func hasClass(n *html.Node, name string) bool {
	for _, a := range n.Attr {
		if a.Key == "class" {
			for _, c := range strings.Fields(a.Val) {
				if c == name {
					return true
				}
			}
		}
	}
	return false
}

func hasStyleContaining(n *html.Node, substr string) bool {
	for _, a := range n.Attr {
		if a.Key == "style" && strings.Contains(a.Val, substr) {
			return true
		}
	}
	return false
}

func collectText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var b strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b.WriteString(collectText(c))
	}
	return b.String()
}
