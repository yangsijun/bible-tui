package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/yangsijun/bible-tui/internal/bible"
)

type BookSelectedMsg struct {
	Book bible.BookInfo
}

type bookItem struct {
	info bible.BookInfo
}

func (i bookItem) Title() string       { return i.info.NameKo }
func (i bookItem) Description() string { return fmt.Sprintf("%s • %d장", i.info.AbbrevKo, i.info.ChapterCount) }
func (i bookItem) FilterValue() string { return i.info.NameKo + " " + i.info.AbbrevKo }

type BookListModel struct {
	list list.Model
}

func NewBookList(width, height int) BookListModel {
	books := bible.AllBooks()
	items := make([]list.Item, len(books))
	for i, b := range books {
		items[i] = bookItem{info: b}
	}
	l := list.New(items, list.NewDefaultDelegate(), width, height-3)
	l.Title = "성경 책 목록"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	return BookListModel{list: l}
}

func (m BookListModel) Update(msg tea.Msg) (BookListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" && !m.list.SettingFilter() {
			if item, ok := m.list.SelectedItem().(bookItem); ok {
				return m, func() tea.Msg { return BookSelectedMsg{Book: item.info} }
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m BookListModel) View() string {
	return m.list.View()
}
