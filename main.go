package main

import (
	"os"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type UpdateInfo struct {
	bot      *tg.BotAPI
	update   tg.Update
	text     string
	id       int64
	user     string
	fullName string
}

func Listen(bot *tg.BotAPI, update tg.Update) (err error) {
	if update.Message == nil {
		return
	}

	info := UpdateInfo{
		bot:    bot,
		update: update,
	}

	info.GetMessage()

	if err := info.SendAnswerWithName("Принято"); err != nil {
		log.Error("Answer is failed")
		return err
	}

	return
}

func (info *UpdateInfo) GetMessage() {
	msg := info.update.Message
	info.id = msg.From.ID
	info.text = msg.Text
	info.user = msg.From.UserName
	info.fullName = msg.From.FirstName + " " + msg.From.LastName
}

func (info *UpdateInfo) SendAnswerWithName(text string) error {
	return info.SendAnswer(text + ", " + info.fullName)
}

func (info *UpdateInfo) SendAnswer(text string) error {
	_, err := info.bot.Send(tg.NewMessage(info.id, text))
	return err
}

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
		if err = Listen(bot, update); err != nil {
			log.WithError(err).Error("Failed to process update")
			continue
		}
	}

	// some code

	log.Info("All Ok")
	os.Exit(0)
}
