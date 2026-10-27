package db

import "strings"

// insertQuery builds an INSERT for the given columns; ph renders the
// placeholder of the n-th (1-based) argument.
func insertQuery(table string, cols []string, vals []Value, ph func(n int) string) (string, []any) {
	if len(cols) == 0 {
		return "INSERT INTO " + table + " DEFAULT VALUES", nil
	}
	names := make([]string, len(cols))
	marks := make([]string, len(cols))
	args := make([]any, len(cols))
	for i, c := range cols {
		names[i], marks[i] = quote(c), ph(i+1)
		if !vals[i].Null {
			args[i] = vals[i].Text
		}
	}
	return "INSERT INTO " + table + " (" + strings.Join(names, ", ") + ") VALUES (" + strings.Join(marks, ", ") + ")", args
}
