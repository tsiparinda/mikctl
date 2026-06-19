package db

import (
	"database/sql"
	"mikctl/src/config"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func OpenDB() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", config.DBPath)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	return db, nil
}
