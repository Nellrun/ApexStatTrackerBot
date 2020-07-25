package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	db "github.com/heroku/go-apex-tracker/postgresql-db"
)

var (
	bot *tgbotapi.BotAPI
)

func webhookHandler(c *gin.Context) {
	defer c.Request.Body.Close()

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return
	}

	if update.Message == nil {
		return
	}

	if update.Message.Photo != nil {
		imageURL := (*update.Message.Photo)[0].FileID
		command := ParseCommand(update.Message.Caption)
		if command.name == "uploadimage" {
			if len(command.args) < 1 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "you must provide image tag")
				bot.Send(msg)
				return
			}
			err := db.AddImage(strings.ToLower(command.args[0]), imageURL)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "something went wrong")
				bot.Send(msg)
				return
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("image for tag %s updated", command.args[0]))
			bot.Send(msg)
			return
		}
	}

	if !strings.HasPrefix(update.Message.Text, "/") {
		return
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

func main() {
	err := db.CreateTables()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	//Create bot
	bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://apex-stat-tracker-bot.herokuapp.com/" + bot.Token))
	if err != nil {
		log.Fatal(err)
	}

	// gin router
	router := gin.New()
	router.Use(gin.Logger())

	router.POST("/"+bot.Token, webhookHandler)

	err = router.Run(":" + port)
	if err != nil {
		log.Println(err)
	}

}
