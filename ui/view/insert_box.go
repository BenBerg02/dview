package view

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"dview/ui/model"
	"dview/ui/style"
)

func insertBox(m *model.Model) string {
	lw, tw, iw := m.InsertLayout()
	legend := style.OK.Render("green") + style.Dim.Render(" auto/default · ") +
		style.Err.Render("*") + style.Dim.Render(" required")
	lines := []string{style.Header.Render("Insert into " + m.Table), legend, ""}

	off, vis := m.InsertOff, m.InsertVisible()
	end := min(len(m.Cols), off+vis)
	for i := off; i < end; i++ {
		c := m.Cols[i]
		name, mark := lipgloss.NewStyle(), " "
		if c.AutoFill() {
			name = style.OK
		} else if c.Required() {
			mark = style.Err.Render("*")
		}
		lines = append(lines, name.Render(style.Pad(c.Name, lw-1))+mark+" "+
			style.Dim.Render(style.Pad(c.Type, tw))+" "+m.Inserts[i].View())
	}
	if len(m.Cols) > vis {
		lines = append(lines, style.Dim.Render(fmt.Sprintf("columns %d-%d of %d", off+1, end, len(m.Cols))))
	} else {
		lines = append(lines, "")
	}

	var btns []string
	for i, b := range model.InsertButtons {
		st := style.Button
		if m.InsertCur == len(m.Cols)+i {
			st = style.ActiveBtn
		}
		btns = append(btns, st.Render(" "+b+" "))
	}
	lines = append(lines, "", strings.Join(btns, " "), "")

	// the error line is always present so the dialog does not jump
	if m.InsertErr != "" {
		lines = append(lines, style.Err.Render(ansi.Truncate(m.InsertErr, lw+tw+iw+5, "…")))
	} else {
		lines = append(lines, "")
	}
	return style.Box.Render(strings.Join(lines, "\n"))
}
