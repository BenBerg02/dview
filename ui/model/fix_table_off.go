package model

func (m *Model) fixTableOff() {
	h := m.PaneHeight()
	if m.TableCur < m.TableOff {
		m.TableOff = m.TableCur
	}
	if m.TableCur >= m.TableOff+h {
		m.TableOff = m.TableCur - h + 1
	}
}
