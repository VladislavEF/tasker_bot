package main

import (
	"os"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Start TG bot")

	apiKey := os.Getenv("TASKER_KEY")
	if apiKey == "" {
		log.Error("No API key")
		os.Exit(1)
	}

	bot, err := tg.NewBotAPI(apiKey)
	if err != nil {
		log.WithError(err).Error("Failed to create bot API")
		os.Exit(1)
	}

	config := tg.NewUpdate(0)
	config.Timeout = 60

	for update := range bot.GetUpdatesChan(config) {
		msg, err := ListenMessage(update)
		if err != nil {
			log.WithError(err).Error("Failed to process update")
			continue
		}
		ans := GetAnswerOnMessage(msg)
		err = SendAnswer(ans, bot)
		if err != nil {
			log.WithError(err).Error("Failed to send answer")
			continue
		}
	}

	log.Info("All Ok")
	os.Exit(0)
}
