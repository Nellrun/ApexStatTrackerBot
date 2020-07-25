package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tracker "github.com/heroku/go-apex-tracker/apex-tracker"
	db "github.com/heroku/go-apex-tracker/postgresql-db"
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

func formatStatDiff(oldStat float64, newStat float64) string {
	diff := newStat - oldStat
	return fmt.Sprintf("%.0f (%+.0f)", newStat, diff)
}

func findLegend(stats []tracker.Segment, legend string) *tracker.Segment {
	for _, segment := range stats {
		if strings.ToLower(segment.Metadata.Name) == strings.ToLower(legend) {
			return &segment
		}
	}
	return nil
}

func formatDailyStats(oldStats []tracker.Segment, newStats []tracker.Segment, legend string) string {
	if legend == "" {
		text := fmt.Sprintf(
			"===Total===\nKills: %s\nDamage: %s\nRank RP: %s\n",
			formatStatDiff(oldStats[0].Stats.Kills.Value, newStats[0].Stats.Kills.Value),
			formatStatDiff(oldStats[0].Stats.Damage.Value, newStats[0].Stats.Damage.Value),
			formatStatDiff(oldStats[0].Stats.RankScore.Value, newStats[0].Stats.RankScore.Value))

		for _, newStat := range newStats[1:] {
			oldStat := findLegend(oldStats, newStat.Metadata.Name)

			if oldStat == nil {
				continue
			}

			text += fmt.Sprintf(
				"\n\n===%s===\nKills: %s\nDamage: %s\nRank RP: %s",
				newStat.Metadata.Name,
				formatStatDiff(oldStat.Stats.Kills.Value, newStat.Stats.Kills.Value),
				formatStatDiff(oldStat.Stats.Damage.Value, newStat.Stats.Damage.Value),
				formatStatDiff(oldStat.Stats.RankScore.Value, newStat.Stats.RankScore.Value))
		}

		return text
	}

	oldLegend := findLegend(oldStats, legend)
	newLegend := findLegend(newStats, legend)

	if oldLegend == nil || newLegend == nil {
		return "legend not found"
	}

	return fmt.Sprintf(
		"Kills: %s\nDamage: %s\nRank RP: %s",
		formatStatDiff(oldLegend.Stats.Kills.Value, newLegend.Stats.Kills.Value),
		formatStatDiff(oldLegend.Stats.Damage.Value, newLegend.Stats.Damage.Value),
		formatStatDiff(oldStats[0].Stats.RankScore.Value, newStats[0].Stats.RankScore.Value))
}

func getStat(stat tracker.Stat) string {
	if stat.DisplayValue == "" {
		return "unknown"
	}
	return stat.DisplayValue
}

func formatRankedValue(stats tracker.Stats) string {
	return fmt.Sprintf("%s (%s)", stats.RankScore.DisplayValue, stats.RankScore.Metadata.RankName)
}

func formatUserInfo(stats tracker.Stats, globalStats tracker.Stats) string {
	return fmt.Sprintf("Kills: %s\nDamage: %s\nRank RP: %s", getStat(stats.Kills), getStat(stats.Damage), formatRankedValue(globalStats))
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

	username := strings.ToLower(command.args[0])
	err := db.Subscribe(username, chatID)

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

	username := strings.ToLower(command.args[0])
	err := db.Unsubscribe(username, chatID)

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

	errMessage := CheckPlatform(platform)
	if errMessage != nil {
		msg := tgbotapi.NewMessage(chatID, *errMessage)
		bot.Send(msg)
		return
	}

	segments, err := tracker.GetStats(username, platform)
	if err != nil || len(segments) == 0 {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("user %s not found for platfrom %s", username, platform))
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, formatUserInfo(segments[0].Stats, segments[0].Stats))
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

	errMessage := CheckPlatform(platform)
	if errMessage != nil {
		msg := tgbotapi.NewMessage(chatID, *errMessage)
		bot.Send(msg)
		return
	}

	segments, err := tracker.GetStats(username, platform)

	if err != nil || len(segments) == 0 {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("user %s not found for platfrom %s", username, platform))
		bot.Send(msg)
		return
	}

	for _, segment := range segments {
		if strings.ToLower(segment.Metadata.Name) == legend {
			msg := tgbotapi.NewPhotoUpload(chatID, nil)
			image, err := db.GetImage(legend)
			if err == nil && image != nil {
				msg.FileID = *image
			} else {
				msg.FileID = segment.Metadata.TallImageURL
			}
			msg.UseExisting = true
			msg.Caption = formatUserInfo(segment.Stats, segments[0].Stats)
			bot.Send(msg)
			return
		}
	}

	msg := tgbotapi.NewMessage(chatID, "legend not found")
	bot.Send(msg)
}

// DailyStatsHandler handler
func DailyStatsHandler(bot *tgbotapi.BotAPI, chatID int64, command Command) {
	if len(command.args) < 1 {
		msg := tgbotapi.NewMessage(chatID, "you must provide username as argument")
		bot.Send(msg)
		return
	}

	username := command.args[0]
	legend := ""
	if len(command.args) >= 2 {
		legend = strings.ToLower(command.args[1])
	}

	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	yesterday := today.AddDate(0, 0, -1)

	rawStats, err := db.GetLog(username, today)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "today's stats retrieve operation failed")
		bot.Send(msg)
		return
	}

	oldRawStats, err := db.GetLog(username, yesterday)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "yesterday's stats retrieve operation failed")
		bot.Send(msg)
		return
	}

	var todayStats []tracker.Segment
	var yesterdayStats []tracker.Segment

	err = json.Unmarshal([]byte(rawStats), &todayStats)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "failed to parse today's stats")
		bot.Send(msg)
		return
	}

	err = json.Unmarshal([]byte(oldRawStats), &yesterdayStats)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "failed to parse yesterday's stats")
		bot.Send(msg)
		return
	}

	text := formatDailyStats(yesterdayStats, todayStats, legend)
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}

// HelpHandler help command
func HelpHandler(bot *tgbotapi.BotAPI, chatID int64, command Command) {
	helpMessage := `
	Основные комманды:
	/help - вывести список доступных комманд
	
	/rank <username> [<platform>] - вывести статистику игрока: количество киллов, урона и очков рейтинга
	/stats <username> <legend> [<platform>] - вывести стату легенды

	/subscribe <username> [<platform>] - добавить пользователя в список ежедневных рассылок статистики в данный чат
	/unsubscribe <username> - удалить игрока из списка ежедневных рассылок статистики 
	/dailystats <username> [<legend>] - получить общую\легендную статистику за сегодня

	Дебаг комманды:
	/chat_id - получить идентификатор чата
	`
	msg := tgbotapi.NewMessage(chatID, helpMessage)
	bot.Send(msg)
}

// DeleteImageHandler handler
func DeleteImageHandler(bot *tgbotapi.BotAPI, chatID int64, command Command) {
	if len(command.args) < 1 {
		msg := tgbotapi.NewMessage(chatID, "you must provide imagetag as argument")
		bot.Send(msg)
		return
	}

	imageTag := command.args[0]
	err := db.DeleteImage(strings.ToLower(imageTag))
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "something went wrong, please try later")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "image deleted")
	bot.Send(msg)
}
