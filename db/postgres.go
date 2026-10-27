package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const pgTimeout = 10 * time.Second

// Postgres is the Database implementation backed by a PostgreSQL server.
type Postgres struct {
	conn   *sql.DB
	name   string
	tables map[string][2]string // display name -> schema, table
}

// pgRef identifies a row through its primary key, or through ctid when the
// table has none.
type pgRef struct {
	keys []string
	vals []any
	ctid bool
}

// PostgresConfig holds the parameters of a server connection.
type PostgresConfig struct {
	Host, Port, User, Password, Database, SSLMode string
}

// DSN renders the config as a connection URL.
func (c PostgresConfig) DSN() string {
	u := url.URL{Scheme: "postgres", Path: "/" + c.Database}
	u.Host = net.JoinHostPort(orDefault(c.Host, "localhost"), orDefault(c.Port, "5432"))
	switch {
	case c.User != "" && c.Password != "":
		u.User = url.UserPassword(c.User, c.Password)
	case c.User != "":
		u.User = url.User(c.User)
	}
	q := url.Values{}
	q.Set("sslmode", orDefault(c.SSLMode, "prefer"))
	q.Set("connect_timeout", fmt.Sprint(int(pgTimeout.Seconds())))
	u.RawQuery = q.Encode()
	return u.String()
}

func orDefault(s, def string) string {
	if s = strings.TrimSpace(s); s != "" {
		return s
	}
	return def
}

// Name is the display name of the connection, without the password.
func (c PostgresConfig) Name() string {
	return fmt.Sprintf("%s@%s:%s/%s", c.User, orDefault(c.Host, "localhost"), orDefault(c.Port, "5432"), c.Database)
}

// OpenPostgres connects using a URL or key=value connection string. An empty
// name is derived from the string.
func OpenPostgres(dsn, name string) (*Postgres, error) {
	if name == "" {
		cfg, err := pgconn.ParseConfig(dsn)
		if err != nil {
			return nil, err
		}
		name = fmt.Sprintf("%s@%s:%d/%s", cfg.User, cfg.Host, cfg.Port, cfg.Database)
	}
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), pgTimeout)
	defer cancel()
	if err := conn.PingContext(ctx); err != nil {
		conn.Close()
		return nil, err
	}
	return &Postgres{conn: conn, name: name, tables: map[string][2]string{}}, nil
}

// OpenPostgresConfig connects using explicit parameters.
func OpenPostgresConfig(c PostgresConfig) (*Postgres, error) {
	if strings.TrimSpace(c.Database) == "" {
		return nil, errors.New("database name is required")
	}
	return OpenPostgres(c.DSN(), c.Name())
}

// TestPostgres checks that the server accepts the connection.
func TestPostgres(c PostgresConfig) error {
	p, err := OpenPostgresConfig(c)
	if err != nil {
		return err
	}
	return p.Close()
}

func (p *Postgres) Name() string { return p.name }

func (p *Postgres) Close() error { return p.conn.Close() }

// qualified returns the quoted schema.table for a display name.
func (p *Postgres) qualified(table string) string {
	st, ok := p.tables[table]
	if !ok {
		st = [2]string{"public", table}
	}
	return quote(st[0]) + "." + quote(st[1])
}

func (p *Postgres) Tables() ([]string, error) {
	rs, err := p.conn.Query(`SELECT table_schema, table_name FROM information_schema.tables
		WHERE table_type = 'BASE TABLE'
		AND table_schema NOT IN ('pg_catalog', 'information_schema')
		AND table_schema NOT LIKE 'pg\_toast%'
		ORDER BY 1, 2`)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var out []string
	p.tables = map[string][2]string{}
	for rs.Next() {
		var schema, name string
		if err := rs.Scan(&schema, &name); err != nil {
			return nil, err
		}
		display := name
		if schema != "public" {
			display = schema + "." + name
		}
		p.tables[display] = [2]string{schema, name}
		out = append(out, display)
	}
	return out, rs.Err()
}

func (p *Postgres) Columns(table string) ([]Column, error) {
	rs, err := p.conn.Query(`SELECT a.attname, format_type(a.atttypid, a.atttypmod), a.attnotnull,
			EXISTS (SELECT 1 FROM pg_index i
				WHERE i.indrelid = a.attrelid AND i.indisprimary AND a.attnum = ANY(i.indkey)),
			pg_get_expr(d.adbin, d.adrelid), a.attidentity <> ''
		FROM pg_attribute a
		LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
		WHERE a.attrelid = $1::regclass AND a.attnum > 0 AND NOT a.attisdropped
		ORDER BY a.attnum`, p.qualified(table))
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var out []Column
	for rs.Next() {
		var (
			c   Column
			def sql.NullString
		)
		if err := rs.Scan(&c.Name, &c.Type, &c.NotNull, &c.PrimaryKey, &def, &c.Auto); err != nil {
			return nil, err
		}
		c.Default, c.HasDefault = def.String, def.Valid
		out = append(out, c)
	}
	return out, rs.Err()
}

func (p *Postgres) Count(table, where string) (int, error) {
	q := "SELECT COUNT(*) FROM " + p.qualified(table)
	if strings.TrimSpace(where) != "" {
		q += " WHERE (" + where + ")"
	}
	var n int
	err := p.conn.QueryRow(q).Scan(&n)
	return n, err
}

// keyColumns returns the primary key columns, or nil when the table has none.
func (p *Postgres) keyColumns(table string) ([]string, error) {
	cols, err := p.Columns(table)
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

func (p *Postgres) Rows(q Query) ([]Row, error) {
	keys, err := p.keyColumns(q.Table)
	if err != nil {
		return nil, err
	}
	ctid := len(keys) == 0
	var exprs []string
	if ctid {
		exprs = []string{"ctid::text"}
	} else {
		for _, k := range keys {
			exprs = append(exprs, quote(k))
		}
	}
	nk := len(exprs)
	query := "SELECT " + strings.Join(exprs, ", ") + ", * FROM " + p.qualified(q.Table)
	if strings.TrimSpace(q.Where) != "" {
		query += " WHERE (" + q.Where + ")"
	}
	query += " ORDER BY " + strings.Join(exprs, ", ") + " LIMIT $1 OFFSET $2"

	rs, err := p.conn.Query(query, q.Limit, q.Offset)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	names, err := rs.Columns()
	if err != nil {
		return nil, err
	}
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
		row := Row{Ref: pgRef{keys: keys, vals: raw[:nk:nk], ctid: ctid}, Values: make([]Value, 0, len(names)-nk)}
		for _, v := range raw[nk:] {
			row.Values = append(row.Values, toValue(v))
		}
		out = append(out, row)
	}
	return out, rs.Err()
}

// where builds the condition matching the row, numbering placeholders from
// first.
func (r pgRef) where(first int) (string, []any) {
	if r.ctid {
		return fmt.Sprintf("ctid = $%d::tid", first), r.vals
	}
	parts := make([]string, len(r.keys))
	for i, k := range r.keys {
		parts[i] = fmt.Sprintf("%s = $%d", quote(k), first+i)
	}
	return strings.Join(parts, " AND "), r.vals
}

func (p *Postgres) exec1(query string, args ...any) error {
	res, err := p.conn.Exec(query, args...)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", n)
	}
	return nil
}

func (p *Postgres) UpdateCell(table string, ref any, column string, v Value) error {
	r, ok := ref.(pgRef)
	if !ok {
		return errors.New("row has no usable identity")
	}
	where, args := r.where(2)
	var val any
	if !v.Null {
		val = v.Text
	}
	return p.exec1("UPDATE "+p.qualified(table)+" SET "+quote(column)+" = $1 WHERE "+where,
		append([]any{val}, args...)...)
}

func (p *Postgres) InsertRow(table string, cols []string, vals []Value) error {
	query, args := insertQuery(p.qualified(table), cols, vals, func(n int) string { return fmt.Sprintf("$%d", n) })
	return p.exec1(query, args...)
}

func (p *Postgres) DeleteRow(table string, ref any) error {
	r, ok := ref.(pgRef)
	if !ok {
		return errors.New("row has no usable identity")
	}
	where, args := r.where(1)
	return p.exec1("DELETE FROM "+p.qualified(table)+" WHERE "+where, args...)
}
