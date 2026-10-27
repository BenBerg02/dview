package view

import (
	"fmt"

	"dview/ui/model"
	"dview/ui/style"
)

func confirmBox(m *model.Model) string {
	return style.Box.Render(style.Header.Render("Delete row?") + "\n\n" +
		fmt.Sprintf("This removes the selected row from %s.", m.Table) + "\n" +
		style.Dim.Render("y = delete · n = cancel"))
}
