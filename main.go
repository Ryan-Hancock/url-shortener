package main

import (
	_ "github.com/mattn/go-sqlite3"
)

const dbURI = "./url.db"

func main() {
	db := NewDB(dbURI)

	h := Handler{}
	h.initialise(db)
	h.run("127.0.0.1:8000")
}
