package model

import (
	tea "github.com/charmbracelet/bubbletea"

	"dview/db"
)

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.W, m.H = msg.Width, msg.Height
		m.resize()
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return tea.Quit
		}
		m.Status = ""
		switch m.Mode {
		case ModeEdit:
			return m.updateInput(msg, func(v string) { m.setCell(db.Value{Text: v}) })
		case ModeSearch:
			return m.updateSearch(msg)
		case ModeMenu:
			return m.updateMenu(msg)
		case ModeStructure:
			m.updateStructure(msg)
		case ModeInsert:
			return m.updateInsert(msg)
		case ModeConfirm:
			m.updateConfirm(msg)
		default:
			return m.updateFocus(msg)
		}
	}
	return nil
}
