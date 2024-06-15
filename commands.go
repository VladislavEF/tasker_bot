package main

func (this *AnswerInfo) Start() {
	this.SetAnswer("Привет.\nЯ бот для заведения и хранения твоих задач. Приступим?")
	this.keyboard = GetStartButtons()
}

func (this *AnswerInfo) NewTask(line *Line, db IDatabase) {
	if line == nil {
		line = &Line{
			Command:  "/newTask",
			Argument: "",
		}
		this.SetAnswer("Что нужно сделать?")
	} else {
		task := NewTask(line.Argument)
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

func (this *AnswerInfo) MyTasks(db IDatabase) {
	tasks := db.GetUserTasks(this.userId)
	if len(tasks) == 0 {
		this.SetAnswer("Нет текущих задач")
	} else {
		this.SetAnswers(tasks)
	}
}
