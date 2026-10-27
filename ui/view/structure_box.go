package view

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"

	"dview/ui/model"
	"dview/ui/style"
)

func structureBox(m *model.Model) string {
	nameW, typeW := 4, 4
	for _, c := range m.Cols {
		nameW = max(nameW, ansi.StringWidth(c.Name))
		typeW = max(typeW, ansi.StringWidth(c.Type))
	}
	lines := []string{style.Header.Render("Structure of " + m.Table), "",
		style.Dim.Render(style.Pad("Name", nameW) + "  " + style.Pad("Type", typeW) + "  Constraints")}
	limit := max(1, m.H-10)
	for i, c := range m.Cols {
		if i == limit {
			lines = append(lines, style.Dim.Render(fmt.Sprintf("… %d more", len(m.Cols)-limit)))
			break
		}
		var flags []string
		if c.PrimaryKey {
			flags = append(flags, "PK")
		}
		if c.Auto {
			flags = append(flags, "AUTO")
		}
		if c.NotNull {
			flags = append(flags, "NOT NULL")
		}
		if c.HasDefault {
			flags = append(flags, "DEFAULT "+c.Default)
		}
		lines = append(lines, style.Pad(c.Name, nameW)+"  "+style.Pad(c.Type, typeW)+"  "+strings.Join(flags, " "))
	}
	return style.Box.Render(strings.Join(lines, "\n"))
}
