// Package view renders a model.Model as a string.
package view

import (
	"github.com/charmbracelet/lipgloss"

	"dview/ui/model"
)

func Render(m *model.Model) string {
	if m.W == 0 || m.H == 0 {
		return ""
	}
	body := panes(m)
	if box := overlay(m); box != "" {
		body = lipgloss.Place(m.W, m.H-2, lipgloss.Center, lipgloss.Center, box)
	}
	return globalBar(m) + "\n" + body + "\n" + footer(m)
}
