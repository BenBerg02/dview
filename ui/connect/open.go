package connect

func (m *Model) open() {
	conn, err := m.Drivers[m.Cur].Open(m.Input.Value())
	if err != nil {
		m.Err = err.Error()
		return
	}
	m.Conn = conn
}
