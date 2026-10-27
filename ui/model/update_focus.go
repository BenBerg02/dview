package model

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) updateFocus(msg tea.KeyMsg) tea.Cmd {
	switch m.Focus {
	case FocusTables:
		m.updateTables(msg)
	case FocusData:
		return m.updateData(msg)
	case FocusContext:
		return m.updateContext(msg)
	case FocusGlobal:
		switch msg.String() {
		case "enter":
			return tea.Quit
		case "esc", "down":
			m.Focus = FocusTables
		}
	}
	return nil
}
