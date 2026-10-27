package model

import (
	tea "github.com/charmbracelet/bubbletea"

	"dview/db"
)

func (m *Model) updateMenu(msg tea.KeyMsg) tea.Cmd {
	n := len(MenuItems)
	switch msg.String() {
	case "up":
		m.MenuCur = cycle(m.MenuCur, -1, n)
	case "down":
		m.MenuCur = cycle(m.MenuCur, 1, n)
	case "esc":
		m.Mode = ModeNone
	case "enter":
		m.Mode = ModeNone
		switch m.MenuCur {
		case 0:
			return m.startEdit()
		case 1:
			m.setCell(db.Value{Null: true})
		case 2:
			m.Mode = ModeConfirm
		}
	}
	return nil
}
