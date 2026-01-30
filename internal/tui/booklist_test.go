package tui

import (
	"testing"
)

func TestBookListInit(t *testing.T) {
	m := NewBookList(80, 24)
	items := m.list.Items()
	if len(items) != 66 {
		t.Errorf("expected 66 items, got %d", len(items))
	}
}

func TestBookListTitle(t *testing.T) {
	m := NewBookList(80, 24)
	if m.list.Title != "성경 책 목록" {
		t.Errorf("unexpected title: %q", m.list.Title)
	}
}

func TestBookListItemTitle(t *testing.T) {
	m := NewBookList(80, 24)
	item, ok := m.list.Items()[0].(bookItem)
	if !ok {
		t.Fatal("expected bookItem type")
	}
	if item.Title() != "창세기" {
		t.Errorf("expected 창세기, got %q", item.Title())
	}
}

func TestBookListView(t *testing.T) {
	m := NewBookList(80, 24)
	v := m.View()
	if v == "" {
		t.Error("expected non-empty view")
	}
}
