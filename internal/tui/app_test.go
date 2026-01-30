package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestAppInit(t *testing.T) {
	m := New(nil)
	if m.state != StateBookList {
		t.Errorf("expected StateBookList, got %d", m.state)
	}
}

func TestAppKeyQ(t *testing.T) {
	m := New(nil)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("expected quit command")
	}
}

func TestAppKeyCtrlC(t *testing.T) {
	m := New(nil)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected quit command")
	}
}

func TestAppKeyHelp(t *testing.T) {
	m := New(nil)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	model := updated.(AppModel)
	if model.state != StateHelp {
		t.Errorf("expected StateHelp, got %d", model.state)
	}
	if model.prevState != StateBookList {
		t.Errorf("expected prevState StateBookList, got %d", model.prevState)
	}
}

func TestAppEscFromHelp(t *testing.T) {
	m := New(nil)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	model := updated.(AppModel)
	updated2, _ := model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	model2 := updated2.(AppModel)
	if model2.state != StateBookList {
		t.Errorf("expected StateBookList after Esc, got %d", model2.state)
	}
}

func TestAppKeySearch(t *testing.T) {
	m := New(nil)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model := updated.(AppModel)
	if model.state != StateSearch {
		t.Errorf("expected StateSearch, got %d", model.state)
	}
}

func TestAppKeyBookmarks(t *testing.T) {
	m := New(nil)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model := updated.(AppModel)
	if model.state != StateBookmarks {
		t.Errorf("expected StateBookmarks, got %d", model.state)
	}
}

func TestAppKeyB(t *testing.T) {
	m := New(nil)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model := updated.(AppModel)
	updated2, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	model2 := updated2.(AppModel)
	if model2.state != StateBookList {
		t.Errorf("expected StateBookList, got %d", model2.state)
	}
}

func TestAppWindowSize(t *testing.T) {
	m := New(nil)
	if m.ready {
		t.Error("expected not ready before WindowSizeMsg")
	}
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := updated.(AppModel)
	if !model.ready {
		t.Error("expected ready after WindowSizeMsg")
	}
	if model.width != 80 || model.height != 24 {
		t.Errorf("expected 80x24, got %dx%d", model.width, model.height)
	}
}

func TestAppViewNotReady(t *testing.T) {
	m := New(nil)
	v := m.View()
	if v != "Loading..." {
		t.Errorf("expected 'Loading...', got %q", v)
	}
}

func TestAppViewReady(t *testing.T) {
	m := New(nil)
	m.ready = true
	m.width = 80
	m.height = 24
	v := m.View()
	if v == "" {
		t.Error("expected non-empty view")
	}
}
