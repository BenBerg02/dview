package model

import tea "github.com/charmbracelet/bubbletea"

// updateInput handles the keys shared by the edit and search dialogs; apply
// receives the typed text on enter.
func (m *Model) updateInput(msg tea.KeyMsg, apply func(string)) tea.Cmd {
	switch msg.String() {
	case "esc":
		m.Mode = ModeNone
		return nil
	case "enter":
		m.Mode = ModeNone
		apply(m.Input.Value())
		return nil
	}
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return cmd
}
