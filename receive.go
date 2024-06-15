package main

import (
	"strconv"
	"strings"

	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageType struct {
	text         string
	userId       int64
	userName     string
	userFullName string
}

type CallbackType struct {
	msg  *MessageType
	id   string
	data string
}

type LineType struct {
	Command  string
	Text     string
	Executor int64
}

var functions = []string{"/start", "/tasks", "/my_tasks", "/new_task", "/stats", "/finish", "/delete", "/users"}

func ListenMessage(update tgApi.Update) (info *MessageType, err error) {
	msg := update.Message
	info = &MessageType{
		text:         msg.Text,
		userId:       msg.From.ID,
		userName:     msg.From.UserName,
		userFullName: msg.From.FirstName + " " + msg.From.LastName,
	}
	return info, nil
}

func ListenCallback(update tgApi.Update) (callback *CallbackType, err error) {
	callbackInfo := update.CallbackQuery
	msg := callbackInfo.Message
	msgInfo := &MessageType{
		text:         callbackInfo.Data,
		userId:       callbackInfo.From.ID,
		userName:     msg.From.UserName,
		userFullName: msg.From.FirstName + " " + msg.From.LastName,
	}

	callback = &CallbackType{
		msg:  msgInfo,
		id:   callbackInfo.ID,
		data: callbackInfo.Data,
	}

	return callback, nil
}

func (this *MessageType) IsCommand() bool {
	return strings.HasPrefix(this.text, "/")
}

func (this *MessageType) IsMenu() bool {
	return this.text == "Меню" || this.text == "Команды" || this.text == "Menu" || this.text == "Commands" || this.text == "Старт"
}

func GetAnswerOnMessage(msg *MessageType, db IDatabase) *Answer {
	answer := NewUserAnswer(msg.userId)

	if userId := db.GetUserId(msg.userName); userId == 0 {
		if msg.userName != selfName {
			db.AddUser(msg.userName, msg.userId)
		}
	}

	if errMessage := msg.Validate(); errMessage != "" {
		answer.text = []string{errMessage}
		return answer
	}

	msg.ProcessComand(answer, db)

	return answer
}

func (this *MessageType) Validate() (errMessage string) {
	if this.text == "" {
		return "Получено пустое сообщение"
	}

	if !this.IsCommand() {
		return ""
	}

	isNownCommand := false
	function := GetFunction(this.text)
	for _, knownFunctions := range functions {
		if function == knownFunctions {
			isNownCommand = true
			break
		}
	}
	if !isNownCommand {
		return "Неизвестная команда"
	}
	if function == "/delete" || function == "/finish" {
		if GetArgument(this.text) == "" {
			return "Операция с неизвестной задачей"
		}
	}

	return ""
}

func (this *MessageType) ProcessComand(answer *Answer, db IDatabase) {
	user := this.userId

	if this.IsMenu() {
		answer.SetAnswer("Что нужно?")
		answer.SetButton(GetStartButtons())
		return
	}

	function, argument := GetFunction(this.text), GetArgument(this.text)

	line := db.GetLine(user)
	if db.IsOpenLine(user) {
		if this.IsCommand() {
			db.ChangeLine(user, *line)
		} else {
			function = line.Command
			line.Text = this.text
			if argument != "" {
				userId, _ := strconv.Atoi(argument)
				line.Executor = int64(userId)
			} else {
				line.Executor = user
			}
		}
	}

	switch function {
	case "/start":
		answer.Start()
	case "/tasks":
		answer.SetAnswer("Раздел в разработке")
	case "/my_tasks":
		answer.MyTasks(db)
	case "/new_task":
		answer.NewTask(line, db)
	case "/stats":
		answer.SetAnswer("Раздел в разработке")
	case "/finish":
		answer.FinishTask(argument, db)
	case "/delete":
		answer.DeleteTask(argument, db)
	case "/users":
		answer.UsersList(db)
	default:
		answer.SetAnswer("Неизвестная команда")
		answer.SetButton(GetStartButtons())
	}
}

func GetFunction(command string) string {
	return strings.Split(command, "@")[0]
}

func GetArgument(command string) string {
	comm := strings.Split(command, "@")
	if len(comm) > 1 {
		return comm[1]
	}
	return ""
}
