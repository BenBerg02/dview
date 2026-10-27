package model

func (m *Model) openTable(i int) {
	m.TableCur = i
	m.Table = m.Tables[i]
	cols, err := m.DB.Columns(m.Table)
	if err != nil {
		m.Cols, m.Rows, m.Total = nil, nil, 0
		m.fail(err)
		return
	}
	m.Cols = cols
	m.Where, m.Page, m.Row, m.Col, m.ColOff = "", 0, 0, 0, 0
	if err := m.load(); err != nil {
		m.Rows, m.Total = nil, 0
		m.fail(err)
	}
}
