package main

import (
	tgApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetStartButtons() tgApi.InlineKeyboardMarkup {
	keyboard := tgApi.NewInlineKeyboardMarkup(
		tgApi.NewInlineKeyboardRow(
			tgApi.NewInlineKeyboardButtonData("Новая задача", "/newTask"),
		),
		tgApi.NewInlineKeyboardRow(
			tgApi.NewInlineKeyboardButtonData("Мои дела", "/myTasks"),
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
