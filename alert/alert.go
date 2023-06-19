package alert

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"os"
	"strconv"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	timeDateFormat = "2006-01-02 15:04:05"
)

var botToken string
var chatID int64

func init() {
    botToken = os.Getenv("bottoken")
    chatIDStr := os.Getenv("chatid")
    chatID, _ = strconv.ParseInt(chatIDStr, 10, 64)
	log.Println("sending to chat: ", chatID)
}

type alertmanagerAlert struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []struct {
		Status string `json:"status"`
		Labels struct {
			Name      string `json:"name"`
			Instance  string `json:"instance"`
			Alertname string `json:"alertname"`
			Service   string `json:"service"`
			Severity  string `json:"severity"`
		} `json:"labels"`
		Annotations struct {
			Info        string `json:"info"`
			Description string `json:"description"`
			Summary     string `json:"summary"`
		} `json:"annotations"`
		StartsAt     time.Time `json:"startsAt"`
		EndsAt       time.Time `json:"endsAt"`
		GeneratorURL string    `json:"generatorURL"`
		Fingerprint  string    `json:"fingerprint"`
	} `json:"alerts"`
	GroupLabels struct {
		Alertname string `json:"alertname"`
	} `json:"groupLabels"`
	CommonLabels struct {
		Alertname string `json:"alertname"`
		Service   string `json:"service"`
		Severity  string `json:"severity"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
		Summary string `json:"summary"`
	} `json:"commonAnnotations"`
	ExternalURL string `json:"externalURL"`
	Version     string `json:"version"`
	GroupKey    string `json:"groupKey"`
}

// ToTelegram function responsible to send msg to telegram
func ToTelegram(w http.ResponseWriter, r *http.Request) {

	var alerts alertmanagerAlert

	bot, err := botapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	_ = json.NewDecoder(r.Body).Decode(&alerts)

	for _, alert := range alerts.Alerts {
		telegramMsg := "Status: " + alerts.Status + "\n"
		if alert.Labels.Name != "" {
			telegramMsg += "Instance: " + alert.Labels.Instance + "(" + alert.Labels.Name + ")\n"
		}
		if alert.Annotations.Info != "" {
			telegramMsg += "Info: " + alert.Annotations.Info + "\n"
		}
		if alert.Annotations.Summary != "" {
			telegramMsg += "Summary: " + alert.Annotations.Summary + "\n"
		}
		if alert.Annotations.Description != "" {
			telegramMsg += "Description: " + alert.Annotations.Description + "\n"
		}
		if alert.Status == "resolved" {
			telegramMsg += "Resolved: " + alert.EndsAt.Format(timeDateFormat)
		} else if alert.Status == "firing" {
			telegramMsg += "Started: " + alert.StartsAt.Format(timeDateFormat)
		}

		msg := botapi.NewMessage(chatID, telegramMsg)
		bot.Send(msg)
	}

	log.Println(alerts)
	json.NewEncoder(w).Encode(alerts)

}
