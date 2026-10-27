package model

func (m *Model) fixColOff() {
	if len(m.Cols) == 0 {
		m.ColOff = 0
		return
	}
	ws := m.ColWidths()
	if m.Col < m.ColOff {
		m.ColOff = m.Col
	}
	for m.ColOff < m.Col {
		sum := 0
		for i := m.ColOff; i <= m.Col; i++ {
			sum += ws[i]
		}
		if sum+3*(m.Col-m.ColOff) <= m.RightInner() {
			break
		}
		m.ColOff++
	}
}
