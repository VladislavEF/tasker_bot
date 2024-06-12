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
	text   []string
	userId int64
}

type Line struct {
	Command  string
	Argument string
}

var commands = []string{"/start", "/tasks", "/myTasks", "/newTask", "/stats"}

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

func GetAnswerOnMessage(msg *MessageInfo, db IDatabase) *AnswerInfo {
	answer := &AnswerInfo{
		userId: msg.userId,
	}

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
	if this.IsCommand(){
		line = nil
	}
	if line != nil && db.IsOpenLine(user) {
		command = line.Command
		line.Argument = this.text
	}

	switch command {
	case "/start":
		answer.SetAnswer("Привет.\nЯ бот для заведения и хранения твоих задач. Приступим?")
	case "/tasks":
		// answer.text = ToString(db.GetUserTasks(user))
		answer.SetAnswer("Раздел в разработке")
	case "/myTasks":
		tasks := db.GetUserTasks(user)
		if len(tasks) == 0 {
			answer.SetAnswer("Нет текущих задач")
		} else {
			answer.SetAnswers(tasks)
		}
	case "/newTask":
		if line == nil {
			line = &Line{
				Command:  command,
				Argument: "",
			}
			answer.SetAnswer("Что нужно сделать?")
		} else {
			task := NewTask(line.Argument)
			if !db.IsTask(task.Id){
				db.AddNewTask(*task)
				db.AddUserTask(user, task.Id)
				answer.SetAnswer("Задача добавлена")
			} else {
				answer.SetAnswer("Задача уже заведена")
			}
		}
		db.ChangeLine(user, *line)
	case "/stats":
		answer.SetAnswer("Раздел в разработке")
	default:
		answer.SetAnswer("Неизвестная команда")
	}
}

func SendAnswer(answer *AnswerInfo, bot *tg.BotAPI) error {
	for _,text := range answer.text{
		_, err := bot.Send(tg.NewMessage(answer.userId, text))
		if err != nil{
			return err
		}
	}
	return nil
}

func (this *AnswerInfo) SetAnswer(text string){
	this.text = append(this.text, text)
}

func (this *AnswerInfo) SetAnswers(text []string){
	this.text = text
}
