package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// SQLite is the Database implementation backed by a SQLite file.
type SQLite struct {
	conn     *sql.DB
	path     string
	hasRowid map[string]bool
}

// sqliteRef identifies a row through its key columns (rowid or primary key).
type sqliteRef struct {
	keys []string
	vals []any
}

func OpenSQLite(path string) (*SQLite, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(1)
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}
	return &SQLite{conn: conn, path: path, hasRowid: map[string]bool{}}, nil
}

func (s *SQLite) Name() string { return s.path }

func (s *SQLite) Close() error { return s.conn.Close() }

func quote(id string) string { return `"` + strings.ReplaceAll(id, `"`, `""`) + `"` }

func (s *SQLite) Tables() ([]string, error) {
	rs, err := s.conn.Query(`SELECT name FROM sqlite_master
		WHERE type = 'table' AND substr(name, 1, 7) <> 'sqlite_' ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var out []string
	for rs.Next() {
		var n string
		if err := rs.Scan(&n); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rs.Err()
}

func (s *SQLite) Columns(table string) ([]Column, error) {
	rs, err := s.conn.Query("PRAGMA table_info(" + quote(table) + ")")
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var out []Column
	for rs.Next() {
		var (
			cid, notNull, pk int
			name, typ        string
			def              sql.NullString
		)
		if err := rs.Scan(&cid, &name, &typ, &notNull, &def, &pk); err != nil {
			return nil, err
		}
		out = append(out, Column{
			Name: name, Type: typ, NotNull: notNull != 0, PrimaryKey: pk > 0,
			Default: def.String, HasDefault: def.Valid,
		})
	}
	pks := 0
	for _, c := range out {
		if c.PrimaryKey {
			pks++
		}
	}
	for i := range out {
		// a lone INTEGER primary key is an alias of the rowid and fills itself
		out[i].Auto = pks == 1 && out[i].PrimaryKey && strings.EqualFold(out[i].Type, "INTEGER")
	}
	return out, rs.Err()
}

func (s *SQLite) Count(table, where string) (int, error) {
	q := "SELECT COUNT(*) FROM " + quote(table)
	if strings.TrimSpace(where) != "" {
		q += " WHERE (" + where + ")"
	}
	var n int
	err := s.conn.QueryRow(q).Scan(&n)
	return n, err
}

// keyColumns returns the columns that uniquely identify a row.
func (s *SQLite) keyColumns(table string) ([]string, error) {
	has, ok := s.hasRowid[table]
	if !ok {
		rs, err := s.conn.Query("SELECT rowid FROM " + quote(table) + " LIMIT 0")
		has = err == nil
		if rs != nil {
			rs.Close()
		}
		s.hasRowid[table] = has
	}
	if has {
		return []string{"rowid"}, nil
	}
	cols, err := s.Columns(table)
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, c := range cols {
		if c.PrimaryKey {
			keys = append(keys, c.Name)
		}
	}
	return keys, nil
}

func (s *SQLite) Rows(q Query) ([]Row, error) {
	keys, err := s.keyColumns(q.Table)
	if err != nil {
		return nil, err
	}
	sel := "*"
	order := ""
	if len(keys) > 0 {
		qk := make([]string, len(keys))
		for i, k := range keys {
			qk[i] = quote(k)
		}
		sel = strings.Join(qk, ", ") + ", *"
		order = " ORDER BY " + strings.Join(qk, ", ")
	}
	query := "SELECT " + sel + " FROM " + quote(q.Table)
	if strings.TrimSpace(q.Where) != "" {
		query += " WHERE (" + q.Where + ")"
	}
	query += order + " LIMIT ? OFFSET ?"

	rs, err := s.conn.Query(query, q.Limit, q.Offset)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	names, err := rs.Columns()
	if err != nil {
		return nil, err
	}
	nk := len(keys)
	var out []Row
	for rs.Next() {
		raw := make([]any, len(names))
		ptrs := make([]any, len(names))
		for i := range raw {
			ptrs[i] = &raw[i]
		}
		if err := rs.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := Row{Values: make([]Value, 0, len(names)-nk)}
		if nk > 0 {
			row.Ref = sqliteRef{keys: keys, vals: raw[:nk:nk]}
		}
		for _, v := range raw[nk:] {
			row.Values = append(row.Values, toValue(v))
		}
		out = append(out, row)
	}
	return out, rs.Err()
}

func toValue(v any) Value {
	switch x := v.(type) {
	case nil:
		return Value{Null: true}
	case []byte:
		return Value{Blob: true, Text: fmt.Sprintf("<blob %d bytes>", len(x))}
	case string:
		return Value{Text: x}
	case int64:
		return Value{Text: strconv.FormatInt(x, 10)}
	case float64:
		return Value{Text: strconv.FormatFloat(x, 'g', -1, 64)}
	case bool:
		return Value{Text: strconv.FormatBool(x)}
	case time.Time:
		return Value{Text: x.Format(time.RFC3339)}
	case map[string]any, []any:
		if b, err := json.Marshal(x); err == nil {
			return Value{Text: string(b)}
		}
		return Value{Text: fmt.Sprint(x)}
	default:
		return Value{Text: fmt.Sprint(x)}
	}
}

func (r sqliteRef) where() (string, []any) {
	parts := make([]string, len(r.keys))
	for i, k := range r.keys {
		parts[i] = quote(k) + " = ?"
	}
	return strings.Join(parts, " AND "), r.vals
}

func asRef(ref any) (sqliteRef, error) {
	r, ok := ref.(sqliteRef)
	if !ok || len(r.keys) == 0 {
		return r, errors.New("row has no usable identity (no rowid or primary key)")
	}
	return r, nil
}

func (s *SQLite) exec1(query string, args ...any) error {
	res, err := s.conn.Exec(query, args...)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", n)
	}
	return nil
}

func (s *SQLite) UpdateCell(table string, ref any, column string, v Value) error {
	r, err := asRef(ref)
	if err != nil {
		return err
	}
	where, args := r.where()
	var val any
	if !v.Null {
		val = v.Text
	}
	return s.exec1("UPDATE "+quote(table)+" SET "+quote(column)+" = ? WHERE "+where,
		append([]any{val}, args...)...)
}

func (s *SQLite) InsertRow(table string, cols []string, vals []Value) error {
	query, args := insertQuery(quote(table), cols, vals, func(int) string { return "?" })
	return s.exec1(query, args...)
}

func (s *SQLite) DeleteRow(table string, ref any) error {
	r, err := asRef(ref)
	if err != nil {
		return err
	}
	where, args := r.where()
	return s.exec1("DELETE FROM "+quote(table)+" WHERE "+where, args...)
}
