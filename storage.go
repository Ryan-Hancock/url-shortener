package main

import (
	"fmt"
	"log"

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

func NewDB(uri string) urlDB {
	return urlDB{
		db: openDB(uri),
	}
}

func (u urlDB) Set(key string, value string) error {
	if key != "" || value != "" {
		return fmt.Errorf("Key/Value is empty")
	}

	_, err := u.db.NamedExec(`INSERT INTO urls (key, url) VALUES (:key,:url)`,
		map[string]interface{}{
			"key": key,
			"url": value,
		})

	if err != nil {
		log.Printf("Set error %v", err.Error())
	}

	return nil
}

func (u urlDB) Get(key string) string {
	var url string
	_ = u.db.Get(&url, "SELECT url FROM urls WHERE key=$1", key)

	return url
}
