package model

func (m *Model) reload() {
	if err := m.load(); err != nil {
		m.fail(err)
	}
}
