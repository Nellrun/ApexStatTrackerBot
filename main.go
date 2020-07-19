package main

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tracker "github.com/heroku/go-apex-tracker/apex-tracker"
)

func MessagesHandler() {
	//Create bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		panic(err)
	}

	//Set update timeout
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//Get updates from bot
	updates, _ := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil || !strings.HasPrefix(update.Message.Text, "/") {
			continue
		}

		//Check if message from user is text
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			switch {
			case update.Message.Text == "/chat_id":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatInt(update.Message.Chat.ID, 10))
				bot.Send(msg)
			case strings.HasPrefix(update.Message.Text, "/rank"):
				username := strings.SplitAfter(update.Message.Text, "/rank")[1]
				clearUsername := strings.TrimSpace(username)

				segments, err := tracker.GetStats(clearUsername, "psn")
				if err != nil || len(segments) == 0 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "something went wrong, please try later")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, segments[0].Stats.RankScore.DisplayValue)
					bot.Send(msg)
				}
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Меня писал очень плохой программист, и он не рассказал мне что значит это сообщения")
				bot.Send(msg)
			}
		} else {
			//Send message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't understand you")
			bot.Send(msg)
		}
	}
}

func main() {
	err := CreateTables()
	if err != nil {
		panic(err)
	}
	MessagesHandler()
}
