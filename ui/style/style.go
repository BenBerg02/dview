// Package style holds the palette and drawing helpers shared by the screens.
package style

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var (
	Accent    = lipgloss.Color("62")
	Dim       = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	Err       = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	OK        = lipgloss.NewStyle().Foreground(lipgloss.Color("114"))
	Header    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	Button    = lipgloss.NewStyle().Background(lipgloss.Color("237"))
	ActiveBtn = lipgloss.NewStyle().Background(Accent).Foreground(lipgloss.Color("230")).Bold(true)
	RowSel    = lipgloss.NewStyle().Background(lipgloss.Color("238"))
	CellSel   = lipgloss.NewStyle().Background(Accent).Foreground(lipgloss.Color("230"))
	Box       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(Accent).Padding(0, 2)
)

// Pad truncates s to w cells and right-pads it with spaces to exactly w.
func Pad(s string, w int) string {
	s = ansi.Truncate(s, w, "…")
	if n := w - ansi.StringWidth(s); n > 0 {
		s += strings.Repeat(" ", n)
	}
	return s
}

// Footer renders a one-line status text in the dim style, indented by one
// cell and cut to w.
func Footer(s string, w int) string {
	return Dim.Render(ansi.Truncate(" "+s, w, "…"))
}
