package model

func (m *Model) changePage(d int) {
	np := m.Page + d
	if np < 0 || np >= m.PageCount() {
		return
	}
	oldPage, oldRow := m.Page, m.Row
	m.Page, m.Row = np, 0
	if err := m.load(); err != nil {
		m.Page, m.Row = oldPage, oldRow
		m.fail(err)
	}
}
