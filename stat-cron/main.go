package main

import (
	"encoding/json"
	"log"
	"time"

	tracker "github.com/heroku/go-apex-tracker/apex-tracker"
	db "github.com/heroku/go-apex-tracker/postgresql-db"
)

func main() {
	usernames, err := db.GetSubscriptions()
	if err != nil {
		log.Print(err)
		return
	}

	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	for _, username := range usernames {
		stats, err := tracker.GetStats(username, "psn")
		if err != nil {
			log.Print(err)
			continue
		}
		bytes, err := json.Marshal(stats)
		if err != nil {
			log.Print(err)
			continue
		}
		err = db.AddLog(username, bytes, today)
		if err != nil {
			log.Print(err)
			continue
		}
	}
}
