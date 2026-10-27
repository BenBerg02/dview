package model

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) startSearch() tea.Cmd {
	return m.startInput(ModeSearch, "e.g. id > 10 AND name LIKE 'a%'", m.Where)
}
