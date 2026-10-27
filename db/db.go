// Package db defines the database abstraction used by the UI.
// Implementations (SQLite today, others later) only need to satisfy Database.
package db

// Column describes a table column.
type Column struct {
	Name       string
	Type       string
	NotNull    bool
	PrimaryKey bool
	Default    string
	HasDefault bool
	Auto       bool // filled in by the database itself (rowid alias, identity column)
}

// AutoFill reports whether the database provides a value when none is given.
func (c Column) AutoFill() bool { return c.HasDefault || c.Auto }

// Required reports whether an insert must provide a value for the column.
func (c Column) Required() bool { return c.NotNull && !c.AutoFill() }

// Value is a single cell value rendered as text.
type Value struct {
	Null bool
	Blob bool // binary data: displayed as a placeholder and not editable
	Text string
}

// Row is a table row. Ref is an opaque identity that the same
// implementation understands in UpdateCell and DeleteRow.
type Row struct {
	Ref    any
	Values []Value // same order as Database.Columns
}

// Query selects a page of rows from a table.
type Query struct {
	Table  string
	Where  string // raw condition without the WHERE keyword; empty means no filter
	Limit  int
	Offset int
}

// Database is the minimal contract a backend must fulfill.
type Database interface {
	// Name is a human readable identifier of the connection (e.g. a file path).
	Name() string
	Tables() ([]string, error)
	Columns(table string) ([]Column, error)
	Count(table, where string) (int, error)
	Rows(q Query) ([]Row, error)
	UpdateCell(table string, ref any, column string, v Value) error
	DeleteRow(table string, ref any) error
	// InsertRow inserts a row setting only the given columns; the others get
	// their default (or NULL). With no columns it inserts a row of defaults.
	InsertRow(table string, cols []string, vals []Value) error
	Close() error
}
