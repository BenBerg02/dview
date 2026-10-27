package model

import "github.com/charmbracelet/x/ansi"

// LeftWidth is the outer width of the tables pane.
func (m *Model) LeftWidth() int {
	w := 0
	for _, t := range m.Tables {
		w = max(w, ansi.StringWidth(t))
	}
	w += 4
	return max(14, min(w, m.W/3))
}
