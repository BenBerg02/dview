package db

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/zalando/go-keyring"
)

// keyringService namespaces dview's passwords in the system keyring.
const keyringService = "dview"

// storeDir is the per-user directory (~/.dview) that holds saved connections.
func storeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".dview"), nil
}

func openStore() (*sql.DB, error) {
	dir, err := storeDir()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	conn, err := sql.Open("sqlite", filepath.Join(dir, "dview.db"))
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(1)
	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS postgres_connections (
		host TEXT NOT NULL, port TEXT NOT NULL, user TEXT NOT NULL, dbname TEXT NOT NULL,
		sslmode TEXT NOT NULL, used_at INTEGER NOT NULL,
		PRIMARY KEY (host, port, user, dbname))`)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

// passwordKey is the keyring account of a connection.
func passwordKey(c PostgresConfig) string {
	return c.Name()
}

// SavePostgres stores the connection in ~/.dview/dview.db and its password
// in the system keyring; the password never reaches the database file.
func SavePostgres(c PostgresConfig) error {
	st, err := openStore()
	if err != nil {
		return err
	}
	defer st.Close()
	_, err = st.Exec(`INSERT INTO postgres_connections (host, port, user, dbname, sslmode, used_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (host, port, user, dbname) DO UPDATE SET sslmode = excluded.sslmode, used_at = excluded.used_at`,
		orDefault(c.Host, "localhost"), orDefault(c.Port, "5432"), c.User, c.Database,
		orDefault(c.SSLMode, "prefer"), time.Now().Unix())
	if err != nil {
		return err
	}
	if c.Password == "" {
		return nil
	}
	if err := keyring.Set(keyringService, passwordKey(c), c.Password); err != nil {
		return errors.New("password not stored in the system keyring: " + err.Error())
	}
	return nil
}

// LastPostgres returns the most recently saved connection with its password,
// or false when nothing was saved. A missing keyring entry leaves the
// password empty.
func LastPostgres() (PostgresConfig, bool) {
	var c PostgresConfig
	dir, err := storeDir()
	if err != nil {
		return c, false
	}
	if _, err := os.Stat(filepath.Join(dir, "dview.db")); err != nil {
		return c, false
	}
	st, err := openStore()
	if err != nil {
		return c, false
	}
	defer st.Close()
	err = st.QueryRow(`SELECT host, port, user, dbname, sslmode FROM postgres_connections
		ORDER BY used_at DESC LIMIT 1`).Scan(&c.Host, &c.Port, &c.User, &c.Database, &c.SSLMode)
	if err != nil {
		return PostgresConfig{}, false
	}
	c.Password, _ = keyring.Get(keyringService, passwordKey(c))
	return c, true
}
