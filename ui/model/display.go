package model

import (
	"strings"

	"dview/db"
)

var sanitizer = strings.NewReplacer("\n", " ", "\r", " ", "\t", " ")

// Display renders a cell value as a single line of text.
func Display(v db.Value) string {
	if v.Null {
		return "NULL"
	}
	return sanitizer.Replace(v.Text)
}
