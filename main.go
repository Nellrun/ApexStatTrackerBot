package main

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
		if update.Message == nil {
			continue
		}

		//Check if message from user is text
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			switch {
			case update.Message.Text == "/chat_id":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatInt(update.Message.Chat.ID, 10))
				bot.Send(msg)
			case strings.HasPrefix(update.Message.Text, "/subscribe"):
				args := strings.SplitAfter(update.Message.Text, "/subscribe")[1]
				clear_args := strings.TrimSpace(args)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, clear_args)
				bot.Send(msg)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello")
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
