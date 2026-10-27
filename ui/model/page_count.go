package model

func (m *Model) PageCount() int {
	return max(1, (m.Total+m.PageSize-1)/m.PageSize)
}
