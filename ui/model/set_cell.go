package model

import "dview/db"

func (m *Model) setCell(v db.Value) {
	r := m.Rows[m.Row]
	if err := m.DB.UpdateCell(m.Table, r.Ref, m.Cols[m.Col].Name, v); err != nil {
		m.fail(err)
		return
	}
	m.reload()
	m.info("Cell updated")
}
