package connect

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"

	"dview/ui/style"
)

// formButtons are the actions below the fields, in focus order.
var formButtons = []string{"Test", "Connect", "Save and connect"}

const (
	actionTest = iota
	actionConnect
	actionSave
)

func (m *Model) startForm() tea.Cmd {
	f := m.Drivers[m.Cur].Form
	var vals []string
	if f.Prefill != nil {
		vals = f.Prefill()
	}
	m.Step, m.Err, m.Info, m.Busy = StepForm, "", "", false
	m.Fields = make([]textinput.Model, len(f.Fields))
	for i, fd := range f.Fields {
		in := textinput.New()
		in.Prompt = "› "
		in.Placeholder = fd.Placeholder
		if fd.Secret {
			in.EchoMode = textinput.EchoPassword
		}
		if i < len(vals) {
			in.SetValue(vals[i])
		}
		m.Fields[i] = in
	}
	m.layoutForm()
	return m.focusForm(0)
}

// focusForm moves the focus, blurring the inputs it leaves.
func (m *Model) focusForm(i int) tea.Cmd {
	m.Focus = i
	var cmd tea.Cmd
	for j := range m.Fields {
		if j == i {
			cmd = m.Fields[j].Focus()
		} else {
			m.Fields[j].Blur()
		}
	}
	return cmd
}

func (m *Model) labelWidth() int {
	w := 0
	for _, fd := range m.Drivers[m.Cur].Form.Fields {
		w = max(w, ansi.StringWidth(fd.Label))
	}
	return w
}

func (m *Model) fieldWidth() int { return max(10, min(34, m.W-14-m.labelWidth()-1)) }

// layoutForm sizes the inputs to the terminal.
func (m *Model) layoutForm() {
	if m.Step != StepForm {
		return
	}
	for i := range m.Fields {
		m.Fields[i].Width = m.fieldWidth()
	}
}

func (m *Model) formValues() []string {
	vals := make([]string, len(m.Fields))
	for i, in := range m.Fields {
		vals[i] = in.Value()
	}
	return vals
}

// buttonsLine renders the action row and returns the x offset (inside the
// box content) and width of each button.
func (m *Model) buttonsLine() (string, []int, []int) {
	var parts []string
	var xs, ws []int
	x := 0
	for i, b := range formButtons {
		st := style.Button
		if m.Focus == len(m.Fields)+i {
			st = style.ActiveBtn
		}
		s := st.Render(" " + b + " ")
		xs, ws = append(xs, x), append(ws, ansi.StringWidth(s))
		x += ansi.StringWidth(s) + 1
		parts = append(parts, s)
	}
	return strings.Join(parts, " "), xs, ws
}

func (m *Model) formContent() string {
	d := m.Drivers[m.Cur]
	lw := m.labelWidth()
	lines := []string{style.Header.Render(d.Name), ""}
	for i, fd := range d.Form.Fields {
		lines = append(lines, style.Dim.Render(style.Pad(fd.Label, lw))+" "+m.Fields[i].View())
	}
	btns, _, _ := m.buttonsLine()
	lines = append(lines, "", btns, "")

	// the status line is always present so the dialog does not jump
	w := lw + 1 + 2 + m.fieldWidth() + 1
	switch {
	case m.Busy:
		lines = append(lines, style.Dim.Render("Connecting…"))
	case m.Err != "":
		lines = append(lines, style.Err.Render(ansi.Truncate(m.Err, w, "…")))
	case m.Info != "":
		lines = append(lines, style.OK.Render(ansi.Truncate(m.Info, w, "…")))
	default:
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}
