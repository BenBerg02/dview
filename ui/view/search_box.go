package view

import (
	"fmt"
	"strings"

	"dview/ui/model"
	"dview/ui/style"
)

func searchBox(m *model.Model) string {
	title := style.Header.Render("Search " + m.Table)
	return style.Box.Render(title + "\n" + style.Dim.Render("WHERE condition (empty clears)") + "\n\n" +
		m.Input.View() + "\n" + searchHints(m))
}

// maxSearchHints is the number of hint lines reserved so the dialog does not jump.
const maxSearchHints = 5

func searchHints(m *model.Model) string {
	hints := m.WhereHints()
	lines := make([]string, maxSearchHints)
	for i := range lines {
		switch {
		case i == maxSearchHints-1 && len(hints) > maxSearchHints:
			lines[i] = style.Dim.Render(fmt.Sprintf("… %d more", len(hints)-i))
		case i < len(hints):
			lines[i] = style.Dim.Render(hints[i])
		}
	}
	return strings.Join(lines, "\n")
}
