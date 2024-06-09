package main

import (
	"os"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

var opts struct {
	DatabasePath string `long:"database" description:"Path to local database"`
}

func main() {
	log.Info("Start TG bot")

	if _, err := flags.Parse(&opts); err != nil {
		log.Error("Failed to parse options")
		os.Exit(1)
	}

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

	if opts.DatabasePath == "" {
		log.Fatal("Only local database is implemented")
		os.Exit(1)
	}
	db, err := GetLocalDatabase(opts.DatabasePath)
	if err != nil {
		log.WithError(err).Fatal("Can't connect to local database")
		os.Exit(1)
	}

	for update := range bot.GetUpdatesChan(config) {
		msg, err := ListenMessage(update)
		if err != nil {
			log.WithError(err).Error("Failed to process update")
			continue
		}
		ans := GetAnswerOnMessage(msg, db)
		err = SendAnswer(ans, bot)
		if err != nil {
			log.WithError(err).Error("Failed to send answer")
			continue
		}
		if err = db.Save(); err != nil{
			log.WithError(err).Fatal("Can't write to local database")
			os.Exit(1)
		}
	}

	log.Info("All Ok")
	os.Exit(0)
}
