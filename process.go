package main

import (
	"errors"
	"strings"

	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageInfo struct {
	text         string
	userId       int64
	userName     string
	userFullName string
}

type CallbackInfo struct {
	msg  *MessageInfo
	id   string
	data string
}

type AnswerInfo struct {
	text     []string
	userId   int64
	keyboard *tgApi.InlineKeyboardMarkup
}

type Line struct {
	Command  string
	Argument string
}

var commands = []string{"/start", "/tasks", "/myTasks", "/newTask", "/stats"}

func ListenMessage(update tgApi.Update) (info *MessageInfo, err error) {
	msg := update.Message
	info = &MessageInfo{
		text:         msg.Text,
		userId:       msg.From.ID,
		userName:     msg.From.UserName,
		userFullName: msg.From.FirstName + " " + msg.From.LastName,
	}
	return info, nil
}

func ListenCallback(update tgApi.Update) (callback *CallbackInfo, err error) {
	callbackInfo := update.CallbackQuery
	msg := callbackInfo.Message
	msgInfo := &MessageInfo{
		text:         callbackInfo.Data,
		userId:       callbackInfo.From.ID,
		userName:     msg.From.UserName,
		userFullName: msg.From.FirstName + " " + msg.From.LastName,
	}

	callback = &CallbackInfo{
		msg:  msgInfo,
		id:   callbackInfo.ID,
		data: callbackInfo.Data,
	}

	return callback, nil
}

func NewAnswer() *AnswerInfo {
	answer := &AnswerInfo{}
	return answer
}

func GetAnswerOnMessage(msg *MessageInfo, db IDatabase) *AnswerInfo {
	answer := NewAnswer()
	answer.userId = msg.userId

	if userId := db.GetUserId(msg.userName); userId == 0 {
		db.AddUser(msg.userName, msg.userId)
	}

	if errMessage := msg.Validate(); errMessage != "" {
		answer.text = []string{errMessage}
		return answer
	}

	msg.ProcessComand(answer, db)

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

func (this *MessageInfo) ProcessComand(answer *AnswerInfo, db IDatabase) {
	user := this.userId
	command := this.text

	line := db.GetLine(user)
	if this.IsCommand() {
		line = nil
	}
	if line != nil && db.IsOpenLine(user) {
		command = line.Command
		line.Argument = this.text
	}

	switch command {
	case "/start":
		answer.Start()
	case "/tasks":
		answer.SetAnswer("Раздел в разработке")
	case "/myTasks":
		answer.MyTasks(db)
	case "/newTask":
		answer.NewTask(line, db)
	case "/stats":
		answer.SetAnswer("Раздел в разработке")
	default:
		answer.SetAnswer("Неизвестная команда")
	}
}

func SendAnswer(answer *AnswerInfo, bot *tgApi.BotAPI) error {
	if answer == nil {
		return errors.New("Empty answer")
	}
	for _, text := range answer.text {
		msg := tgApi.NewMessage(answer.userId, text)
		msg.ReplyMarkup = answer.keyboard
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *AnswerInfo) SetAnswer(text string) {
	this.text = append(this.text, text)
}

func (this *AnswerInfo) SetAnswers(text []string) {
	this.text = text
}
