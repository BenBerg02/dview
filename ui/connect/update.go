package connect

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.W, m.H = msg.Width, msg.Height
		m.layoutForm()
	case formDone:
		return m.finishForm(msg)
	case tea.MouseMsg:
		return m.updateMouse(msg)
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return tea.Quit
		}
		switch m.Step {
		case StepForm:
			return m.updateForm(msg)
		case StepLocation:
			return m.updateLocation(msg)
		}
		return m.updatePick(msg)
	}
	return nil
}
