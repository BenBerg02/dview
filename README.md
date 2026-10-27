# dview

A terminal UI for browsing and editing database tables. Pick a table, page
through its rows, filter them with a `WHERE` clause, and edit, insert or
delete rows without writing SQL.

Built in Go with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Supported databases

| Database   | Driver                                              | Notes                                                                                   |
|------------|-----------------------------------------------------|-----------------------------------------------------------------------------------------|
| SQLite     | [`modernc.org/sqlite`](https://modernc.org/sqlite)  | Pure Go, no CGO or system library needed. Opens an existing database file.              |
| PostgreSQL | [`pgx`](https://github.com/jackc/pgx)               | Lists tables from every non-system schema; tables outside `public` show as `schema.table`. |

Rows are identified by their primary key. Tables without one fall back to
`rowid` in SQLite and `ctid` in PostgreSQL. Binary (BLOB/bytea) values are
shown as a placeholder and can't be edited.

## Usage

```
dview                       # open the connection screen
dview path/to/file.sqlite   # open a SQLite database directly
dview --psql                # open the PostgreSQL connection form
dview --psql <url>          # connect to PostgreSQL directly
```

The PostgreSQL connection string can be a URL or a `key=value` string:

```
dview --psql postgres://user:password@localhost:5432/mydb?sslmode=disable
dview --psql "host=localhost user=postgres dbname=mydb"
```

### Connection screen

Choose a database type with `↑`/`↓` and `Enter` (or click it).

- **SQLite**: type the path to the database file. `Tab` completes paths and
  `~/` is expanded to your home directory.
- **PostgreSQL**: fill in host, port, user, password, database and SSL mode,
  then choose one of:
  - **Test**: check that the connection works.
  - **Connect**: connect without saving anything.
  - **Save and connect**: remember the connection and prefill the form with
    it next time.

Saved connections are stored in `~/.dview/dview.db`. Passwords are never
written to that file. They go to the system keyring (Secret Service on Linux,
Keychain on macOS, Credential Manager on Windows).

### Main screen

The screen has a table list on the left, a toolbar (`Prev`, `Next`,
`Structure`, `Search`, `Insert`) above the data grid, and a `Quit` button at
the top. `Esc` moves focus outward: data grid → toolbar → table list → Quit.
The footer always shows the keys available at that moment.

| Area        | Keys                                                                                                                 |
|-------------|----------------------------------------------------------------------------------------------------------------------|
| Table list  | `↑`/`↓` select a table · `→`/`Enter` go to its data · `Esc` go to Quit                                                |
| Data grid   | arrows move · `Enter` edit cell · `Space` cell menu · `Del` delete row · `PgUp`/`PgDn` change page · `Esc` toolbar |
| Toolbar     | `←`/`→` choose · `Enter` activate · `↓` back to data · `Esc` table list                                             |
| Anywhere    | `Ctrl+C` quit                                                                                                        |

Toolbar actions:

- **Prev / Next**: move between pages. The page size follows the terminal height.
- **Structure**: show the table's columns, types, nullability, defaults and keys.
- **Search**: filter rows with a raw SQL condition, written without the `WHERE`
  keyword, e.g. `age > 30 AND name LIKE 'A%'`. `Tab` completes column names.
  Submit an empty condition to clear the filter.
- **Insert**: open a form with one field per column. Empty fields get the
  column default or `NULL`, and required columns are marked. `Tab`/`↑`/`↓` move,
  `Ctrl+S` inserts, `Esc` cancels.

The cell menu (`Space`) has **Edit**, **Set NULL** and **Delete row**. Deleting
always asks for confirmation (`y`/`n`).

## Building

Requires Go 1.26 or newer. All drivers are pure Go, so CGO is not needed.

```sh
make build            # build/dview for the current platform
make run DB=file.db   # build and open a SQLite file
make test             # run the tests
make vet              # go vet
make clean            # remove build/
```

Static Linux binaries for other architectures:

```sh
make build-linux-amd64   # build/dview-linux-amd64
make build-linux-arm64   # build/dview-linux-arm64
```

Without `make`:

```sh
go build -o build/dview ./cmd/dview
```

### Tests

`make test` runs the unit tests. The PostgreSQL integration test runs only
when `DVIEW_TEST_PG` points at a server it can use. It creates and drops the
tables `dview_pk` and `dview_nopk`.

```sh
DVIEW_TEST_PG=postgres://user:pass@localhost/testdb go test ./db/
```
