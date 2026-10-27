package connect

func (m *Model) refreshHints() {
	m.Hints = nil
	if !m.Drivers[m.Cur].CompletePath {
		return
	}
	_, _, m.Hints = listPath(m.Input.Value())
}
