package main

import (
	"errors"
	"os"

	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

var opts struct {
	DatabasePath string `long:"database" description:"Path to local database"`
}

var selfName = "tasker_list_bot"

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

	bot, err := tgApi.NewBotAPI(apiKey)
	if err != nil {
		log.WithError(err).Error("Failed to create bot API")
		os.Exit(1)
	}

	config := tgApi.NewUpdate(0)
	config.Timeout = 60

	if opts.DatabasePath == "" {
		opts.DatabasePath = os.Getenv("TASKER_LOCAL_DATABASE")
		if opts.DatabasePath == "" {
			log.Fatal("Only local database is implemented")
			os.Exit(1)
		}
	}
	db, err := GetLocalDatabase(opts.DatabasePath)
	if err != nil {
		log.WithError(err).Fatal("Can't connect to local database")
		os.Exit(1)
	}

	// if err := SendStartMsg(bot, db); err != nil {
	// 	log.WithError(err).Error("Failed to send start message")
	// 	os.Exit(1)
	// }

	go Notifyer(db, bot)

	for update := range bot.GetUpdatesChan(config) {
		var msg *MessageType
		var err error
		if update.CallbackQuery != nil {
			callback, _err := ListenCallback(update)
			err = _err
			msg = callback.msg
		} else if update.Message != nil {
			msg, err = ListenMessage(update)
		} else {
			err = errors.New("Unknown update")
		}
		if err != nil {
			log.WithError(err).Error("Failed to process update")
			continue
		}

		err = SendMessage(GetAnswerOnMessage(msg, db), bot)
		if err != nil {
			log.WithError(err).Error("Failed to send answer")
			continue
		}
		if err = db.Save(); err != nil {
			log.WithError(err).Fatal("Can't write to local database")
			os.Exit(1)
		}
	}

	log.Info("All Ok")
	os.Exit(0)
}
