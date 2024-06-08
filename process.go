package main

import (
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageInfo struct {
	update       tg.Update
	text         string
	userId       int64
	userName     string
	userFullName string
}

type AnswerInfo struct {
	text   string
	userId int64
}

var commands = []string{"/start", "/tasks", "/userTasks", "/newTask", "/stats"}

func ListenMessage(update tg.Update) (info *MessageInfo, err error) {
	if update.Message == nil {
		return
	}

	msg := update.Message
	info = &MessageInfo{
		update:       update,
		text:         msg.Text,
		userId:       msg.From.ID,
		userName:     msg.From.UserName,
		userFullName: msg.From.FirstName + " " + msg.From.LastName,
	}

	return info, nil
}

func GetAnswerOnMessage(msg *MessageInfo) *AnswerInfo {
	answer := &AnswerInfo{
		userId: msg.userId,
	}

	if errMessage := msg.Validate(); errMessage != "" {
		answer.text = errMessage
		return answer
	}

	// TODO:
	// 1) Проверка пришла команда или нет
	// 2) обработка команды, создание line
	// 3) создание списка задач
	// 4) выдача списка задач
	// 5) добавление задачи
	// 6) перевод задачи в статус
	// 7) уведомления о задачах

	return answer
}

func (this *MessageInfo) Validate() (errMessage string) {
	if this.text == "" {
		return "Получено пустое сообщение"
	}

	if this.IsCommand() {
		isNownCommand := false
		for _, command := range commands {
			if this.text == command {
				isNownCommand = true
				break
			}
		}
		if !isNownCommand {
			return "Неизвестная команда"
		}
	}

	return ""
}

func (this *MessageInfo) IsCommand() bool {
	return strings.HasPrefix(this.text, "/")
}

func SendAnswer(answer *AnswerInfo, bot *tg.BotAPI) error {
	_, err := bot.Send(tg.NewMessage(answer.userId, answer.text))
	return err
}
