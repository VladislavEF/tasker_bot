package main

import (
	"os"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func main(){
	log.Info("Start TG bot")

	bot, err := tg.NewBotAPI("tasker key")
	if err != nil {
		log.Panic(err)
	}

	// some code
	
	log.Info("All Ok")
	os.Exit(0)
}