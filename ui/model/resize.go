package model

func (m *Model) resize() {
	ps := m.calcPageSize()
	if m.Table != "" {
		first := m.Page * m.PageSize
		m.PageSize = ps
		m.Page = first / ps
		m.reload()
	}
	m.PageSize = ps
	m.fixTableOff()
	m.layoutInsert()
}
