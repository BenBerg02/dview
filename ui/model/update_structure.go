package model

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) updateStructure(msg tea.KeyMsg) {
	switch msg.String() {
	case "esc", "enter", "q":
		m.Mode = ModeNone
	}
}
