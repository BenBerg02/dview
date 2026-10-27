package connect

// tabComplete behaves like a shell: a single match is completed in full
// (directories get a trailing "/"), several matches extend the input to
// their common prefix.
func (m *Model) tabComplete() {
	if !m.Drivers[m.Cur].CompletePath {
		return
	}
	v := m.Input.Value()
	if v == "~" {
		v = "~/"
	} else {
		dir, base, entries := listPath(v)
		switch {
		case len(entries) == 1:
			v = dir + entries[0].Name
			if entries[0].Dir {
				v += "/"
			}
		case len(entries) > 1:
			if p := commonPrefix(entries); len(p) > len(base) {
				v = dir + p
			}
		}
	}
	m.Input.SetValue(v)
	m.Input.CursorEnd()
	m.refreshHints()
}
