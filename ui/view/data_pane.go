package view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"dview/ui/model"
	"dview/ui/style"
)

func dataPane(m *model.Model, w int) string {
	lines := []string{contextBar(m, w)}
	if m.Table == "" {
		return strings.Join(lines, "\n")
	}
	ws := m.ColWidths()
	names := make([]string, len(m.Cols))
	for i, c := range m.Cols {
		names[i] = c.Name
	}
	lines = append(lines, renderLine(m, names, ws, style.Header, -1, style.Header, w))
	if len(m.Rows) == 0 {
		lines = append(lines, style.Dim.Render(" (no rows)"))
	}
	dataFocused := m.Focus == model.FocusData
	for i, r := range m.Rows {
		cells := make([]string, len(r.Values))
		for j, v := range r.Values {
			cells[j] = model.Display(v)
		}
		rowSt, selCol := lipgloss.NewStyle(), -1
		if i == m.Row {
			rowSt = style.RowSel
			if dataFocused {
				selCol = m.Col
			}
		}
		lines = append(lines, renderLine(m, cells, ws, rowSt, selCol, style.CellSel, w))
	}
	return strings.Join(lines, "\n")
}
