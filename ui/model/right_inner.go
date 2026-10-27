package model

// RightInner is the inner width of the data pane.
func (m *Model) RightInner() int { return max(1, m.W-m.LeftWidth()-2) }
