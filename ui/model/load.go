package model

import "dview/db"

// load fetches the current page; on error the previous data is kept.
func (m *Model) load() error {
	ps := m.PageSize
	total, err := m.DB.Count(m.Table, m.Where)
	if err != nil {
		return err
	}
	page := m.Page
	if total == 0 {
		page = 0
	} else if last := (total - 1) / ps; page > last {
		page = last
	}
	rows, err := m.DB.Rows(db.Query{Table: m.Table, Where: m.Where, Limit: ps, Offset: page * ps})
	if err != nil {
		return err
	}
	m.Total, m.Page, m.Rows = total, page, rows
	m.Row = max(0, min(m.Row, len(rows)-1))
	m.Col = max(0, min(m.Col, len(m.Cols)-1))
	m.fixColOff()
	return nil
}
