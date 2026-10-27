package connect

import (
	"github.com/charmbracelet/bubbles/textinput"

	"dview/db"
)

func New() *Model {
	in := textinput.New()
	in.Prompt = "› "
	return &Model{Drivers: db.Drivers, Input: in}
}

// NewAt starts directly on the form or location input of the driver selected
// by a command line flag; unknown flags fall back to the driver list.
func NewAt(flag string) *Model {
	m := New()
	for i, d := range m.Drivers {
		if d.Flag == flag {
			m.Cur = i
			m.startLocation()
		}
	}
	return m
}
