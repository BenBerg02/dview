package model

func (m *Model) info(s string) {
	m.Status, m.StatusErr = s, false
}
