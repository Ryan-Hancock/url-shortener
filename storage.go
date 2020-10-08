package main

import "github.com/jmoiron/sqlx"

const url_db = "./url.db"

// Storage interface
type Storage interface {
	Set(key string, value string) error
	Get(key string) string
}

type urlDB struct {
	db *sqlx.DB
}

schema := `CREATE TABLE urls (
    key varchar(10) NOT NULL PRIMARY KEY,
	url text NOT NULL);`
	
func openDB() *sqlx.DB {
	return sqlx.MustConnect("sqlite3", url_db)
}

func NewDB() *urlDB {
	return &urlDB {
		db: openDB()
	}
}

func (u urlDB) Set(key string, value string) error {
	return nil
}

func (u urlDB) Get(key string) {
	return
}
