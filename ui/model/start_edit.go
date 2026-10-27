package model

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) startEdit() tea.Cmd {
	v := m.Rows[m.Row].Values[m.Col]
	if v.Blob {
		m.fail(errBlob)
		return nil
	}
	return m.startInput(ModeEdit, "NULL", v.Text)
}
