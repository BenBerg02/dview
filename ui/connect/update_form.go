package connect

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"dview/db"
)

// formDone is the result of a test or connection attempt.
type formDone struct {
	action int
	conn   db.Database
	err    error
}

func (m *Model) updateForm(msg tea.KeyMsg) tea.Cmd {
	if m.Busy {
		return nil
	}
	n := len(m.Fields)
	total := n + len(formButtons)
	switch msg.String() {
	case "esc":
		m.Step, m.Err, m.Info = StepPick, "", ""
		return nil
	case "tab", "down":
		return m.focusForm((m.Focus + 1) % total)
	case "shift+tab", "up":
		return m.focusForm((m.Focus + total - 1) % total)
	case "left":
		if m.Focus > n {
			return m.focusForm(m.Focus - 1)
		}
	case "right":
		if m.Focus >= n && m.Focus < total-1 {
			return m.focusForm(m.Focus + 1)
		}
	case "enter":
		if m.Focus < n {
			return m.focusForm(m.Focus + 1)
		}
		return m.runForm(m.Focus - n)
	}
	if m.Focus >= n {
		return nil
	}
	var cmd tea.Cmd
	m.Fields[m.Focus], cmd = m.Fields[m.Focus].Update(msg)
	return cmd
}

func (m *Model) mouseForm(msg tea.MouseMsg) tea.Cmd {
	if m.Busy || msg.Action != tea.MouseActionPress || msg.Button != tea.MouseButtonLeft {
		return nil
	}
	_, ox, oy := m.box()
	line, x := msg.Y-oy-1, msg.X-ox-3
	n := len(m.Fields)
	switch {
	case line >= itemsStart && line < itemsStart+n:
		return m.focusForm(line - itemsStart)
	case line == itemsStart+n+1:
		_, xs, ws := m.buttonsLine()
		for i := range formButtons {
			if x >= xs[i] && x < xs[i]+ws[i] {
				m.focusForm(n + i)
				return m.runForm(i)
			}
		}
	}
	return nil
}

// runForm starts the chosen action in the background.
func (m *Model) runForm(action int) tea.Cmd {
	f := m.Drivers[m.Cur].Form
	vals := m.formValues()
	m.Busy, m.Err, m.Info = true, "", ""
	return func() tea.Msg {
		if action == actionTest {
			return formDone{action: action, err: f.Test(vals)}
		}
		conn, err := f.Open(vals)
		if err != nil {
			return formDone{action: action, err: err}
		}
		if action == actionSave {
			if err := f.Save(vals); err != nil {
				conn.Close()
				return formDone{action: action, err: fmt.Errorf("connected, but could not save: %w", err)}
			}
		}
		return formDone{action: action, conn: conn}
	}
}

func (m *Model) finishForm(r formDone) tea.Cmd {
	m.Busy = false
	switch {
	case r.err != nil:
		m.Err = r.err.Error()
	case r.action == actionTest:
		m.Info = "✓ Connection successful"
	default:
		m.Conn = r.conn
	}
	return nil
}
