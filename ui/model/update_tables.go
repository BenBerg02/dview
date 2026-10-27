package model

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) updateTables(msg tea.KeyMsg) {
	n := len(m.Tables)
	switch msg.String() {
	case "up":
		if n > 0 {
			m.openTable(cycle(m.TableCur, -1, n))
		}
	case "down":
		if n > 0 {
			m.openTable(cycle(m.TableCur, 1, n))
		}
	case "right", "enter":
		if m.Table != "" {
			m.Focus = FocusData
		}
	case "esc":
		m.Focus = FocusGlobal
	}
	m.fixTableOff()
}
