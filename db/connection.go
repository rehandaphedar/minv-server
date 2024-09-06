package db

import (
	"context"
	"database/sql"
	"log"

	_ "embed"

	"git.sr.ht/~rehandaphedar/minv-server/sqlc"
	_ "github.com/mattn/go-sqlite3"
)

var Queries *sqlc.Queries

//go:embed schema/schema.sql
var initQuery string

func Connect() {

	db, err := sql.Open("sqlite3", "file:data/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// create tables
	if _, err := db.ExecContext(context.Background(), initQuery); err != nil {
		log.Fatal(err)
	}

	Queries = sqlc.New(db)
}
