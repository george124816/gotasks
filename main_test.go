package main

import (
	"database/sql"
	"testing"
)

type structureTable struct {
	cid        int
	name       string
	typ        string
	notnull    int
	dflt_value bool
	pk         int
}

func TestRunMigration(t *testing.T) {
	db, err := setupDatabase()
	if err != nil {
		t.Error(err)
	}

	RunMigration(db)

	row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='tasks';")

	var tableName string
	err = row.Scan(&tableName)
	if err == sql.ErrNoRows {
		t.Fatal("tasks table was NOT created")
	} else if err != nil {
		t.Fatal(err)
	}
}

func TestRunMigration_AssertSchea(t *testing.T) {
	db, err := setupDatabase()
	if err != nil {
		t.Error(err)
	}

	RunMigration(db)

	rows, err := db.Query("PRAGMA table_info(tasks);")
	if err != nil {
		t.Fatal(err)
	}

	returnedStructure := []structureTable{}
	expectedStructure := []structureTable{
		{
			cid:     0,
			name:    "id",
			typ:     "INTEGER",
			notnull: 0,
			pk:      0,
		},
		{
			cid:  1,
			name: "name",
			typ:  "TEXT",
		},
		{
			cid:        2,
			name:       "enabled",
			typ:        "boolean",
			dflt_value: false,
		},
		{
			cid:  3,
			name: "description",
			typ:  "TEXT",
		},
	}

	for rows.Next() {
		var structure structureTable

		rows.Scan(&structure.cid, &structure.name, &structure.typ, &structure.notnull, &structure.dflt_value, &structure.pk)

		returnedStructure = append(returnedStructure, structure)
	}

	for i, expected := range expectedStructure {
		got := returnedStructure[i]
		want := expected

		if got.cid != want.cid {
			t.Fatalf("cid mismatch: got=%d want=%d", got.cid, want.cid)
		}

		if got.name != want.name {
			t.Fatalf("name mismatch: got=%s want=%s", got.name, want.name)
		}

		if got.typ != want.typ {
			t.Fatalf("type mismatch: got=%s want=%s", got.typ, want.typ)
		}
	}
}

func setupDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	return db, err
}
