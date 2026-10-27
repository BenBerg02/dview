package view

import (
	"strings"

	"dview/ui/model"
	"dview/ui/style"
)

func menuBox(m *model.Model) string {
	var lines []string
	for i, it := range model.MenuItems {
		if i == m.MenuCur {
			lines = append(lines, style.ActiveBtn.Render(style.Pad(" "+it, 16)))
		} else {
			lines = append(lines, style.Pad(" "+it, 16))
		}
	}
	return style.Box.Render(strings.Join(lines, "\n"))
}
