package connect

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) updateLocation(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		m.Step, m.Err = StepPick, ""
		return nil
	case "enter":
		m.open()
		return nil
	case "tab":
		m.tabComplete()
		return nil
	}
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	m.refreshHints()
	return cmd
}
