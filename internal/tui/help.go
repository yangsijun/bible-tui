package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yangsijun/bible-tui/internal/tui/styles"
)

type HelpModel struct {
	viewport viewport.Model
	theme    *styles.Theme
	width    int
	height   int
}

func NewHelp(theme *styles.Theme, width, height int) HelpModel {
	vp := viewport.New(width, height-2)
	vp.SetContent(renderHelpContent(theme))
	return HelpModel{
		viewport: vp,
		theme:    theme,
		width:    width,
		height:   height,
	}
}

func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m HelpModel) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Primary).Padding(0, 1)
	title := titleStyle.Render("도움말 — 키바인딩")
	return title + "\n" + m.viewport.View()
}

func renderHelpContent(theme *styles.Theme) string {
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.Secondary)
	keyStyle := lipgloss.NewStyle().Foreground(theme.Primary).Width(14)
	descStyle := lipgloss.NewStyle().Foreground(theme.Foreground)

	var b strings.Builder

	b.WriteString(sectionStyle.Render("전역 키바인딩"))
	b.WriteString("\n\n")
	keys := [][2]string{
		{"q, Ctrl+C", "종료"},
		{"?", "도움말 열기/닫기"},
		{"b", "책 목록으로 이동"},
		{"/", "검색"},
		{"m", "책갈피/하이라이트"},
		{"Esc", "이전 화면"},
	}
	for _, kv := range keys {
		b.WriteString("  " + keyStyle.Render(kv[0]) + descStyle.Render(kv[1]) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("책 목록"))
	b.WriteString("\n\n")
	keys2 := [][2]string{
		{"j, k, ↑, ↓", "위/아래 이동"},
		{"Enter", "책 선택"},
		{"/", "필터 검색"},
	}
	for _, kv := range keys2 {
		b.WriteString("  " + keyStyle.Render(kv[0]) + descStyle.Render(kv[1]) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("장 선택"))
	b.WriteString("\n\n")
	keys3 := [][2]string{
		{"h, l, ←, →", "좌/우 이동"},
		{"j, k, ↑, ↓", "위/아래 이동"},
		{"Enter", "장 선택"},
		{"Esc", "책 목록으로"},
	}
	for _, kv := range keys3 {
		b.WriteString("  " + keyStyle.Render(kv[0]) + descStyle.Render(kv[1]) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("읽기 화면"))
	b.WriteString("\n\n")
	keys4 := [][2]string{
		{"h, ←", "이전 장"},
		{"l, →", "다음 장"},
		{"j, k", "위/아래 스크롤"},
		{"g", "맨 위"},
		{"G", "맨 아래"},
		{"Esc", "장 선택으로"},
	}
	for _, kv := range keys4 {
		b.WriteString("  " + keyStyle.Render(kv[0]) + descStyle.Render(kv[1]) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("검색"))
	b.WriteString("\n\n")
	keys5 := [][2]string{
		{"Enter", "검색 실행"},
		{"j, k", "결과 탐색"},
		{"Enter", "결과 선택"},
		{"/", "검색창으로"},
	}
	for _, kv := range keys5 {
		b.WriteString("  " + keyStyle.Render(kv[0]) + descStyle.Render(kv[1]) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("책갈피/하이라이트"))
	b.WriteString("\n\n")
	keys6 := [][2]string{
		{"Tab", "탭 전환"},
		{"j, k", "위/아래 이동"},
		{"d", "삭제"},
		{"Enter", "해당 구절로 이동"},
	}
	for _, kv := range keys6 {
		b.WriteString("  " + keyStyle.Render(kv[0]) + descStyle.Render(kv[1]) + "\n")
	}

	return b.String()
}
