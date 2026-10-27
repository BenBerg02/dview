package model

import (
	"github.com/charmbracelet/bubbles/textinput"

	"dview/db"
)

func New(d db.Database) (*Model, error) {
	tables, err := d.Tables()
	if err != nil {
		return nil, err
	}
	in := textinput.New()
	in.Prompt = "› "
	m := &Model{DB: d, Tables: tables, PageSize: 20, Input: in}
	if len(tables) > 0 {
		m.openTable(0)
	}
	return m, nil
}
