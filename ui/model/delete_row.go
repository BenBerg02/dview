package model

func (m *Model) deleteRow() {
	r := m.Rows[m.Row]
	if err := m.DB.DeleteRow(m.Table, r.Ref); err != nil {
		m.fail(err)
		return
	}
	m.reload()
	m.info("Row deleted")
}
