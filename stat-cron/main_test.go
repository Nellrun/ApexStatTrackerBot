package main

import (
	"fmt"
	"os"
	"testing"

	db "github.com/heroku/go-apex-tracker/postgresql-db"
)

func TestDb(t *testing.T) {
	os.Setenv("DATABASE_URL", "TOKEN")

	usernames, err := db.GetSubscriptions()

	if err != nil {
		t.Error(err)
		return
	}

	if usernames[0] != "LUV_Nellrun" {
		t.Error(fmt.Sprintf("got %v, expected %s", usernames, "luv_nellrun"))
	}
}
