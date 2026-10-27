package view

import (
	"github.com/charmbracelet/lipgloss"

	"dview/ui/model"
	"dview/ui/style"
)

func paneBorder(active bool) lipgloss.Style {
	c := lipgloss.Color("240")
	if active {
		c = style.Accent
	}
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(c)
}

func panes(m *model.Model) string {
	lw := m.LeftWidth()
	ph := m.PaneHeight()
	left := paneBorder(m.Focus == model.FocusTables).Width(lw - 2).Height(ph).Render(tablesPane(m, lw-2, ph))
	right := paneBorder(m.Focus == model.FocusData || m.Focus == model.FocusContext).
		Width(m.RightInner()).Height(ph).Render(dataPane(m, m.RightInner()))
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
