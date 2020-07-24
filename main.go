package main

import (
	"os"
	"reflect"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func messagesHandler() {
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
		if update.Message.Photo != nil {
			imageURL := (*update.Message.Photo)[0].FileID
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, imageUrl)
			bot.Send(msg)
			return
		}

		if update.Message == nil || !strings.HasPrefix(update.Message.Text, "/") {
			continue
		}

		//Check if message from user is text
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			command := ParseCommand(update.Message.Text)
			switch command.name {
			case "chat_id":
				ChatIDHandler(bot, update.Message.Chat.ID, command)
			case "rank":
				RankHandler(bot, update.Message.Chat.ID, command)
			case "stats":
				StatsHandler(bot, update.Message.Chat.ID, command)
			case "subscribe":
				SubscribeHandler(bot, update.Message.Chat.ID, command)
			case "unsubscribe":
				UnsubscribeHandler(bot, update.Message.Chat.ID, command)
			case "help":
				HelpHandler(bot, update.Message.Chat.ID, command)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Меня писал очень плохой программист, и он не рассказал мне что значит эта комманда")
				bot.Send(msg)
			}
		} else {
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
	messagesHandler()
}
