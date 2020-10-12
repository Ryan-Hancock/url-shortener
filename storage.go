package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Storage interface
type Storage interface {
	Set(key string, value string) error
	Get(key string) string
}

type urlDB struct {
	db *sqlx.DB
}

var schema = `CREATE TABLE IF NOT EXISTS urls (
    key varchar(10) NOT NULL PRIMARY KEY,
	url text NOT NULL);`

func openDB(uri string) *sqlx.DB {
	db := sqlx.MustConnect("sqlite3", uri)
	db.MustExec(schema)

	return db
}

func newDB(uri string) urlDB {
	return urlDB{
		db: openDB(uri),
	}
}

func (u urlDB) Set(key string, value string) error {
	if key == "" || value == "" {
		return fmt.Errorf("Key || Value is empty")
	}

	res, err := u.db.NamedExec(`INSERT INTO urls (key, url) VALUES (:key,:url)`,
		map[string]interface{}{
			"key": key,
			"url": value,
		})
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (u urlDB) Get(key string) string {
	row := u.db.QueryRow("SELECT url FROM urls WHERE `key`=$1", key)
	var url string
	row.Scan(&url)

	return url
}
