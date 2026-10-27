package model

func (m *Model) fail(err error) {
	m.Status, m.StatusErr = err.Error(), true
}
