package connect

import (
	"os"
	"path/filepath"
	"strings"
)

// listPath splits input at its last "/" and lists the entries of the
// directory part whose names start with the remainder (base).
func listPath(input string) (dir, base string, entries []Entry) {
	i := strings.LastIndex(input, "/")
	dir, base = input[:i+1], input[i+1:]

	real := dir
	switch {
	case real == "":
		real = "."
	case strings.HasPrefix(real, "~/"):
		if home, err := os.UserHomeDir(); err == nil {
			real = filepath.Join(home, real[2:])
		}
	}

	items, err := os.ReadDir(real)
	if err != nil {
		return dir, base, nil
	}
	for _, it := range items {
		name := it.Name()
		if !strings.HasPrefix(name, base) || (strings.HasPrefix(name, ".") && !strings.HasPrefix(base, ".")) {
			continue
		}
		isDir := it.IsDir()
		if !isDir && it.Type()&os.ModeSymlink != 0 {
			if st, err := os.Stat(filepath.Join(real, name)); err == nil {
				isDir = st.IsDir()
			}
		}
		entries = append(entries, Entry{Name: name, Dir: isDir})
	}
	return dir, base, entries
}
