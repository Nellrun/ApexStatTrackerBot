package main

import (
	db "github.com/heroku/go-apex-tracker/postgresql-db"
)

func main() {
	db.CreateTables()
}
