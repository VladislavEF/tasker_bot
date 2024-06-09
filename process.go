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

type Line struct {
	command  string
	argument string
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

func GetAnswerOnMessage(msg *MessageInfo, db Database) *AnswerInfo {
	answer := &AnswerInfo{
		userId: msg.userId,
	}

	if userId := db.GetUserId(msg.userName); userId == 0 {
		db.AddUser(msg.userName, msg.userId)
	}

	if errMessage := msg.Validate(); errMessage != "" {
		answer.text = errMessage
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

func (this *MessageInfo) ProcessComand(answer *AnswerInfo, db Database) {
	user := this.userId
	command := this.text
	line := db.GetLine(user)
	if db.IsOpenLine(user) {
		command = line.command
	}

	switch command {
	case "/start":
		answer.text = "Привет.\nЯ бот для заведения и хранения твоих задач. Приступим?"
	case "/tasks":
		answer.text = "Раздел в разработке"
	case "/userTasks":
		taskIds := db.GetUserTasks(user)
		tasks := ""
		for _, id := range taskIds {
			task := db.GetTaskInfo(id)
			tasks += task.name + "\n"
		}
		answer.text = tasks
	case "/newTask":
		if line == nil {
			line = &Line{
				command:  command,
				argument: "",
			}
			answer.text = "Что нужно сделать?"
		} else {
			task := NewTask(line.argument)
			db.AddNewTask(*task)
			db.AddUserTask(user, task.id)
			answer.text = "Задача добавлена"
		}
		db.ChangeLine(user, *line)
	case "/stats":
		answer.text = "Раздел в разработке"
	default:
		answer.text = "Неизвестная команда"
	}

	return
}

func SendAnswer(answer *AnswerInfo, bot *tg.BotAPI) error {
	_, err := bot.Send(tg.NewMessage(answer.userId, answer.text))
	return err
}
