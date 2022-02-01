package main

import (
	"fmt"
	"os"
	"testing"

	db "github.com/heroku/go-apex-tracker/postgresql-db"
)

func TestDb(t *testing.T) {
	os.Setenv("DATABASE_URL", "TOKEN")

	subscriptions, err := db.GetSubscriptions()

	if err != nil {
		t.Error(err)
		return
	}

	if subscriptions[0].Username != "LUV_Nellrun" {
		t.Error(fmt.Sprintf("got %v, expected %s", subscriptions, "luv_nellrun"))
	}
}
