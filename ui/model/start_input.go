package model

import tea "github.com/charmbracelet/bubbletea"

// startInput opens the shared text input in the given mode.
func (m *Model) startInput(mode Mode, placeholder, value string) tea.Cmd {
	m.Input.Placeholder = placeholder
	m.Input.SetValue(value)
	m.Input.CursorEnd()
	m.Input.Width = max(20, min(60, m.W-14))
	m.Mode = mode
	return m.Input.Focus()
}
