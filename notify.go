package main

import (
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jasonlvhit/gocron"
)

func Notification(db IDatabase, bot *tgApi.BotAPI) {
	for userId := range db.GetAllUsers() {
		answer := NewUserAnswer(userId)
		answer.MyTasks(db)
		if len(answer.text) != 0{
			SendMessage(answer, bot)
		}
	}
}

func Notifyer(db IDatabase, bot *tgApi.BotAPI) {
	s := gocron.NewScheduler()
	s.Every(1).Day().At("11:00:00").Do(Notification, db, bot)
	s.Every(1).Day().At("19:00:00").Do(Notification, db, bot)
	<-s.Start()
}
