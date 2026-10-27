package connect

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) updatePick(msg tea.KeyMsg) tea.Cmd {
	n := len(m.Drivers)
	switch msg.String() {
	case "up":
		m.Cur = (m.Cur + n - 1) % n
	case "down":
		m.Cur = (m.Cur + 1) % n
	case "enter":
		return m.startLocation()
	case "esc", "q":
		return tea.Quit
	}
	return nil
}
