package connect

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) startLocation() tea.Cmd {
	d := m.Drivers[m.Cur]
	if d.Form != nil {
		return m.startForm()
	}
	m.Step, m.Err = StepLocation, ""
	m.Input.Placeholder = d.Placeholder
	m.Input.SetValue("")
	m.Input.Width = m.inputWidth()
	m.refreshHints()
	return m.Input.Focus()
}
