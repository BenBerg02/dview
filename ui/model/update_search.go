package model

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) updateSearch(msg tea.KeyMsg) tea.Cmd {
	if msg.String() == "tab" {
		m.completeWhere()
		return nil
	}
	return m.updateInput(msg, m.applyWhere)
}
