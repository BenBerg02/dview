package db

import (
	"os"
	"testing"
)

// Runs against a live server: DVIEW_TEST_PG=postgres://user:pass@host/db
func TestPostgresLive(t *testing.T) {
	dsn := os.Getenv("DVIEW_TEST_PG")
	if dsn == "" {
		t.Skip("DVIEW_TEST_PG not set")
	}
	p, err := OpenPostgres(dsn, "")
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	for _, q := range []string{
		`DROP TABLE IF EXISTS dview_pk, dview_nopk`,
		`CREATE TABLE dview_pk (id int PRIMARY KEY, name text, n int)`,
		`CREATE TABLE dview_nopk (name text)`,
		`INSERT INTO dview_pk VALUES (1, 'a', 10), (2, 'b', 20)`,
		`INSERT INTO dview_nopk VALUES ('x'), ('y')`,
	} {
		if _, err := p.conn.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	tables, err := p.Tables()
	if err != nil || len(tables) < 2 {
		t.Fatalf("tables %v %v", tables, err)
	}
	cols, err := p.Columns("dview_pk")
	if err != nil || len(cols) != 3 || !cols[0].PrimaryKey {
		t.Fatalf("columns %+v %v", cols, err)
	}
	for _, tb := range []string{"dview_pk", "dview_nopk"} {
		rows, err := p.Rows(Query{Table: tb, Limit: 10})
		if err != nil || len(rows) != 2 {
			t.Fatalf("%s rows %v %v", tb, rows, err)
		}
		last := len(rows[0].Values) - 1
		if tb == "dview_pk" {
			if err := p.UpdateCell(tb, rows[0].Ref, "n", Value{Text: "99"}); err != nil {
				t.Fatal(err)
			}
			if err := p.UpdateCell(tb, rows[0].Ref, "name", Value{Null: true}); err != nil {
				t.Fatal(err)
			}
		} else if err := p.UpdateCell(tb, rows[0].Ref, "name", Value{Text: "z"}); err != nil {
			t.Fatal(err)
		}
		_ = last
		if err := p.DeleteRow(tb, rows[1].Ref); err != nil {
			t.Fatal(err)
		}
		n, err := p.Count(tb, "")
		if err != nil || n != 1 {
			t.Fatalf("%s count %d %v", tb, n, err)
		}
	}
	rows, _ := p.Rows(Query{Table: "dview_pk", Limit: 10, Where: "n = 99"})
	if len(rows) != 1 || !rows[0].Values[1].Null || rows[0].Values[2].Text != "99" {
		t.Fatalf("after update %+v", rows)
	}

	if _, err := p.conn.Exec(`CREATE TABLE dview_ins (id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		name text NOT NULL, role text NOT NULL DEFAULT 'user', note text)`); err != nil {
		t.Fatal(err)
	}
	cols, err = p.Columns("dview_ins")
	if err != nil || !cols[0].Auto || !cols[0].AutoFill() || cols[0].Required() ||
		!cols[1].Required() || !cols[2].AutoFill() || cols[3].Required() {
		t.Fatalf("insert columns %+v %v", cols, err)
	}
	if err := p.InsertRow("dview_ins", []string{"name"}, []Value{{Text: "ana"}}); err != nil {
		t.Fatal(err)
	}
	if err := p.InsertRow("dview_ins", []string{"name", "note"}, []Value{{Text: "bo"}, {Text: "hi"}}); err != nil {
		t.Fatal(err)
	}
	if err := p.InsertRow("dview_ins", nil, nil); err == nil {
		t.Fatal("expected NOT NULL violation for defaults-only insert")
	}
	rows, _ = p.Rows(Query{Table: "dview_ins", Limit: 10})
	if len(rows) != 2 || rows[0].Values[2].Text != "user" || rows[1].Values[3].Text != "hi" {
		t.Fatalf("inserted %+v", rows)
	}
	p.conn.Exec(`DROP TABLE dview_pk, dview_nopk, dview_ins`)
}
