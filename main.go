package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tracker "github.com/heroku/go-apex-tracker/apex-tracker"
)

func formatUserInfo(stats tracker.Stats) string {
	return fmt.Sprintf("Kills: %s\nDamage: %s\nRank RP: %s", stats.Kills.DisplayValue, stats.Damage.DisplayValue, stats.RankScore.DisplayValue)
}

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
		if update.Message == nil || !strings.HasPrefix(update.Message.Text, "/") {
			continue
		}

		//Check if message from user is text
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			command := ParseCommand(update.Message.Text)
			switch command.name {
			case "chat_id":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatInt(update.Message.Chat.ID, 10))
				bot.Send(msg)
			case "/rank":
				username := strings.SplitAfter(update.Message.Text, "/rank")[1]
				clearUsername := strings.TrimSpace(username)

				segments, err := tracker.GetStats(clearUsername, "psn")
				if err != nil || len(segments) == 0 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "something went wrong, please try later")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, formatUserInfo(segments[0].Stats))
					bot.Send(msg)
				}
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Меня писал очень плохой программист, и он не рассказал мне что значит это сообщение")
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
	messagesHandler()
}
