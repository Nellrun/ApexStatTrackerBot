package main

import (
	"fmt"
	"log"
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

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://apex-stat-tracker-bot.herokuapp.com/" + bot.Token))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/" + bot.Token)
	// go http.ListenAndServeTLS("0.0.0.0", "cert.pem", "key.pem", nil)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if update.Message.Photo != nil {
			imageURL := (*update.Message.Photo)[0].FileID
			command := ParseCommand(update.Message.Caption)
			if command.name == "uploadimage" {
				if len(command.args) < 1 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "you must provide image tag")
					bot.Send(msg)
					continue
				}
				err := AddImage(strings.ToLower(command.args[0]), imageURL)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "something went wrong")
					bot.Send(msg)
					continue
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("image for tag %s updated", command.args[0]))
				bot.Send(msg)
				continue
			}
		}

		if !strings.HasPrefix(update.Message.Text, "/") {
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
			case "deleteimage":
				DeleteImageHandler(bot, update.Message.Chat.ID, command)
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
