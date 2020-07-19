package main

import (
	"os"
	"reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

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
		if update.Message == nil {
			continue
		}

		//Check if message from user is text
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			switch update.Message.Text {
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello")
				bot.Send(msg)
			}
		} else {
			//Send message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don undestand you")
			bot.Send(msg)
		}
	}
}
