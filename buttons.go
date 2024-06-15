package main

import (
	"strconv"

	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetStartButtons() tgApi.InlineKeyboardMarkup {
	keyboard := tgApi.NewInlineKeyboardMarkup(
		tgApi.NewInlineKeyboardRow(
			tgApi.NewInlineKeyboardButtonData("Мои задачи", "/my_tasks"),
		),
		tgApi.NewInlineKeyboardRow(
			tgApi.NewInlineKeyboardButtonData("Новая задача", "/new_task"),
		),
		tgApi.NewInlineKeyboardRow(
			tgApi.NewInlineKeyboardButtonData("Назначить задачу", "/users"),
		))
	return keyboard
}

func GetTaskButtons(taskId string) tgApi.InlineKeyboardMarkup {
	keyboard := tgApi.NewInlineKeyboardMarkup(
		tgApi.NewInlineKeyboardRow(
			tgApi.NewInlineKeyboardButtonData("Завершить", "/finish@"+taskId),
			tgApi.NewInlineKeyboardButtonData("Удалить", "/delete@"+taskId),
		))
	return keyboard
}

func GetUserButton(self int64, users map[int64]string) tgApi.InlineKeyboardMarkup {
	rows := make([]tgApi.InlineKeyboardButton, 0)
	for userId, userName := range users {
		row := tgApi.NewInlineKeyboardButtonData(userName, "/new_task@"+strconv.Itoa(int(userId)))
		rows = append(rows, row)
	}
	keyboard := tgApi.NewInlineKeyboardMarkup(rows)

	return keyboard
}
