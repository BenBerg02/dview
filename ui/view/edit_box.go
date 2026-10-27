package view

import (
	"fmt"

	"dview/ui/model"
	"dview/ui/style"
)

func editBox(m *model.Model) string {
	c := m.Cols[m.Col]
	title := style.Header.Render(fmt.Sprintf("Edit %s.%s", m.Table, c.Name)) + style.Dim.Render(" ("+c.Type+")")
	return style.Box.Render(title + "\n\n" + m.Input.View())
}
