package view

import (
	"github.com/charmbracelet/x/ansi"

	"dview/ui/model"
	"dview/ui/style"
)

func footer(m *model.Model) string {
	if m.Status != "" {
		st := style.OK
		if m.StatusErr {
			st = style.Err
		}
		return st.Render(ansi.Truncate(" "+m.Status, m.W, "…"))
	}
	return style.Footer(hint(m), m.W)
}
