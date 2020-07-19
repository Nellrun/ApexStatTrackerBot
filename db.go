package main;

import (
	"os"
	"database/sql"
	_ "github.com/lib/pq"
 )


func Subscribe(username string, chatId int64) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO subscriptions (username, chat_id) VALUES ($1, $2);`


	_, err = db.Exec(query, username, chatId)
	if err != nil {
		return err
	}

	return nil;
}


func CreateTables() error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS subscriptions (id SERIAL PRIMARY KEY, username TEXT, platform TEXT DEFAULT psn, chat_id INT);`)
	if err != nil {
		return err
	}

	return nil;
}