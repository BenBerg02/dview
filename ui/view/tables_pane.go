package view

import (
	"strings"

	"dview/ui/model"
	"dview/ui/style"
)

func tablesPane(m *model.Model, w, h int) string {
	var lines []string
	for i := m.TableOff; i < len(m.Tables) && i < m.TableOff+h; i++ {
		line := style.Pad(" "+m.Tables[i], w)
		switch {
		case i == m.TableCur && m.Focus == model.FocusTables:
			line = style.ActiveBtn.Render(line)
		case i == m.TableCur:
			line = style.RowSel.Render(line)
		}
		lines = append(lines, line)
	}
	if len(m.Tables) == 0 {
		lines = append(lines, style.Dim.Render(" (no tables)"))
	}
	return strings.Join(lines, "\n")
}
