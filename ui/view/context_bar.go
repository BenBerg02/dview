package view

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"

	"dview/ui/model"
	"dview/ui/style"
)

func contextBar(m *model.Model, w int) string {
	var parts []string
	for i, it := range model.CtxItems {
		st := style.Button
		if m.Focus == model.FocusContext && i == m.CtxCur {
			st = style.ActiveBtn
		}
		parts = append(parts, st.Render(" "+it+" "))
	}
	bar := strings.Join(parts, " ")
	if m.Table != "" {
		info := fmt.Sprintf("Page %d/%d · %d rows", m.Page+1, m.PageCount(), m.Total)
		if m.Where != "" {
			info += " · WHERE " + m.Where
		}
		bar += "  " + style.Dim.Render(info)
	}
	return ansi.Truncate(bar, w, "…")
}
