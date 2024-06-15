package main

import (
	"errors"
	"fmt"

	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Answer struct {
	text     []string
	userId   int64
	keyboard []tgApi.InlineKeyboardMarkup
}

func NewAnswer() *Answer {
	return &Answer{}
}

func NewUserAnswer(userId int64) *Answer {
	answer := NewAnswer()
	answer.userId = userId
	return answer
}

func (this *Answer) Start() {
	this.SetAnswer("Привет.\nЯ бот для заведения и хранения твоих задач. Приступим?")
	this.SetButton(GetStartButtons())
}

func (this *Answer) NewTask(line *LineType, db IDatabase) {
	if line == nil {
		line = &LineType{
			Command: "/newTask",
			Text:    "",
		}
		this.SetAnswer("Что нужно сделать?")
	} else {
		task := NewTask(line.Text)
		if !db.IsTask(task.Id) {
			if task.Name == "" {
				this.SetAnswer("Нет текста задачи")
				return
			}
			db.AddNewTask(*task)
			db.AddUserTask(this.userId, task.Id)
			this.SetAnswer("Задача добавлена")
		} else {
			this.SetAnswer("Задача уже заведена")
		}
	}
	db.ChangeLine(this.userId, *line)
}

func (this *Answer) DeleteTask(id string, db IDatabase) {
	if db.IsTask(id) {
		task := db.GetTaskInfo(id)
		task.Cancelled()
		fmt.Println(task)
		db.ChangeTask(id, task)
	}
	this.SetAnswer("Удалена")
}

func (this *Answer) FinishTask(id string, db IDatabase) {
	if db.IsTask(id) {
		task := db.GetTaskInfo(id)
		task.Done()
		db.ChangeTask(id, task)
		this.SetAnswer("Поздравляю!")
	} else {
		this.SetAnswer("Нет такой задачи")
	}
}

func (this *Answer) MyTasks(db IDatabase) {
	tasks := db.GetUserTasks(this.userId)
	for _, task := range tasks {
		if task.Name == "" || task.Status != Backlog {
			continue
		}
		this.SetAnswer(task.Name)
		this.SetButton(GetTaskButtons(task.Id))
	}
	if len(this.text) == 0 {
		this.SetAnswer("Нет текущих задач")
	}
}

func (this *Answer) SetButton(button tgApi.InlineKeyboardMarkup) {
	this.keyboard = append(this.keyboard, button)
}

func (this *Answer) SetAnswer(text string) {
	this.text = append(this.text, text)
}

func (this *Answer) SetAnswers(text []string) {
	this.text = text
}

func SendStartMsg(bot *tgApi.BotAPI, db IDatabase) error {
	baseMsg := NewAnswer()
	baseMsg.Start()
	return SendMessageToUsers(db.GetAllUsers(), baseMsg, bot)
}

func SendMessageToUsers(users []int64, answer *Answer, bot *tgApi.BotAPI) error {
	if answer == nil {
		return errors.New("Empty answer")
	}
	for _, user := range users {
		answer.userId = user
		SendMessage(answer, bot)
	}
	return nil
}

func SendMessage(outMsg *Answer, bot *tgApi.BotAPI) error {
	if outMsg == nil {
		return errors.New("Empty answer")
	}
	for i, text := range outMsg.text {
		msg := tgApi.NewMessage(outMsg.userId, text)
		if len(outMsg.text) == len(outMsg.keyboard) {
			msg.ReplyMarkup = outMsg.keyboard[i]
		}
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
