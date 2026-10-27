package view

import "dview/ui/model"

func hint(m *model.Model) string {
	switch m.Mode {
	case model.ModeEdit:
		return "enter apply · esc cancel"
	case model.ModeSearch:
		return "tab complete · enter apply · esc cancel"
	case model.ModeMenu:
		return "↑↓ choose · enter run · esc close"
	case model.ModeInsert:
		return "tab/↑↓ move · enter next/select · ctrl+s insert · esc cancel"
	case model.ModeStructure:
		return "esc close"
	case model.ModeConfirm:
		return "y confirm · n cancel"
	}
	switch m.Focus {
	case model.FocusTables:
		return "↑↓ select table · →/enter open · esc global"
	case model.FocusData:
		return "↑↓←→ move · enter edit · space menu · del delete · pgup/pgdn page · esc context"
	case model.FocusContext:
		return "←→ choose · enter activate · ↓ data · esc tables"
	case model.FocusGlobal:
		return "enter quit · esc back"
	}
	return ""
}
