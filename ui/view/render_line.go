package view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"dview/ui/model"
	"dview/ui/style"
)

// renderLine draws one table line, scrolled horizontally from m.ColOff.
func renderLine(m *model.Model, cells []string, ws []int, base lipgloss.Style, selCol int, sel lipgloss.Style, avail int) string {
	var b strings.Builder
	used := 0
	for i := m.ColOff; i < len(ws); i++ {
		if i > m.ColOff {
			if avail-used-3 < 1 {
				break
			}
			b.WriteString(base.Render(" │ "))
			used += 3
		}
		w := min(ws[i], avail-used)
		st := base
		if i == selCol {
			st = sel
		}
		b.WriteString(st.Render(style.Pad(cells[i], w)))
		used += w
		if used >= avail {
			break
		}
	}
	if used < avail {
		b.WriteString(base.Render(strings.Repeat(" ", avail-used)))
	}
	return b.String()
}
