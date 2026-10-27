package model

// PaneHeight is the inner height of both panes.
func (m *Model) PaneHeight() int { return max(1, m.H-4) }
