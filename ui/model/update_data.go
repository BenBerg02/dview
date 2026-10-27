package model

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) updateData(msg tea.KeyMsg) tea.Cmd {
	n := len(m.Rows)
	switch msg.String() {
	case "up":
		if n > 0 {
			m.Row = cycle(m.Row, -1, n)
		}
	case "down":
		if n > 0 {
			m.Row = cycle(m.Row, 1, n)
		}
	case "left":
		m.Col = max(0, m.Col-1)
		m.fixColOff()
	case "right":
		m.Col = min(len(m.Cols)-1, m.Col+1)
		m.fixColOff()
	case "pgup":
		m.changePage(-1)
	case "pgdown":
		m.changePage(1)
	case "enter":
		if n > 0 {
			return m.startEdit()
		}
	case " ":
		if n > 0 {
			m.MenuCur, m.Mode = 0, ModeMenu
		}
	case "delete":
		if n > 0 {
			m.Mode = ModeConfirm
		}
	case "esc":
		m.Focus = FocusContext
	}
	return nil
}
