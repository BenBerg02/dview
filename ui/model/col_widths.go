package model

import "github.com/charmbracelet/x/ansi"

// ColWidths returns the display width of each column for the current page.
func (m *Model) ColWidths() []int {
	avail := m.RightInner()
	ws := make([]int, len(m.Cols))
	for i, c := range m.Cols {
		w := ansi.StringWidth(c.Name)
		for _, r := range m.Rows {
			w = max(w, ansi.StringWidth(Display(r.Values[i])))
		}
		ws[i] = min(w, maxColWidth, avail)
	}
	return ws
}
