package main

import (
	"fmt"
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

func formatUserInfo(stats tracker.Stats) string {
	return fmt.Sprintf("Kills: %s\nDamage: %s\nRank RP: %s", stats.Kills.DisplayValue, stats.Damage.DisplayValue, stats.RankScore.DisplayValue)
}

// ChatIDHandler handler for command chat_id
func ChatIDHandler(bot *tgbotapi.BotAPI, chatID int64, command Command) {
	msg := tgbotapi.NewMessage(chatID, strconv.FormatInt(chatID, 10))
	bot.Send(msg)
}

// SubscribeHandler handler
func SubscribeHandler(bot *tgbotapi.BotAPI, chatID int64, command Command) {
	if len(command.args) < 1 {
		msg := tgbotapi.NewMessage(chatID, "you must provide username as argument")
		bot.Send(msg)
		return
	}

	username := command.args[0]
	err := Subscribe(username, chatID)

	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "something went wrong, please try later")
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("user %s subscribed for this chat", username))
		bot.Send(msg)
	}
}

// UnsubscribeHandler handler
func UnsubscribeHandler(bot *tgbotapi.BotAPI, chatID int64, command Command) {
	if len(command.args) < 1 {
		msg := tgbotapi.NewMessage(chatID, "you must provide username as argument")
		bot.Send(msg)
		return
	}

	username := command.args[0]
	err := Unsubscribe(username, chatID)

	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "something went wrong, please try later")
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("user %s unsubscribed for this chat", username))
		bot.Send(msg)
	}
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

// StatsHandler handler
func StatsHandler(bot *tgbotapi.BotAPI, chatID int64, command Command) {
	if len(command.args) < 2 {
		msg := tgbotapi.NewMessage(chatID, "you must provide username and legend name as argument")
		bot.Send(msg)
		return
	}

	username := command.args[0]
	legend := strings.ToLower(command.args[1])
	platform := "psn"
	if len(command.args) >= 3 {
		platform = command.args[2]
	}

	segments, err := tracker.GetStats(username, platform)

	if err != nil || len(segments) == 0 {
		msg := tgbotapi.NewMessage(chatID, "something went wrong, please try later")
		bot.Send(msg)
		return
	}

	for _, segment := range segments {
		if strings.ToLower(segment.Metadata.Name) == legend {
			msg := tgbotapi.NewPhotoUpload(chatID, nil)
			msg.FileID = segment.Metadata.TallImageURL
			msg.UseExisting = true
			msg.Caption = formatUserInfo(segment.Stats)
			bot.Send(msg)
			return
		}
	}

	msg := tgbotapi.NewMessage(chatID, "legend not found")
	bot.Send(msg)
}

// HelpHandler help command
func HelpHandler(bot *tgbotapi.BotAPI, chatID int64, command Command) {
	helpMessage := `
	Основные комманды:
	/help - вывести список доступных комманд
	
	/rank <username> [<platform>] - вывести статистику игрока: количество киллов, урона и очков рейтинга
	/stats <username> <class> [<platform>] - вывести стату легенды

	/subscribe <username> [<platform>] - добавить пользователя в список ежедневных рассылок статистики в данный чат
	/unsubscribe <username> - удалить игрока из списка ежедневных рассылок статистики 


	Дебаг комманды:
	/chat_id - получить идентификатор чата
	`
	msg := tgbotapi.NewMessage(chatID, helpMessage)
	bot.Send(msg)
}
