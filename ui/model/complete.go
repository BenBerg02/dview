package model

import (
	"strings"
	"unicode"
)

// whereKeywords are offered next to the column names while typing a filter.
var whereKeywords = []string{
	"AND", "OR", "NOT", "LIKE", "ILIKE", "IN", "IS", "NULL", "BETWEEN",
	"TRUE", "FALSE", "EXISTS", "ESCAPE",
}

func isWordRune(r rune) bool { return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) }

// quoteIfNeeded wraps names that are not plain lowercase identifiers.
func quoteIfNeeded(name string) string {
	for _, r := range name {
		if r != '_' && !unicode.IsDigit(r) && !(unicode.IsLetter(r) && unicode.IsLower(r)) {
			return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
		}
	}
	return name
}

// whereWord returns the word before the cursor of the search input and where
// it starts, in runes. It is empty inside a string literal.
func (m *Model) whereWord() (rs []rune, start int) {
	rs = []rune(m.Input.Value())
	end := min(m.Input.Position(), len(rs))
	quotes := 0
	for _, r := range rs[:end] {
		if r == '\'' {
			quotes++
		}
	}
	if quotes%2 == 1 {
		return rs, end
	}
	start = end
	for start > 0 && isWordRune(rs[start-1]) {
		start--
	}
	return rs, start
}

// WhereHints lists the column names and keywords that complete the word
// being typed in the search input.
func (m *Model) WhereHints() []string {
	if m.Mode != ModeSearch {
		return nil
	}
	rs, start := m.whereWord()
	end := min(m.Input.Position(), len(rs))
	word := strings.ToLower(string(rs[start:end]))
	if word == "" {
		return nil
	}
	var out []string
	for _, c := range m.Cols {
		if strings.HasPrefix(strings.ToLower(c.Name), word) {
			out = append(out, quoteIfNeeded(c.Name))
		}
	}
	for _, k := range whereKeywords {
		if strings.HasPrefix(strings.ToLower(k), word) {
			out = append(out, k)
		}
	}
	return out
}

// completeWhere behaves like a shell: a single candidate is completed in
// full, several extend the word to their common prefix.
func (m *Model) completeWhere() {
	hints := m.WhereHints()
	if len(hints) == 0 {
		return
	}
	rs, start := m.whereWord()
	end := min(m.Input.Position(), len(rs))
	word := string(rs[start:end])

	completion := hints[0]
	if len(hints) > 1 {
		completion = commonPrefixFold(hints)
		if len([]rune(completion)) <= len([]rune(word)) {
			return
		}
	}
	// a quoted candidate replaces the whole word, a plain one keeps its case
	rep := []rune(completion)
	m.Input.SetValue(string(rs[:start]) + string(rep) + string(rs[end:]))
	m.Input.SetCursor(start + len(rep))
}

// commonPrefixFold is the longest prefix shared by all names, ignoring case.
func commonPrefixFold(names []string) string {
	prefix := []rune(names[0])
	for _, n := range names[1:] {
		r := []rune(n)
		i := 0
		for i < len(prefix) && i < len(r) && unicode.ToLower(prefix[i]) == unicode.ToLower(r[i]) {
			i++
		}
		prefix = prefix[:i]
	}
	return string(prefix)
}
