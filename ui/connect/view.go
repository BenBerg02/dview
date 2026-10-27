package connect

import (
	"strings"

	"dview/ui/style"
)

func (m *Model) View() string {
	if m.W == 0 || m.H == 0 {
		return ""
	}
	b, x, y := m.box()
	pre := strings.Repeat(" ", x)
	lines := make([]string, 0, m.H)
	for i := 0; i < y; i++ {
		lines = append(lines, "")
	}
	for _, l := range strings.Split(b, "\n") {
		lines = append(lines, pre+l)
	}
	for len(lines) < m.H-1 {
		lines = append(lines, "")
	}
	lines = lines[:m.H-1]
	return strings.Join(lines, "\n") + "\n" + style.Footer(hint(m), m.W)
}
