package model

func (m *Model) calcPageSize() int { return max(1, m.H-6) }
