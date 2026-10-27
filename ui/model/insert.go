package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"

	"dview/db"
)

// InsertButtons are the actions below the inputs, in focus order.
var InsertButtons = []string{"Insert", "Cancel"}

func insertPlaceholder(c db.Column) string {
	switch {
	case c.Auto:
		return "auto"
	case c.HasDefault:
		return "default: " + c.Default
	case c.Required():
		return "required"
	}
	return "NULL"
}

func (m *Model) startInsert() tea.Cmd {
	if len(m.Cols) == 0 {
		return nil
	}
	m.Inserts = make([]textinput.Model, len(m.Cols))
	for i, c := range m.Cols {
		in := textinput.New()
		in.Prompt = "› "
		in.Placeholder = insertPlaceholder(c)
		m.Inserts[i] = in
	}
	m.InsertErr, m.InsertOff = "", 0
	m.Mode = ModeInsert
	m.layoutInsert()
	return m.focusInsert(0)
}

// InsertLayout returns the widths of the label, type and input columns of
// the insert dialog.
func (m *Model) InsertLayout() (lw, tw, iw int) {
	for _, c := range m.Cols {
		lw = max(lw, ansi.StringWidth(c.Name))
		tw = max(tw, ansi.StringWidth(c.Type))
	}
	lw, tw = min(lw, 24)+1, min(tw, 16) // lw has room for the required mark
	iw = max(10, min(30, m.W-14-lw-tw-2))
	return
}

// InsertVisible is the number of inputs shown at once.
func (m *Model) InsertVisible() int { return max(3, m.H-13) }

func (m *Model) layoutInsert() {
	if m.Mode != ModeInsert {
		return
	}
	_, _, iw := m.InsertLayout()
	for i := range m.Inserts {
		m.Inserts[i].Width = iw
	}
	m.fixInsertOff()
}

func (m *Model) fixInsertOff() {
	v := m.InsertVisible()
	if m.InsertCur >= len(m.Inserts) {
		return
	}
	if m.InsertCur < m.InsertOff {
		m.InsertOff = m.InsertCur
	}
	if m.InsertCur >= m.InsertOff+v {
		m.InsertOff = m.InsertCur - v + 1
	}
}

// focusInsert moves the focus, blurring the inputs it leaves.
func (m *Model) focusInsert(i int) tea.Cmd {
	m.InsertCur = i
	var cmd tea.Cmd
	for j := range m.Inserts {
		if j == i {
			cmd = m.Inserts[j].Focus()
		} else {
			m.Inserts[j].Blur()
		}
	}
	m.fixInsertOff()
	return cmd
}

func (m *Model) updateInsert(msg tea.KeyMsg) tea.Cmd {
	n := len(m.Inserts)
	total := n + len(InsertButtons)
	switch msg.String() {
	case "esc":
		m.Mode = ModeNone
		return nil
	case "ctrl+s":
		return m.insertRow()
	case "tab", "down":
		return m.focusInsert((m.InsertCur + 1) % total)
	case "shift+tab", "up":
		return m.focusInsert((m.InsertCur + total - 1) % total)
	case "left":
		if m.InsertCur > n {
			return m.focusInsert(m.InsertCur - 1)
		}
	case "right":
		if m.InsertCur >= n && m.InsertCur < total-1 {
			return m.focusInsert(m.InsertCur + 1)
		}
	case "enter":
		switch {
		case m.InsertCur < n:
			return m.focusInsert(m.InsertCur + 1)
		case m.InsertCur == n:
			return m.insertRow()
		}
		m.Mode = ModeNone
		return nil
	}
	if m.InsertCur >= n {
		return nil
	}
	var cmd tea.Cmd
	m.Inserts[m.InsertCur], cmd = m.Inserts[m.InsertCur].Update(msg)
	return cmd
}

// insertRow inserts the typed values. Empty inputs are left out so the
// database applies the default (or NULL).
func (m *Model) insertRow() tea.Cmd {
	var cols []string
	var vals []db.Value
	for i, c := range m.Cols {
		v := m.Inserts[i].Value()
		if v == "" {
			if c.Required() {
				m.InsertErr = fmt.Sprintf("%s is required", c.Name)
				return m.focusInsert(i)
			}
			continue
		}
		cols, vals = append(cols, c.Name), append(vals, db.Value{Text: v})
	}
	if err := m.DB.InsertRow(m.Table, cols, vals); err != nil {
		m.InsertErr = err.Error()
		return nil
	}
	m.Mode = ModeNone
	m.reload()
	m.info("Row inserted")
	return nil
}
