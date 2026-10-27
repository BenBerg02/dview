// Package connect implements the startup screen that asks for the database
// type and its location.
package connect

import (
	"github.com/charmbracelet/bubbles/textinput"

	"dview/db"
)

type Step int

const (
	StepPick Step = iota
	StepLocation
	StepForm // driver with several fields, e.g. PostgreSQL
)

type Model struct {
	W, H    int
	Drivers []db.Driver
	Cur     int
	Step    Step
	Input   textinput.Model
	Err     string
	Hints   []Entry // path completions for the current input

	// Form step state.
	Fields []textinput.Model
	Focus  int    // 0..len(Fields)-1 are inputs, then the buttons
	Info   string // last successful action message
	Busy   bool   // a test or connection attempt is running

	// Conn is set once a connection was opened successfully.
	Conn db.Database
}
