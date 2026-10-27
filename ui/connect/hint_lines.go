package connect

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"dview/ui/style"
)

var dirStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))

// maxHints is the number of hint lines reserved so the dialog does not jump.
func (m *Model) maxHints() int {
	return max(0, min(6, m.H-12))
}

func (m *Model) hintLines() []string {
	if !m.Drivers[m.Cur].CompletePath || m.maxHints() == 0 {
		return nil
	}
	width := m.inputWidth() + 2
	n := m.maxHints()
	lines := make([]string, 0, n+1)
	lines = append(lines, "")
	for i, e := range m.Hints {
		if i == n-1 && len(m.Hints) > n {
			lines = append(lines, style.Dim.Render(fmt.Sprintf("… %d more", len(m.Hints)-i)))
			break
		}
		if i == n {
			break
		}
		name, st := e.Name, style.Dim
		if e.Dir {
			name, st = name+"/", dirStyle
		}
		lines = append(lines, st.Render(ansi.Truncate(name, width, "…")))
	}
	for len(lines) < n+1 {
		lines = append(lines, "")
	}
	return lines
}
