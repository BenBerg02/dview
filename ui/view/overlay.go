package view

import "dview/ui/model"

func overlay(m *model.Model) string {
	switch m.Mode {
	case model.ModeEdit:
		return editBox(m)
	case model.ModeSearch:
		return searchBox(m)
	case model.ModeMenu:
		return menuBox(m)
	case model.ModeInsert:
		return insertBox(m)
	case model.ModeConfirm:
		return confirmBox(m)
	case model.ModeStructure:
		return structureBox(m)
	}
	return ""
}
