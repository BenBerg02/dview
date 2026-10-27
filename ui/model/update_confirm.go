package model

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) updateConfirm(msg tea.KeyMsg) {
	switch msg.String() {
	case "y", "Y", "enter":
		m.Mode = ModeNone
		m.deleteRow()
	case "n", "N", "esc":
		m.Mode = ModeNone
	}
}
