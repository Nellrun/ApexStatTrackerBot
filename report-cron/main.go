package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tracker "github.com/heroku/go-apex-tracker/apex-tracker"
	db "github.com/heroku/go-apex-tracker/postgresql-db"
)

// Stat Describe one param
type Stat struct {
	OldStat float64
	NewStat float64
	Diff    float64
}

// StatElem elem
type StatElem struct {
	Type     string
	Kills    Stat
	Damage   Stat
	Rank     Stat
	ImageURL string
}

//TotalStat total
type TotalStat struct {
	total   StatElem
	legends []StatElem
}

func calculateStartPeriod() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, -1)
}

func calculateEndPeriod() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func findLegend(stats []tracker.Segment, legend string) *tracker.Segment {
	for _, segment := range stats {
		if strings.ToLower(segment.Metadata.Name) == strings.ToLower(legend) {
			return &segment
		}
	}
	return nil
}

func getLogStats(username string, date time.Time) *[]tracker.Segment {
	rawStats, err := db.GetLog(username, date)
	if err != nil {
		log.Print(err)
		return nil
	}

	var stats []tracker.Segment

	err = json.Unmarshal([]byte(rawStats), &stats)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &stats
}

func calcStat(newValue, oldValue tracker.Stat) Stat {
	return Stat{oldValue.Value, newValue.Value, newValue.Value - oldValue.Value}
}

func calcLegend(newStats, oldStats tracker.Segment) StatElem {
	return StatElem{
		strings.ToLower(newStats.Metadata.Name),
		calcStat(newStats.Stats.Kills, oldStats.Stats.Kills),
		calcStat(newStats.Stats.Damage, oldStats.Stats.Damage),
		calcStat(newStats.Stats.RankScore, oldStats.Stats.RankScore),
		newStats.Metadata.TallImageURL,
	}
}

func calcTotalStats(new, old []tracker.Segment) TotalStat {
	var stat TotalStat
	stat.total = calcLegend(new[0], old[0])

	var defaultLegend tracker.Segment

	for _, newLegend := range new[1:] {
		oldLegend := findLegend(old, newLegend.Metadata.Name)
		var legendScore StatElem
		if oldLegend == nil {
			legendScore = calcLegend(newLegend, defaultLegend)
		} else {
			legendScore = calcLegend(newLegend, *oldLegend)
		}

		if legendScore.Damage.Diff > 0 || legendScore.Kills.Diff > 0 {
			stat.legends = append(stat.legends, legendScore)
		}

	}

	sort.Slice(stat.legends, func(a, b int) bool {
		return stat.legends[a].Kills.Diff > stat.legends[b].Kills.Diff
	})

	return stat
}

func calculateStat(username string, startDate time.Time, endDate time.Time) *TotalStat {
	oldStats := getLogStats(username, startDate)
	newStats := getLogStats(username, endDate)

	if oldStats == nil || newStats == nil {
		log.Print("error while restoring stats")
		return nil
	}

	stats := calcTotalStats(*newStats, *oldStats)

	return &stats
}

func formatStatDiff(stat Stat) string {
	return fmt.Sprintf("%.0f (%+.0f)", stat.NewStat, stat.Diff)
}

func formatMessage(username string, stats TotalStat) string {
	text := fmt.Sprintf(
		"<<%s>>\nKills: %s\nDamage: %s\nRP: %s\n",
		username,
		formatStatDiff(stats.total.Kills),
		formatStatDiff(stats.total.Damage),
		formatStatDiff(stats.total.Rank),
	)

	for _, legendStat := range stats.legends {
		text += fmt.Sprintf(
			"\n<<%s>>\nKills: %s\nDamage: %s",
			legendStat.Type,
			formatStatDiff(legendStat.Kills),
			formatStatDiff(legendStat.Damage),
		)
	}

	return text
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		panic(err)
	}

	log.Print("getting subscription list")
	subscriptions, err := db.GetSubscriptionsToSend()
	if err != nil {
		log.Print(err)
		return
	}

	subscriptionsGrouped := make(map[string][]int64)
	for _, elem := range subscriptions {
		subscriptionsGrouped[elem.Username] = append(subscriptionsGrouped[elem.Username], elem.ChatID)
	}

	for username, chatIDs := range subscriptionsGrouped {
		log.Printf("processing %s", username)
		startDate := calculateStartPeriod()
		endDate := calculateEndPeriod()
		stats := calculateStat(username, startDate, endDate)

		if stats == nil {
			log.Print("failed to build stats, skipping")
			continue
		}

		if len(stats.legends) == 0 {
			log.Print("user without updates, skipping")
			continue
		}

		for _, chatID := range chatIDs {
			msg := tgbotapi.NewPhotoUpload(chatID, nil)

			imagePath := new(string)
			if len(stats.legends) == 0 {
				imagePath, _ = db.GetImage("default")
				log.Print("fallback to default image")
			} else {
				log.Print("retrieving image from base")
				image, _ := db.GetImage(strings.ToLower(stats.legends[0].Type))
				if image == nil || *image == "" {
					log.Print("fail, fallback to url")
					*imagePath = stats.legends[0].ImageURL
				} else {
					imagePath = image
				}
			}

			log.Printf("image path %s", *imagePath)

			msg.FileID = *imagePath
			msg.UseExisting = true

			text := formatMessage(username, *stats)
			msg.Caption = text
			bot.Send(msg)
		}
	}

}
