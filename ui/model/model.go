// Package model holds the UI state and the logic that changes it.
package model

import (
	"github.com/charmbracelet/bubbles/textinput"

	"dview/db"
)

type Model struct {
	DB   db.Database
	W, H int

	Tables   []string
	TableCur int
	TableOff int

	Focus Focus
	Mode  Mode

	Table    string
	Cols     []db.Column
	Rows     []db.Row
	Total    int
	Page     int
	PageSize int
	Where    string
	Row, Col int
	ColOff   int

	CtxCur  int
	MenuCur int
	Input   textinput.Model

	Inserts   []textinput.Model // one input per column while inserting
	InsertCur int               // focused input, then the Insert and Cancel buttons
	InsertOff int               // first visible input
	InsertErr string

	Status    string
	StatusErr bool
}
