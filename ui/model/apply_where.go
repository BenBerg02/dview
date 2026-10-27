package model

import "strings"

func (m *Model) applyWhere(where string) {
	oldWhere, oldPage, oldRow := m.Where, m.Page, m.Row
	m.Where, m.Page, m.Row = strings.TrimSpace(where), 0, 0
	if err := m.load(); err != nil {
		m.Where, m.Page, m.Row = oldWhere, oldPage, oldRow
		m.fail(err)
		return
	}
	if m.Where == "" {
		m.info("Filter cleared")
	} else {
		m.info("Filter applied")
	}
}
