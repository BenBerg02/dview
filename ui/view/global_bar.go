package view

import (
	"strings"

	"github.com/charmbracelet/x/ansi"

	"dview/ui/model"
	"dview/ui/style"
)

func globalBar(m *model.Model) string {
	left := style.Header.Render(" dview ") + style.Dim.Render(m.DB.Name())
	quit := style.Button.Render(" Quit ")
	if m.Focus == model.FocusGlobal {
		quit = style.ActiveBtn.Render(" Quit ")
	}
	gap := m.W - ansi.StringWidth(left) - ansi.StringWidth(quit) - 1
	if gap < 1 {
		left = ansi.Truncate(left, max(0, m.W-ansi.StringWidth(quit)-2), "…")
		gap = max(1, m.W-ansi.StringWidth(left)-ansi.StringWidth(quit)-1)
	}
	return left + strings.Repeat(" ", gap) + quit + " "
}
