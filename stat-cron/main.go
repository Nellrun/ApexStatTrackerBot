package main

import (
	"encoding/json"
	"log"
	"time"

	tracker "github.com/heroku/go-apex-tracker/apex-tracker"
	db "github.com/heroku/go-apex-tracker/postgresql-db"
)

func main() {
	subscriptions, err := db.GetSubscriptions()
	if err != nil {
		log.Print(err)
		return
	}

	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	for _, subscription := range subscriptions {
		stats, err := tracker.GetStats(subscription.Username, subscription.Platform)
		if err != nil {
			log.Print(err)
			continue
		}
		bytes, err := json.Marshal(stats)
		if err != nil {
			log.Print(err)
			continue
		}
		err = db.AddLog(subscription.Username, bytes, today)
		if err != nil {
			log.Print(err)
			continue
		}
	}
}
