package connect

func commonPrefix(entries []Entry) string {
	prefix := []rune(entries[0].Name)
	for _, e := range entries[1:] {
		r := []rune(e.Name)
		n := 0
		for n < len(prefix) && n < len(r) && prefix[n] == r[n] {
			n++
		}
		prefix = prefix[:n]
	}
	return string(prefix)
}
