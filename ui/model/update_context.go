package model

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) updateContext(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "left":
		m.CtxCur = max(0, m.CtxCur-1)
	case "right":
		m.CtxCur = min(len(CtxItems)-1, m.CtxCur+1)
	case "down":
		m.Focus = FocusData
	case "esc":
		m.Focus = FocusTables
	case "enter":
		switch m.CtxCur {
		case 0:
			m.changePage(-1)
		case 1:
			m.changePage(1)
		case 2:
			if m.Table != "" {
				m.Mode = ModeStructure
			}
		case 3:
			if m.Table != "" {
				return m.startSearch()
			}
		case 4:
			if m.Table != "" {
				return m.startInsert()
			}
		}
	}
	return nil
}
