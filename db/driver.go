package db

import (
	"os"
	"path/filepath"
	"strings"
)

// Driver describes a selectable database backend and how to open it.
type Driver struct {
	Name          string
	LocationLabel string // what Open expects, e.g. "Database file path"
	Placeholder   string
	CompletePath  bool // location is a filesystem path: offer tab completion and hints
	Open          func(location string) (Database, error)

	// Flag is the command line switch that selects the driver (without "--").
	Flag string
	// Form, when set, replaces the single location input by several fields.
	Form *Form
}

// Field is one input of a connection form.
type Field struct {
	Label       string
	Placeholder string
	Secret      bool // masked while typing
}

// Form describes a driver that needs several parameters to connect. Values
// are passed in the order of Fields.
type Form struct {
	Fields  []Field
	Prefill func() []string // initial values, nil entries meaning empty
	Test    func(vals []string) error
	Open    func(vals []string) (Database, error)
	Save    func(vals []string) error
}

// DriverByFlag returns the driver selected by a command line switch.
func DriverByFlag(flag string) (Driver, bool) {
	for _, d := range Drivers {
		if d.Flag == flag {
			return d, true
		}
	}
	return Driver{}, false
}

// Drivers lists every backend offered by the connection screen.
var Drivers = []Driver{
	{
		Name:          "SQLite",
		LocationLabel: "Database file path",
		Placeholder:   "/path/to/file.sqlite",
		CompletePath:  true,
		Open:          openSQLiteFile,
	},
	{
		Name: "PostgreSQL",
		Flag: "psql",
		Form: &Form{
			Fields: []Field{
				{Label: "Host", Placeholder: "localhost"},
				{Label: "Port", Placeholder: "5432"},
				{Label: "User", Placeholder: "postgres"},
				{Label: "Password", Secret: true},
				{Label: "Database", Placeholder: "mydb"},
				{Label: "SSL mode", Placeholder: "prefer"},
			},
			Prefill: prefillPostgres,
			Test:    func(v []string) error { return TestPostgres(postgresConfig(v)) },
			Open: func(v []string) (Database, error) {
				p, err := OpenPostgresConfig(postgresConfig(v))
				if err != nil {
					return nil, err
				}
				return p, nil
			},
			Save: func(v []string) error { return SavePostgres(postgresConfig(v)) },
		},
	},
}

func postgresConfig(v []string) PostgresConfig {
	return PostgresConfig{Host: v[0], Port: v[1], User: v[2], Password: v[3], Database: v[4], SSLMode: v[5]}
}

func prefillPostgres() []string {
	c, ok := LastPostgres()
	if !ok {
		return nil
	}
	return []string{c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode}
}

func openSQLiteFile(path string) (Database, error) {
	path = strings.TrimSpace(path)
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			path = filepath.Join(home, path[2:])
		}
	}
	if path == "" {
		return nil, os.ErrInvalid
	}
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	return OpenSQLite(path)
}
