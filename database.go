package main

import (
	"encoding/json"
	"os"
)

type MemoryStorage struct {
	// UserName on UserId
	UserIds map[int64]string
	// UserId on TaskIds
	TaskList map[int64][]string
	// TaskId on TaskInfo
	Tasks map[string]TaskInfo
	// UserId on current Line
	Lines map[int64]LineType
	path  string
}

type IDatabase interface {
	GetAllUsers() map[int64]string
	GetUserId(userName string) int64
	GetUserTasks(userId int64) []TaskInfo
	GetTaskInfo(taskId string) TaskInfo
	GetLine(userId int64) *LineType

	IsTask(taskId string) bool

	AddUser(userName string, userId int64)
	AddUserTask(userId int64, taskId string)
	AddNewTask(task TaskInfo)

	ChangeTask(taskId string, task TaskInfo)
	ChangeLine(userId int64, line *LineType)

	Save() error
}

func (this *MemoryStorage) GetAllUsers() map[int64]string {
	return this.UserIds
}
func (this *MemoryStorage) GetUserId(id string) int64 {
	return 0
}
func (this *MemoryStorage) GetUserTasks(userId int64) []TaskInfo {
	taskIds := this.TaskList[userId]
	tasks := make([]TaskInfo, 0)
	for _, taskId := range taskIds {
		task := this.Tasks[taskId]
		tasks = append(tasks, task)
	}
	return tasks
}
func (this *MemoryStorage) GetTaskInfo(taskId string) TaskInfo {
	if _, ok := this.Tasks[taskId]; ok {
		return this.Tasks[taskId]
	}
	return TaskInfo{}
}
func (this *MemoryStorage) GetLine(userId int64) *LineType {
	if _, ok := this.Lines[userId]; ok {
		line := this.Lines[userId]
		return &line
	}
	return nil
}

func (this *MemoryStorage) IsTask(taskId string) bool {
	_, ok := this.Tasks[taskId]
	return ok
}

func (this *MemoryStorage) AddUser(userName string, userId int64) {
	if this.UserIds == nil {
		this.UserIds = make(map[int64]string)
	}
	this.UserIds[userId] = userName
}
func (this *MemoryStorage) AddUserTask(userId int64, taskId string) {
	if this.TaskList == nil {
		this.TaskList = make(map[int64][]string)
		this.TaskList[userId] = []string{}
	}
	tasks := append(this.TaskList[userId], taskId)
	this.TaskList[userId] = tasks
}
func (this *MemoryStorage) AddNewTask(task TaskInfo) {
	if this.Tasks == nil {
		this.Tasks = make(map[string]TaskInfo)
	}
	this.Tasks[task.Id] = task
}

func (this *MemoryStorage) ChangeTask(taskId string, task TaskInfo) {
	if _, ok := this.Tasks[taskId]; ok {
		this.Tasks[taskId] = task
	}
}

func (this *MemoryStorage) ChangeLine(userId int64, line *LineType) {
	if this.Lines == nil {
		this.Lines = make(map[int64]LineType)
	}
	if line == nil {
		return
	}
	if _, ok := this.Lines[userId]; !ok {
		this.Lines[userId] = *line
	} else {
		delete(this.Lines, userId)
	}
}

func (this *MemoryStorage) Save() error {
	data, err := json.MarshalIndent(&this, "", "	")
	if err != nil {
		return err
	}
	return os.WriteFile(this.path, data, 0644)
}

func GetLocalDatabase(path string) (IDatabase, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	file.Close()
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	db := &MemoryStorage{}
	if len(bytes) == 0 {
		db.path = path
	} else {
		err = json.Unmarshal(bytes, db)
		db.path = path
	}

	return IDatabase(db), err
}
