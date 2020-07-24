package main

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tracker "github.com/heroku/go-apex-tracker/apex-tracker"
)

// Command bot command structure
type Command struct {
	name string
	args []string
}

// ParseCommand parsing text commands
func ParseCommand(message string) Command {
	var command Command
	trimmedMessage := strings.TrimSpace(message)
	if !strings.HasPrefix(trimmedMessage, "/") {
		return command
	}

	splited := strings.Split(trimmedMessage, " ")
	command.name = splited[0][1:]

	for _, elem := range splited[1:] {
		if elem != "" {
			command.args = append(command.args, elem)
		}
	}

	return command
}

// ChatIDHandler handler for command chat_id
func ChatIDHandler(bot *tgbotapi.BotAPI, chatID int64, command Command) {
	msg := tgbotapi.NewMessage(chatID, strconv.FormatInt(chatID, 10))
	bot.Send(msg)
}

// RankHandler handler
func RankHandler(bot *tgbotapi.BotAPI, chatID int64, command Command) {
	if len(command.args) < 1 {
		msg := tgbotapi.NewMessage(chatID, "you must provide username as argument")
		bot.Send(msg)
		return
	}

	username := command.args[0]
	platform := "psn"
	if len(command.args) >= 2 {
		platform = command.args[1]
	}
	segments, err := tracker.GetStats(username, platform)
	if err != nil || len(segments) == 0 {
		msg := tgbotapi.NewMessage(chatID, "something went wrong, please try later")
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, formatUserInfo(segments[0].Stats))
		bot.Send(msg)
	}
}
