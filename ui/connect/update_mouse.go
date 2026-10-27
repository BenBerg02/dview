package connect

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateMouse(msg tea.MouseMsg) tea.Cmd {
	if m.Step == StepForm {
		return m.mouseForm(msg)
	}
	if m.Step != StepPick || msg.Action != tea.MouseActionPress || msg.Button != tea.MouseButtonLeft {
		return nil
	}
	_, ox, oy := m.box()
	for i := range m.Drivers {
		row := oy + 1 + itemsStart + i
		if msg.Y == row && msg.X > ox && msg.X < ox+m.boxWidth()-1 {
			m.Cur = i
			return m.startLocation()
		}
	}
	return nil
}
