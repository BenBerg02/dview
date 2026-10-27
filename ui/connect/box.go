package connect

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"dview/ui/style"
)

// itemsStart is the content line of the first driver in the pick dialog.
const itemsStart = 2

func (m *Model) content() string {
	if m.Step == StepPick {
		w := 24
		for _, d := range m.Drivers {
			w = max(w, ansi.StringWidth(d.Name)+2)
		}
		lines := []string{style.Header.Render("Select database type"), ""}
		for i, d := range m.Drivers {
			item := style.Pad(" "+d.Name, w)
			if i == m.Cur {
				item = style.ActiveBtn.Render(item)
			}
			lines = append(lines, item)
		}
		return strings.Join(lines, "\n")
	}
	if m.Step == StepForm {
		return m.formContent()
	}
	d := m.Drivers[m.Cur]
	lines := []string{style.Header.Render(d.Name), style.Dim.Render(d.LocationLabel), "", m.Input.View()}
	lines = append(lines, m.hintLines()...)
	if m.Err != "" {
		lines = append(lines, "", style.Err.Render(ansi.Truncate(m.Err, m.inputWidth()+2, "…")))
	}
	return strings.Join(lines, "\n")
}

func (m *Model) boxWidth() int {
	return lipgloss.Width(style.Box.Render(m.content()))
}

// box returns the rendered dialog and its top-left screen position.
func (m *Model) box() (string, int, int) {
	b := style.Box.Render(m.content())
	x := max(0, (m.W-lipgloss.Width(b))/2)
	y := max(0, (m.H-1-lipgloss.Height(b))/2)
	return b, x, y
}

func hint(m *Model) string {
	if m.Step == StepPick {
		return "click or enter select · ↑↓ move · esc quit"
	}
	if m.Step == StepForm {
		return "tab/↑↓ move · enter select · esc back"
	}
	if m.Drivers[m.Cur].CompletePath {
		return "tab complete · enter connect · esc back"
	}
	return "enter connect · esc back"
}

// inputWidth is the width of the location input.
func (m *Model) inputWidth() int { return max(20, min(50, m.W-14)) }
